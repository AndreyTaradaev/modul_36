package rss

import (
	 "github.com/SlyMarbo/rss"
	 "news/pkg/db"	 
	 tags "github.com/grokify/html-strip-tags-go" 
	 
)

func GetNewsRss(url string) ([]db.New,error){
	

	feed, err := rss.Fetch(url)
	if(err != nil) {
		return nil,err
	}
	var news []db.New

	for _, v := range feed.Items {
		var n db.New 
		n.Title = v.Title
		n.Description = tags.StripTags(v.Summary)
		n.Url = v.Link
		n.Time = v.Date.Unix()
		n.Guid = v.ID
		news = append(news,n )		
	}
	return news,nil
}