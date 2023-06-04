package channel

import (
	//"log"
	"bufio"
	"fmt"
	"news/pkg/db"
	"news/pkg/rss"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ReadConsole(exit chan bool, cherr chan<- error) {
	// run thread read console
	fmt.Println("консоль для обработки команд, \nКоманды:  \n\t add URL - добавление ссылки на ленту")
	fmt.Println("\texit  -  выход\n\t list - список подписок")
	fmt.Println("\tdel URL  -  оключить подписку ")
	fmt.Println("\tchange   <number>  изменить время обновления ленты")
	go func() {
		defer close(exit)
		scanner := bufio.NewScanner(os.Stdin)
		var data string
		for {
			scanner.Scan()
			data = scanner.Text()
			if strings.EqualFold(data, "exit") {
				return
			}

			if strings.EqualFold(data, "list") {
				fmt.Println("Список RSS\n")
				Db, err := db.Create()
				if err == nil {
					Db.GetActive()
					Db.Close()
				} else {
					cherr <- err
				}
				continue
			}

			sp := strings.Split(data, " ")
			if len(sp) != 2 {
				cherr <- fmt.Errorf("неизвестная команда\n")
				continue
			}

			if strings.EqualFold(sp[0], "add") {
				Db, err := db.Create()
				if err != nil {
					cherr <- err
					continue
				}
				id, err := Db.AddUrls(sp[1])
				if err != nil {
					cherr <- err
					Db.Close()
					continue
				}
				Db.Close()
				fmt.Printf(" Добавлена ссылка id [%d]   %s\n", id, sp[1])
				fmt.Println("Требуется перезагрузка программы !")
				continue
			}
			if strings.EqualFold(sp[0], "del") {

				i, err := strconv.Atoi(sp[1])
				if err != nil {
					cherr <- err
					continue
				}
				Db, err := db.Create()
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, err = Db.ChangeRssStatus(i, false)
				if err != nil {
					cherr <- err
					Db.Close()
					continue
				}
				Db.Close()
				fmt.Println(" удалена ссылка " + sp[1] + "\n")
				fmt.Println("Требуется перезагрузка программы !")
				continue
			}
			if strings.EqualFold(sp[0], "change") {
				i, err := strconv.Atoi(sp[1])
				if err != nil {
					cherr <- err
					continue
				}
				Db, err := db.Create()
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, err = Db.ChangeDuration(i)
				if err != nil {
					cherr <- err
					Db.Close()
					continue
				}
				Db.Close()
				fmt.Println(" Изменено время обновления " + sp[1] + "\n")
				fmt.Println("Требуется перезагрузка программы !")
				continue
			}
		}
	}()
}

func WriteConsole(exit <-chan bool, cherr <-chan error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		for {
			select {
			case err := <-cherr:
				Log(err.Error())
			case <-exit:
				defer wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	os.Exit(0)
}

func GetNews(urls []string, D int, exit chan bool, cherr chan<- error, News chan<- []db.New) {
	connstr := os.Getenv("Debug")
	debug := connstr != ""
	var c <-chan time.Time = time.After(time.Duration(D) * time.Second)
	for {
		select {
		case <-c:

			for _, v := range urls {
				news, err := rss.GetNewsRss(v)
				if err != nil {
					cherr <- err
					continue
				}
				News <- news
				if debug {
					cherr <- fmt.Errorf("Load News in %s  count %d", v, len(news))
				}
			}

		case <-exit:
			if debug {
				cherr <- fmt.Errorf("quit GetUrls")
			}
			return
		}
	}

}

func LoadNews(exit chan bool, cherr chan<- error, News <-chan []db.New) {
	connstr := os.Getenv("Debug")
	debug := connstr != ""
	for {
		select {
		case <-exit:
			if debug {
				cherr <- fmt.Errorf("quit LoadNews")
			}
			return
		case d := <-News:
			Db, err := db.Create()
			if err != nil {
				cherr <- err
				continue
			}
			err = Db.AddNews(d)
			if err != nil {
				cherr <- err
				continue
			}
			cherr <- fmt.Errorf("Load News succes")

		}
	}

}
