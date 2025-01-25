package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Abiti0233/memo-app/app"
	"github.com/Abiti0233/memo-app/configs"

	_ "github.com/lib/pq"
)

func main() {
	config := configs.LoadConfig()

	// データベース接続
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// DBのPingで接続確認
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// アプリケーションの初期化と起動
	a := app.App{}
	a.Initialize(db, config)
	a.Run(fmt.Sprintf(":%s", config.ServerPort))
}
