package main

import (
	"fmt"
	"godiff/internal/app"
	"godiff/internal/db"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	f, err := os.OpenFile("debug.log", os.O_TRUNC|os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)

	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	log.SetOutput(f)
	log.SetPrefix("debug")
	defer f.Close()

	// reset the log file after every start, just so we don't accidentally put hundreds of GBs of logs while developing
	f.Truncate(0)
	f.Seek(0, 0)

	log.SetReportCaller(true)
	log.SetReportTimestamp(true)
	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)

	log.Info("Starting application...")
	err = db.InitDatabase()
	if err != nil {
		log.Error("Something went wrong while trying to open the database: %v \n", err)
		os.Exit(1)
	}
	log.Debug("Database initialized")

	err = db.MigrateDatabase()
	if err != nil {
		log.Error("Something went wrong while trying to migrate the database: %v \n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(app.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Error("Something went wrong while trying to run the program: %v \n", err)
		os.Exit(1)
	}
}
