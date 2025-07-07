package config

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"testing"
)

func init() {
	log.Println("Running init()")

	if err := godotenv.Load("../../.env"); err != nil {
		log.Print("No .env file found")
	}
}

func TestNewPgPool(t *testing.T) {
	ctx := context.Background()

	pool, err := NewPgPool(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pool)
	defer pool.Close()
}
