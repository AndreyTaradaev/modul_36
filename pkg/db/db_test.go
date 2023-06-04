package db

import (
	"math/rand"
	
	"strconv"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	db, err := Create()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
}

func TestDb_GetNews(t *testing.T) {

	rand.Seed(time.Now().UnixNano())
	posts := []New{
		{
			Title: "Test Post",
			Url:   strconv.Itoa(rand.Intn(1_000_000_000)),
			Time:  time.Now().Unix(),
			Guid:  strconv.Itoa(1_000_000_000),
		},
	}
	db, err := Create()
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddNews(posts)
	if err != nil {
		t.Fatal(err)
	}
	news, err := db.GetNews(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", news)
}

func TestDb_LoadSetting(t *testing.T) {
	db, err := Create()
	if err != nil {
		t.Fatal(err)
	}
	id,err := db.AddUrls("eeeeeee")
	if err != nil {
		t.Fatal(err)
	}
	config, err := db.LoadSetting()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", config)
	id,err = db.ChangeRssStatus(id,false)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", config)
}
