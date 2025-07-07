package main

import (
	"bountyboard/cmd/server/servrun"
	"bountyboard/internal/domain/task"
	"encoding/gob"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
		os.Exit(1)
	}
	gob.Register(&task.Task{})
}

func main() {
	if err := servrun.Run(); err != nil {
		log.Fatal(err)
	}
}
