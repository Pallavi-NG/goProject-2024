package main

import (
	"goProject-2024/database"
	transportHTTP "goProject-2024/httpTransport"
	"goProject-2024/models"

	log "github.com/sirupsen/logrus"
)

// Run - sets up our application
func Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Setting Up Our APP")

	var err error
	store, err := database.NewDatabase()
	if err != nil {
		log.Error("failed to setup connection to the database")
		return err
	}

	studentService := models.NewService(store)
	handler := transportHTTP.NewHandler(studentService)

	if err := handler.Serve(); err != nil {
		log.Error("failed to gracefully serve our application")
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Error(err)
		log.Fatal("Error starting up our REST API")
	}
}
