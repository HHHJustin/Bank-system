package main

import (
	"bank_system/api"
	"bank_system/util"
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

	server, err := api.NewServer(db, config)
	if err != nil {
		log.Fatal("cannot create server : ", err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
