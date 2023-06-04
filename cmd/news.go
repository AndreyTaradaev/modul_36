package main

import (
	"fmt"
	"log"
	"news/pkg/channel"
	"news/pkg/db"
	"news/pkg/api"
	"os"
	"net/http"
)

func main() {
	//Connect DB
	DB, err := db.Create()
	if err != nil {
		log.Fatal(err)
	}
	
	defer DB.Close()

	api := api.New(DB)
	//load setting
	set, err := DB.LoadSetting()
	if err != nil {
		log.Fatal(err)
	}
	connstr := os.Getenv("Debug")
	if connstr != "" {
		fmt.Println(set) 
	}

	
	var chexit = make(chan bool)
	var cherr = make(chan error)
	var chNews = make(chan []db.New)
	if (set.Refresh==0){
		set.Refresh =300
	}
	go channel.GetNews(set.Urls,set.Refresh,chexit,cherr,chNews)
	go channel.LoadNews(chexit,cherr,chNews)
	channel.ReadConsole(chexit, cherr)
	go channel.WriteConsole(chexit, cherr)
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}


	
}
