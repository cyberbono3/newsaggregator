package main

import (
	"net/http"
	"fmt"
	"html/template"
	"io/ioutil"
	"encoding/xml"
	"sync"

	)

	var wg sync.WaitGroup


	type SitemapIndex struct {
		Locations []string `xml:"sitemap>loc"`
		}
	
	type News struct {
		 Titles []string `xml:"url>news>title"`
		 Keywords []string `xml:"url>news>keywords"`
		 Locations []string `xml:"url>loc"`
	}
	
	type NewsMap struct {
		Keyword string
		Location string
	}
	
	

type NewsAggPage struct {
	Title string
	News map[string]NewsMap
}

func index_handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"whoa, go is neat")
}

func newsRoutine(c chan News, Location string) {
    defer wg.Done()
	var n News
	resp, _ := http.Get(Location)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &n)
	resp.Body.Close()  //free resourcers
	c <- n
}





func newsAggHandler(w http.ResponseWriter, r *http.Request){
	// see xmlpasingcont

    
	var s SitemapIndex

	news_map := make(map[string]NewsMap)


	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()  //free resourcers
    xml.Unmarshal(bytes, &s)

	
	queue := make(chan News, 30)
	for _, Location := range s.Locations {
		wg.Add(1)
	    go newsRoutine(queue, Location)
	}

	wg.Wait()
	close(queue)
	// iterate through news items received from the channel
    for item := range(queue) {
		for idx,_ := range item.Keywords {
			news_map[item.Titles[idx]] = NewsMap{item.Keywords[idx], item.Locations[idx] }
		}
	}
	
	p := NewsAggPage{Title: "Amazing news aggregator", News: news_map}
	t, err := template.ParseFiles("newsaggtemplate.html")
	if (err!=nil) {
		fmt.Println(err)
		}
	fmt.Println(t.Execute(w,p))
}

func main() {
	http.HandleFunc("/", index_handler)
	http.HandleFunc("/agg/", newsAggHandler)
	http.ListenAndServe(":8000", nil)
}