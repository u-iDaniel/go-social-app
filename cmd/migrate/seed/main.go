package main

import (
	"log"

	"github.com/u-iDaniel/go-social-app/internal/db"
	"github.com/u-iDaniel/go-social-app/internal/env"
	"github.com/u-iDaniel/go-social-app/internal/store"
)

func main() {
	// You should truncate all tables (i.e. delete all data but keep the columns) before seeding
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5431/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
