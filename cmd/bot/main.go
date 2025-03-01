package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/iteplenky/employee-attendance/database"
	"github.com/iteplenky/employee-attendance/internal/bot"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := database.NewPostgresDB(dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func(db *database.PostgresDB) {
		if err = db.Close(); err != nil {
			log.Fatal("Failed to close database:", err)
		}
	}(db)

	dbAttendanceURL := os.Getenv("DATABASE_ATTENDANCE_URL")
	if dbAttendanceURL == "" {
		log.Fatal("DATABASE_ATTENDANCE_URL environment variable is not set")
	}

	listener, err := database.NewListener(dbAttendanceURL)
	if err != nil {
		log.Fatal("Failed to initialize listener:", err)
	}
	defer listener.Close()

	go listener.StartListening()

	b := bot.NewBot(db)
	updater := ext.NewUpdater(b.Dispatcher, nil)
	if err = bot.StartPolling(b, updater); err != nil {
		log.Fatal("Failed to start polling:", err)
	}

	log.Printf("%s has been started...\n", b.Bot.User.Username)
	updater.Idle()
}
