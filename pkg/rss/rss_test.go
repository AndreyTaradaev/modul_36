package rss

import (
	//"news/pkg/db"
	//"reflect"
	"testing"
)

func TestGetNewsRss(t *testing.T) {
	url :="https://habr.com/ru/rss/hub/go/all/?fl=ru"
	news,err := GetNewsRss(url)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", news)
	t.Logf("%d", len(news))
}

