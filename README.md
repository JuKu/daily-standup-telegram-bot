# Telegram Daily StandUp Bot

Ein Telegram-Bot, der tägliche StandUp-Erinnerungen in Telegram-Kanälen postet. Der Bot erinnert Teammitglieder jeden Tag (außer am Wochenende) um 12:00 Uhr an ihr Daily StandUp und verfolgt, ob alle ihre Updates gegeben haben.

## Funktionen

- **Tägliche Erinnerung:** Der Bot postet jeden Wochentag um 12:00 Uhr eine StandUp-Erinnerung in den Kanälen, in denen er hinzugefügt wurde.
- **Nachfassen:** Um 17:00 Uhr erinnert der Bot automatisch alle Teammitglieder, die ihre StandUp-Fragen noch nicht beantwortet haben.
- **Motivation:** Sobald alle Teammitglieder ihre Antworten gegeben haben, postet der Bot einen zufällig ausgewählten Motivationsspruch.
- **Ermahnung:** Um 22:00 Uhr erinnert der Bot alle, die noch nicht geantwortet haben, mit einer freundlichen, aber bestimmten Nachricht.

## Anforderungen

- **Go:** Dieses Projekt ist in Go geschrieben. Du benötigst [Go](https://golang.org/dl/) in Version 1.16 oder höher, um es zu bauen.
- **Docker:** Um den Bot in einem Docker-Container auszuführen, benötigst du [Docker](https://www.docker.com/products/docker-desktop).
- **Telegram Bot API Key:** Du benötigst einen Telegram Bot API Key, den du über den BotFather auf Telegram erhalten kannst.

## Installation und Ausführung

### 1. Bot einrichten

1. Klone dieses Repository:

    ```bash
    git clone https://github.com/dein-benutzername/telegram-bot.git
    cd telegram-bot
    ```

2. Erstelle eine `.env`-Datei im Projektverzeichnis und füge deinen Telegram Bot API Key hinzu:

    ```plaintext
    TELEGRAM_BOT_API_KEY=dein_api_key
    ```

### 2. Mit Go ausführen

1. Installiere die Abhängigkeiten:

    ```bash
    go mod tidy
    ```

2. Baue und starte den Bot:

    ```bash
    go build -o telegram-bot
    ./telegram-bot
    ```

### 3. Mit Docker ausführen

1. Erstelle ein Docker-Image:

    ```bash
    docker build -t telegram-bot .
    ```

2. Führe den Bot mit Docker Compose aus:

    ```bash
    docker-compose up -d
    ```

3. Überwache die Logs:

    ```bash
    docker-compose logs -f
    ```

### 4. Mit GitLab CI/CD

1. Richte GitLab CI/CD ein, um das Projekt automatisch zu bauen und zu deployen. Verwende die bereitgestellte `.gitlab-ci.yml`.

## Verzeichnisstruktur

```plaintext
telegram-bot/
├── Dockerfile
├── docker-compose.yml
├── .env (optional)
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Lizenz

Dieses Projekt ist unter der MIT-Lizenz lizenziert. Siehe die [LICENSE](LICENSE.md) Datei für weitere Details.

## Mitwirkende

  - @JuKu
  - ChatGPT

## Kontakt

Bei Fragen oder Problemen, zögere nicht, ein Issue zu erstellen oder mich direkt zu kontaktieren.