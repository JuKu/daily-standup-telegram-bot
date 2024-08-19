package main

import (
	"testing"
	"time"
)

// Mock-Daten für Tests
var testUserID int64 = 123456
var testUserName string = "testuser"

func TestRecordActivity(t *testing.T) {
	// Simuliere Montag
	now := time.Date(2024, 8, 19, 12, 0, 0, 0, time.UTC) // 19. August 2024 ist ein Montag

	// Prüfe, ob die Map leer ist
	if len(userActivities) != 0 {
		t.Fatalf("Erwartete leere userActivities Map, aber es sind %d Einträge vorhanden", len(userActivities))
	}

	// Simuliere die Aktivität des Benutzers
	RecordActivity(testUserID, testUserName, now)

	// Überprüfe, ob die Aktivität korrekt gespeichert wurde
	activity, exists := userActivities[testUserID]
	if !exists {
		t.Fatalf("Aktivität für Benutzer %d wurde nicht korrekt gespeichert", testUserID)
	}

	if !activity.Days[0] { // Montag ist der erste Tag (Index 0)
		t.Errorf("Erwartete Aktivität für Montag, aber sie wurde nicht korrekt gespeichert")
	}

	// Simuliere Dienstag
	now = now.Add(24 * time.Hour) // Ein Tag später

	// Simuliere die Aktivität des Benutzers erneut
	RecordActivity(testUserID, testUserName, now)

	// Überprüfe, ob die Aktivität für Dienstag ebenfalls gespeichert wurde
	if !activity.Days[1] { // Dienstag ist der zweite Tag (Index 1)
		t.Errorf("Erwartete Aktivität für Dienstag, aber sie wurde nicht korrekt gespeichert")
	}
}

func TestWeeklySummary(t *testing.T) {
	// Setze Aktivitäten für die Woche
	userActivities[testUserID] = &UserActivity{
		UserID:   testUserID,
		UserName: testUserName,
		Days:     [5]bool{true, true, true, true, true}, // Alle Tage sind wahr (aktiv)
	}

	summary := GenerateWeeklySummary()

	// Überprüfe die erzeugte Zusammenfassung
	expected := "**Wochenstatistik** 📊\n@testuser: 5 Tage 🚀🤩\n"
	if summary != expected {
		t.Errorf("Erwartete Zusammenfassung '%s', aber es wurde '%s' erzeugt", expected, summary)
	}
}

func TestWeeklySummaryPartial(t *testing.T) {
	// Setze Aktivitäten für die Woche
	userActivities[testUserID] = &UserActivity{
		UserID:   testUserID,
		UserName: testUserName,
		Days:     [5]bool{true, true, true, false, false}, // Nur 3 Tage aktiv
	}

	summary := GenerateWeeklySummary()

	// Überprüfe die erzeugte Zusammenfassung
	expected := "**Wochenstatistik** 📊\n@testuser: 3 Tage 🙏\n"
	if summary != expected {
		t.Errorf("Erwartete Zusammenfassung '%s', aber es wurde '%s' erzeugt", expected, summary)
	}
}
