package main

import (
	"bank_system/util"
	"fmt"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Not able to load config", err)
	}

	db, err := util.ConnectDB(config)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	// 使用原生 SQL 查詢
	var results []map[string]interface{}
	result := db.Raw("SELECT * FROM users").Scan(&results)
	if result.Error != nil {
		log.Fatalf("failed to execute query: %v", result.Error)
	}

	for _, user := range results {
		fmt.Printf("User: %+v\n", user)
	}

}
