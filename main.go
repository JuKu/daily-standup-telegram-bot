package main

import (
	"bufio"
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

// Struktur zur Speicherung der Benutzeraktivitäten
type UserActivity struct {
	UserID   int64
	UserName string
	Days     [5]bool
}

// Map zur Speicherung der Aktivitäten
var userActivities = make(map[int64]*UserActivity)

// map to store the chat ids
var chatIDs = make(map[int64]struct{})

var motivationalQuotes = []string{
	"Super Teamarbeit heute! Ihr seid spitze! 💪😊",
	"Einfach klasse, wie diszipliniert alle heute waren! 👏✨",
	"Ihr seid wirklich fleißig! Weiter so! 🛠️🔥",
	"Heute war wieder ein Tag voller Erfolg. Gut gemacht! 🌟👍",
	"Gemeinsam sind wir unschlagbar! Tolle Arbeit heute! 🤝💯",
	"Jeden Tag ein bisschen besser – weiter so! 🚀✨",
	"Das war wieder ein erfolgreicher Tag. Klasse Leistung, Team! 🥇🎉",
	"Mit so einem Einsatz kommen wir ganz nach oben! ⛰️🏆",
	"Heute wieder alles gegeben – das ist Teamgeist! 💪🙌",
	"Eure Disziplin und Motivation beeindrucken mich jeden Tag aufs Neue! 🧠💥",
	"Das war Spitzenarbeit! Ihr seid wirklich ein starkes Team! 💪🛡️",
	"Jeder von euch bringt uns ein Stück weiter – danke dafür! 🧩🙏",
	"Mit so einer Einstellung können wir alles erreichen! 🎯✨",
	"Wieder ein Schritt näher an unseren Zielen. Super gemacht! 🏃‍♂️🥅",
	"Teamwork makes the dream work – und ihr habt das heute wieder bewiesen! 💭🤝",
	"Toll, wie jeder von euch seinen Beitrag leistet. Gemeinsam sind wir unschlagbar! 🤜🤛💪",
	"Jeder von euch macht den Unterschied – danke für euren Einsatz! 🌟👏",
	"Heute war wieder ein Tag, an dem alles gepasst hat. Ihr seid großartig! 🎉🙌",
	"Euer Engagement ist ansteckend – weiter so! 🔥🚀",
	"Mit so viel Elan und Energie geht's steil nach oben! 📈💪",
}

func main() {
	// Get API Key from environmental variables
	apiKey := os.Getenv("TELEGRAM_BOT_API_KEY")
	if apiKey == "" {
		log.Fatal("TELEGRAM_BOT_API_KEY environment variable not set")
		return
	}

	log.Println("TELEGRAM_BOT_API_KEY is available")

	// Bot initialisieren
	pref := telebot.Settings{
		Token:  apiKey,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("loadChatIDs")
	loadChatIDs()

	// Function to collect the activity
	b.Handle(telebot.OnText, func(c telebot.Context) error {
		user := c.Sender()
		now := time.Now()

		log.Printf("Received message: %s", c.Message().Text)

		chatID := c.Chat().ID
		chatIDs[chatID] = struct{}{} // Speichern der Chat-ID
		saveChatIDs()

		// Save the activity for the current day (Monday - Friday)
		RecordActivity(user.ID, user.Username, now)

		return nil
	})

	// store the chat id the telegram bot was added to
	b.Handle(telebot.OnAddedToGroup, func(c telebot.Context) error {
		chatID := c.Chat().ID
		chatIDs[chatID] = struct{}{} // Speichern der Chat-ID
		saveChatIDs()
		sendWelcomeMessage(b, c.Chat())
		return nil
	})

	// Scheduler initialisieren
	initScheduler(b)

	// Lies die Nachrichten der Woche erneut aus, falls der Bot neu gestartet wurde
	// Diese Funktion müsste implementiert werden, um Nachrichten erneut auszulesen und die Aktivitäten zu rekonstruieren

	log.Println("Bot started successfully")

	// Bot starten
	b.Start()
}

// Funktion zum Initialisieren des Schedulers
func initScheduler(b *telebot.Bot) {
	s := gocron.NewScheduler(time.UTC)

	_, _ = s.Every(1).Monday().At("12:00").Do(func() {
		sendDailyStandUp(b)
	})
	_, _ = s.Every(1).Tuesday().At("12:00").Do(func() {
		sendDailyStandUp(b)
	})
	_, _ = s.Every(1).Wednesday().At("12:00").Do(func() {
		sendDailyStandUp(b)
	})
	_, _ = s.Every(1).Thursday().At("12:00").Do(func() {
		sendDailyStandUp(b)
	})
	_, _ = s.Every(1).Friday().At("12:00").Do(func() {
		sendDailyStandUp(b)
	})

	_, _ = s.Every(1).Saturday().At("12:00").Do(func() {
		sendWeeklySummary(b)
	})

	// Erinnerungen um 17:00 Uhr für diejenigen, die noch nichts geschrieben haben
	_, _ = s.Every(1).Monday().At("17:00").Do(func() {
		remindUsers(b)
	})
	_, _ = s.Every(1).Tuesday().At("17:00").Do(func() {
		remindUsers(b)
	})
	_, _ = s.Every(1).Wednesday().At("17:00").Do(func() {
		remindUsers(b)
	})
	_, _ = s.Every(1).Thursday().At("17:00").Do(func() {
		remindUsers(b)
	})
	_, _ = s.Every(1).Friday().At("17:00").Do(func() {
		remindUsers(b)
	})

	// Flamen um 22:00 Uhr, wenn jemand noch nicht geschrieben hat
	_, _ = s.Every(1).Monday().At("22:00").Do(func() {
		flameUsers(b)
	})
	_, _ = s.Every(1).Tuesday().At("22:00").Do(func() {
		flameUsers(b)
	})
	_, _ = s.Every(1).Wednesday().At("22:00").Do(func() {
		flameUsers(b)
	})
	_, _ = s.Every(1).Thursday().At("22:00").Do(func() {
		flameUsers(b)
	})
	_, _ = s.Every(1).Friday().At("22:00").Do(func() {
		flameUsers(b)
	})

	s.StartAsync()
}

// Funktion zum Auslesen der Chat-IDs, in denen der Bot Mitglied ist
func getChatIDs() []int64 {
	var ids []int64
	for chatID := range chatIDs {
		ids = append(ids, chatID)
	}
	return ids
}

// Funktion zum Laden der Chat-IDs aus einer Datei
func loadChatIDs() {
	file, err := os.Open("chat-ids.txt")
	if err != nil {
		if os.IsNotExist(err) {
			// Datei existiert nicht, also keine IDs zu laden
			return
		}
		log.Fatalf("Fehler beim Öffnen der Datei chat-ids.txt: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if chatID, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64); err == nil {
			chatIDs[chatID] = struct{}{}
		} else {
			log.Printf("Fehler beim Parsen der Chat-ID aus der Zeile %q: %v", line, err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Fehler beim Lesen der Datei chat-ids.txt: %v", err)
	}
}

// Funktion zum Speichern der Chat-IDs in eine Datei
func saveChatIDs() {
	log.Println("saveChatIDs")

	file, err := os.Create("chat-ids.txt")
	if err != nil {
		log.Fatalf("Fehler beim Erstellen der Datei chat-ids.txt: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for chatID := range chatIDs {
		_, err := writer.WriteString(fmt.Sprintf("%d\n", chatID))
		if err != nil {
			log.Printf("Fehler beim Schreiben der Chat-ID %d: %v", chatID, err)
		}
	}
	writer.Flush()
}

func sendWelcomeMessage(b *telebot.Bot, chat *telebot.Chat) {
	message := "Willkommen im Daily StandUp! 🎉 Hier ist, wie wir arbeiten:\n1. Täglich um 12:00 Uhr bekommst du eine Nachricht mit den Fragen für unser StandUp.\n2. Bitte beantworte die Fragen bis 17:00 Uhr, damit wir wissen, was du gemacht hast und ob es Hindernisse gibt.\n3. Wenn du bis 22:00 Uhr nicht geantwortet hast, werde ich dich daran erinnern.\n\nLass uns gemeinsam einen produktiven Tag haben! 🚀"
	log.Printf("Send Welcome Message to %s", chat.Title)

	_, err := b.Send(chat, message)
	if err != nil {
		log.Printf("Fehler beim Senden der Willkommensnachricht an Chat %d: %v", chat.ID, err)
	}
}

// Funktion zum Senden der täglichen StandUp-Nachricht
func sendDailyStandUp(b *telebot.Bot) {
	message := fmt.Sprintf("**Daily StandUp %s %s**\n1. Was hast du gestern gemacht?\n2. Was machst du heute?\n3. Gibt es Unklarheiten oder Hindernisse?", time.Now().Weekday(), time.Now().Format("02.01.2006"))

	// Sende die Nachricht an die entsprechenden Channels
	// Beispiel: Nutze hier eine Liste von Chat-IDs, in die die Nachricht gesendet werden soll
	chatIDs := getChatIDs() // Ersetze durch die realen Chat-IDs
	for _, chatID := range chatIDs {
		_, err := b.Send(&telebot.Chat{ID: chatID}, message)
		if err != nil {
			log.Printf("Fehler beim Senden der Nachricht an Chat %d: %v", chatID, err)
		}
	}
}

// Funktion, die Benutzer erinnert, wenn sie bis 17:00 Uhr noch nichts geschrieben haben
func remindUsers(b *telebot.Bot) {
	ids := getChatIDs()

	for _, chatID := range ids {
		chat := &telebot.Chat{ID: chatID}
		allWritten := checkAllUsersWritten(b, chat)
		if !allWritten {
			// Falls noch nicht alle geantwortet haben, sende eine Erinnerung
			/*message := "Erinnerung: Bitte beantworte die Fragen für den Daily StandUp!"
			_, err := b.Send(chat, message)
			if err != nil {
				log.Printf("Fehler beim Senden der Erinnerung an Chat %d: %v", chatID, err)
			}*/

			for _, activity := range userActivities {
				today := time.Now().Weekday()
				if today >= time.Monday && today <= time.Friday && !activity.Days[today-time.Monday] {
					// Erinnere den Benutzer
					message := fmt.Sprintf("@%s, bitte beantworte die StandUp-Fragen!", activity.UserName)
					for chatID := range userActivities {
						chat := &telebot.Chat{ID: chatID}
						_, err := b.Send(chat, message)
						if err != nil {
							log.Printf("Fehler beim Senden der Erinnerung an Chat %d: %v", chatID, err)
						}
					}
				}
			}
		}
	}

	/*for _, activity := range userActivities {
		today := time.Now().Weekday()
		if today >= time.Monday && today <= time.Friday && !activity.Days[today-time.Monday] {
			// Erinnere den Benutzer
			message := fmt.Sprintf("@%s, bitte beantworte die StandUp-Fragen!", activity.UserName)
			for chatID := range userActivities {
				chat := &telebot.Chat{ID: chatID}
				_, err := b.Send(chat, message)
				if err != nil {
					log.Printf("Fehler beim Senden der Erinnerung an Chat %d: %v", chatID, err)
				}
			}
		}
	}*/
}

// Funktion, die prüft, ob alle Benutzer geschrieben haben
func checkAllUsersWritten(b *telebot.Bot, chat *telebot.Chat) bool {
	allWritten := true
	for _, activity := range userActivities {
		today := time.Now().Weekday()
		if today >= time.Monday && today <= time.Friday && !activity.Days[today-time.Monday] {
			allWritten = false
			break
		}
	}
	if allWritten {
		message := motivationalQuotes[time.Now().Day()%len(motivationalQuotes)]
		_, err := b.Send(chat, message)
		if err != nil {
			log.Printf("Fehler beim Senden der Motivationsnachricht: %v", err)
		}
	}
	return allWritten
}

// Funktion zum Flamen der Benutzer, die bis 22:00 Uhr noch nicht geschrieben haben
func flameUsers(b *telebot.Bot) {
	ids := getChatIDs()
	for _, chatID := range ids {
		chat := &telebot.Chat{ID: chatID}
		allWritten := checkAllUsersWritten(b, chat)
		if !allWritten {
			for _, activity := range userActivities {
				today := time.Now().Weekday()
				if today >= time.Monday && today <= time.Friday && !activity.Days[today-time.Monday] {
					message := fmt.Sprintf("@%s, ich bin enttäuscht, dass du noch nicht geantwortet hast! 😠", activity.UserName)
					_, err := b.Send(chat, message)
					if err != nil {
						log.Printf("Fehler beim Senden der Flame-Nachricht an Chat %d: %v", chatID, err)
					}
				}
			}
		}
	}
}

// Funktion zum Aufzeichnen der Benutzeraktivität
func RecordActivity(userID int64, userName string, now time.Time) {
	if now.Weekday() >= time.Monday && now.Weekday() <= time.Friday {
		if activity, exists := userActivities[userID]; exists {
			activity.Days[now.Weekday()-time.Monday] = true
		} else {
			var days [5]bool
			days[now.Weekday()-time.Monday] = true
			userActivities[userID] = &UserActivity{
				UserID:   userID,
				UserName: userName,
				Days:     days,
			}
		}
	}
}

// Funktion zum Senden der Wochenstatistik
func sendWeeklySummary(b *telebot.Bot) {
	summary := GenerateWeeklySummary()

	// Sende die Nachricht an die entsprechenden Channels
	/// chatIDs := []int64{-1234567890, -987654321} // Ersetze durch die realen Chat-IDs

	chatIDs := getChatIDs()

	for _, chatID := range chatIDs {
		_, err := b.Send(&telebot.Chat{ID: chatID}, summary)
		if err != nil {
			log.Printf("Fehler beim Senden der Wochenstatistik an Chat %d: %v", chatID, err)
		}
	}
}

// Funktion zur Generierung der Wochenstatistik
func GenerateWeeklySummary() string {
	var sb strings.Builder
	sb.WriteString("**Wochenstatistik** 📊\n")

	for _, activity := range userActivities {
		sb.WriteString(fmt.Sprintf("@%s: ", activity.UserName))
		daysActive := 0
		for _, active := range activity.Days {
			if active {
				daysActive++
			}
		}

		switch daysActive {
		case 1, 2:
			sb.WriteString(fmt.Sprintf("%d Tage 👎\n", daysActive))
		case 3:
			sb.WriteString(fmt.Sprintf("%d Tage 🙏\n", daysActive))
		case 4:
			sb.WriteString(fmt.Sprintf("%d Tage ✈️🔥\n", daysActive))
		case 5:
			sb.WriteString(fmt.Sprintf("%d Tage 🚀🤩\n", daysActive))
		}
	}

	return sb.String()
}
