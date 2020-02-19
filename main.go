package main

import (
	"crypto/md5"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"go-crapy/message"
	"log"
	"net/http"
	"go-crapy/config"
)

type Worker struct {
	Cron       *cron.Cron
	CronSpec   string
	HashString string
	Url        string
	Name       string
}

func (w *Worker) Monitor() {
	// Request the HTML page.
	res, err := http.Get(w.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	allHrefs := ""
	doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("href")
		//linkText := linkTag.Text()
		//fmt.Printf("Link #%d: '%s' - '%s'\n", index, linkText, link)
		allHrefs += link
	})

	// Hash the hrefs
	h := md5.New()
	h.Write([]byte(allHrefs))
	result := fmt.Sprintf("%x\n\n", h.Sum([]byte("")))
	if w.HashString == "" {
		w.HashString = result
		return
	}

	if w.HashString != result {
		log.Print(w.Name + "diff!!!")
		message.SendMsg(w.Name)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	log.Print(w.Name + " same")
}

func StartWorker(cronSpec, url, workerName string) {
	worker := &Worker{
		Cron:     cron.New(),
		CronSpec: cronSpec,
		Url:      url,
		Name:     workerName,
	}
	worker.Cron.AddFunc(cronSpec, func() {
		worker.Monitor()
	})
	worker.Cron.Start()
	defer worker.Cron.Stop()
	select {}

}

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	go StartWorker(conf.CronSpec, "http://81rc.81.cn/index.htm", "军队人才网首页")
	//go StartWorker(conf.CronSpec, "http://81rc.81.cn/Civilianpost/index.htm", "军队人才网文职人员")

	select {}
}
