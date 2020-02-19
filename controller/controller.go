package controller

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"os"
	"fmt"
	"io/ioutil"
	"io"
	"log"
	"crypto/md5"
	"net/http"
	"go-crapy/message"
)

type Worker struct {
	Cron       *cron.Cron
	CronSpec   string
	HashString string
	Url        string
	Name       string
	HashFile   string
}

func StartWorker(cronSpec, url, workerName, hashFile string) {
	worker := &Worker{
		Cron:     cron.New(),
		CronSpec: cronSpec,
		Url:      url,
		Name:     workerName,
		HashFile: hashFile,
	}
	oldValue, err := worker.LoadHashFromFile()
	if err != nil {
		panic(err)
	}
	worker.HashString = oldValue
	worker.Cron.AddFunc(cronSpec, func() {
		worker.Monitor()
	})
	worker.Cron.Start()
	defer worker.Cron.Stop()
	select {}
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
		w.SaveHashToFile()
		return
	}

	if w.HashString != result {
		log.Print(w.Name + "diff!!!")
		w.HashString = result
		w.SaveHashToFile()
		message.SendMsg(w.Name)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	log.Print(w.Name + " same")
}

func (w *Worker) SaveHashToFile() error {
	var err error
	var f *os.File

	if CheckFileIsExist(w.HashFile) { //如果文件存在
		f, err = os.OpenFile(w.HashFile, os.O_TRUNC|os.O_WRONLY, 0666) //打开文件
		//fmt.Println("文件存在");
	} else {
		f, err = os.Create(w.HashFile) //创建文件
		//fmt.Println("文件不存在");
	}
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, w.HashString) //写入文件(字符串)
	if err != nil {
		return err
	}
	//fmt.Printf("写入 %d 个字节\n", n);
	return nil
}

func (w *Worker) LoadHashFromFile() (string, error) {
	var err error
	var f *os.File
	if CheckFileIsExist(w.HashFile) { //如果文件存在
		f, err = os.OpenFile(w.HashFile, os.O_RDONLY, 0666) //打开文件
		//fmt.Println("文件存在");
	} else {
		return "", nil
	}
	defer f.Close()
	if err != nil {
		return "", fmt.Errorf("Fail to open file ")
	}
	contents, err := ioutil.ReadAll(f)
	if err != nil {

		return "", fmt.Errorf("Fail to read file ")
	}
	return string(contents), nil
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
