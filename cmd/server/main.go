package main

import (
	"bountyboard/cmd/server/servrun"
	"bountyboard/internal/domain/task"
	"encoding/gob"
	"log"
)

func init() {
	gob.Register(&task.Task{})
}

func main() {
	if err := servrun.Run(); err != nil {
		log.Fatal(err)
	}
}
