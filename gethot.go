package main

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net/http"
	"strings"
	"time"
)

func main(){
	//关闭gin debug
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/gethot", gethot)
	//测试api是否正常访问
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Run(":8010")
}

func gethot(c *gin.Context) {
	var allData []map[string]interface{}
	sort_url := "https://tophub.today/c/developer"
	sort_timeout := time.Duration(15 * time.Second)
	sort_client := &http.Client{
		Timeout: sort_timeout,
	}
	var sort_Body io.Reader
	sort_request, err := http.NewRequest("GET", sort_url, sort_Body)
	if err != nil {
		message = "抓取失败"
		c.String(200, message)
		return
	}
	sort_res, err := sort_client.Do(sort_request)
	if err != nil {
		message = "抓取失败"
		c.String(200, message)
		return
	}
	defer sort_res.Body.Close()
	sort_document, err := goquery.NewDocumentFromReader(sort_res.Body)
	if err != nil {
		message = "抓取失败"
		c.String(200, message)
		return
	}
	sort_document.Find(".bc-cc div").Each(func(i int, selection *goquery.Selection) {
		id, _ := selection.Attr("id")
		typeName := selection.Find(".cc-cd-lb").Text()
		typeName = strings.Replace(typeName, " ", "", -1)
		document.Find("#"+id+" .nano-content a").Each(func(i int, selection *goquery.Selection) {
			url, boolUrl := selection.Attr("href")
			text := selection.Find(".t").Text()
			if boolUrl {
				allData = append(allData, map[string]interface{}{"sort": "开发","type": typeName,"title": text, "url": url})
			}
		})
	})

	od, _ := JSONMarshal(allData,true)
	c.String(200, string(od))
	return
}
//替换转译字符
func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)
	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}
