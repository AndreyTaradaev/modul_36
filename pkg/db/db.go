package db

import (
	"context"
	"fmt"

	//"errors"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Db struct {
	pool *pgxpool.Pool
}

// структура Данных из таблицы новостей

type New struct {
	ID          int `json:"ID"`
	Title       string `json:"Title"`
	Description string  `json:"Content"`
	Time        int64  `json:"PubTime"`
	Url         string `json:"Link"`
	Guid        string `json:"Guid"`
}

const (
	cUrls int = iota + 1
	cRefresh
)

type config struct {
	Db     string `json:"database"`
	User   string `json:"user"`
	Passwd string `json:"password"`
	Host   string `json:"host"`
}

var fileconf config

type RssConfig struct {
	Urls    []string
	Refresh int
}

// ctor

func Create() (*Db, error) {
	f, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(f, &fileconf)
	if err != nil {
		log.Fatal(err)
	}
	connstr := "postgres://" + fileconf.User + ":" + fileconf.Passwd + "@" + fileconf.Host + ":5432/" + fileconf.Db
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	db := Db{pool: pool}

	//fmt.Println(connstr)
	// читаем настройки из БД
	return &db, nil
}

func (db *Db) GetActive() {

	rows, err := db.pool.Query(context.Background(), `
	SELECT "ID", "ValueStr"
			 FROM "Settings"
				WHERE  "Active"=true and "TypeSetting" = $1;`, cUrls,
	)
	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var url *string

		var Id *int
		err = rows.Scan(
			&Id,
			&url,
		)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("[%d]  %s\n", *Id, *url)

	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
}

func (db *Db) LoadSetting() (*RssConfig, error) {

	rows, err := db.pool.Query(context.Background(), `
	SELECT  "ValueStr", "TypeSetting", "ValueInt"
			 FROM "Settings"
				WHERE  "Active"=true;`,
	)
	if err != nil {
		return nil, err
	}
	var rss RssConfig
	rss.Urls = make([]string, 0, 10)
	// итерирование по результату выполнения запроса
	defer rows.Close()
	for rows.Next() {
		var url *string
		var time *int
		var typeset *int
		err = rows.Scan(
			&url,
			&typeset,
			&time,
		)
		if err != nil {
			return nil, err
		}
		if *typeset == cUrls {
			rss.Urls = append(rss.Urls, *url)
		} else if *typeset == cRefresh {
			rss.Refresh = *time
		}

	}
	// ВАЖНО не забыть проверить rows.Err()
	return &rss, rows.Err()
}

func (db *Db) AddUrls(u string) (int, error) {
	var id int
	err := db.pool.QueryRow(context.Background(), `
	INSERT INTO "Settings"
	( "Description", "ValueStr", "TypeSetting","Active")
	VALUES( 'URL', $1, $2,  true) RETURNING "ID";`,

		u,
		cUrls,
	).Scan(&id)
	return id, err
}

func (db *Db) ChangeRssStatus(i int, status bool) (int, error) {
	var id int
	err := db.pool.QueryRow(context.Background(), `
	UPDATE "Settings"
SET "Active"=$2
WHERE "ID"=$1  RETURNING "ID";				`,
		i,
		status,
	).Scan(&id)
	return id, err
}

func (db *Db) ChangeDuration(i int) (int, error) {
	var id int
	err := db.pool.QueryRow(context.Background(), `
	UPDATE "Settings"
SET "ValueInt"=$2
WHERE "TypeSetting"=$1  RETURNING "ID";				`,
		cRefresh,
		i,
	).Scan(&id)
	return id, err
}

func (db *Db) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *Db) AddNew(n New) (int, error) {
	var id int
	err := db.pool.QueryRow(context.Background(), `
	INSERT INTO "NewsTable"
    ( "Title", "Description", "Time", "Url","GUID")
    values ($1, $2, $3, $4,$5)
       ON CONFLICT DO nothing RETURNING "Id";`,
		n.Title,
		n.Description,
		n.Time,
		n.Url,
		n.Guid,
	).Scan(&id)
	return id, err
}

func (db *Db) AddNews(n []New) error {
	for _, v := range n {
		id, err := db.AddNew(v)
		if err != nil && id!=0 {
			return err
		}

	}
	return nil
}

func (db *Db) GetNews(n int) ([]New, error) {
	rows, err := db.pool.Query(context.Background(), `
	select nt."Id" ,nt."Title" ,NT."Description" ,nt."Url",NT."Time"  from "NewsTable" nt 
order  by NT."Time" desc 
limit  $1`,
		n)
	if err != nil {
		return nil, err
	}
	var news []New = make([]New, 0)
	defer rows.Close()
	for rows.Next() {
		var rss New
		var id *int
		var title *string
		var desc *string
		var url *string
		var time *int64
		err = rows.Scan(
			&id,
			&title,
			&desc,
			&url,
			&time,
		)
		if err != nil {
			return nil, err
		}
		if id != nil {
			rss.ID = *id
		}
		if title != nil {
			rss.Title = *title
		}
		if desc != nil {
			rss.Description = *desc
		}
		if url != nil {
			rss.Url = *url
		}
		if time != nil {
			rss.Time = *time
		}
		news = append(news, rss)
	}
	// ВАЖНО не забыть проверить rows.Err()
	return news, rows.Err()
}
