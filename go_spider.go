package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

var aUrl = make(map[string]string)
var urls = make(map[string]string)
var number int

var countPage = 5
func main() {

	for i := countPage; i > 0; i-- {
		getPage(i)
	}
	generatePdf()
}

func getPage(page int) {
	url := fmt.Sprintf("https://www.jianshu.com/u/ae840e18f653?order_by=shared_at&page=%d", page)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		fmt.Println(resp.StatusCode)
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024)
	bufString := ""
	for {
		n, _ := resp.Body.Read(buf)
		if 0 == n {
			break
		}
		bufString += string(buf[:n])
	}

    hrefRegexp := regexp.MustCompile("<a.*>(.*)</a>")
	match := hrefRegexp.FindAllString(bufString, -1)

	if match != nil {
		for _, v := range match {

            hrefRegexpHref := regexp.MustCompile("/p/\\w*")
            matchHref := hrefRegexpHref.FindAllString(v, -1)

            hrefRegexpName := regexp.MustCompile(">(.*)<")
            matchName := hrefRegexpName.FindAllString(v, 1)

            hrefUrl := ""
            hrefName := ""
            if len(matchHref) > 0 {
                for _, href := range matchHref {
                    hrefUrl = href
                }

                for _, name := range matchName {
                    name = strings.Replace(name, ">", "", -1)
                    name = strings.Replace(name, "<", "", -1)
                    name = strings.Replace(name, " ", "", -1)
                    name = strings.Replace(name, "、", "_", -1)
                    name = strings.Replace(name, "：", ":", -1)
                    hrefName = name
                }
                _, ok := aUrl[hrefUrl]
                if ok {
                    continue
                }
                aUrl[hrefUrl] = hrefUrl
                urls[hrefName] = "https://www.jianshu.com" + hrefUrl
            }
		}
	}
}

func generatePdf() {
	for key, url := range urls {
		fileName := fmt.Sprintf("%s.pdf", key)
		fmt.Println(fileName)
		cmd := exec.Command("wkhtmltopdf", url, fileName)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Execute Command failed:" + err.Error())
			return
		}
		fmt.Println("Execute Command finished.")
	}

}
