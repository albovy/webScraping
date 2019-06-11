package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//Esto es una prueba de webScraping

func main() {
	complexUrl := "https://www.tripadvisor.es/Hotel_Review-g187514-d10190864-Reviews-Gran_Hotel_Ingles-Madrid.html"
	splited := strings.Split(complexUrl, "Reviews")

	parsedUrl, err := url.Parse(complexUrl)

	if err != nil {
		log.Fatal(err)
	}
	response, err := http.Get(parsedUrl.String())
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	numPags := document.Find(".pageNum").Last().Text()
	pages, _ := strconv.Atoi(numPags)
	mapa := make(map[string]string)
	for i := 0; i < pages; i++ {
		num := i * 5
		paginacion := splited[0] + "Reviews-or" + strconv.Itoa(num) + splited[1]
		fmt.Println(paginacion)
		nextPage, _ := url.Parse(paginacion)
		nextResponse, _ := http.Get(nextPage.String())
		defer response.Body.Close()
		nextPageDoc, _ := goquery.NewDocumentFromReader(nextResponse.Body)

		nextPageDoc.Find(".hotels-review-list-parts-ReviewTitle__reviewTitleText--3QrTy").Each(func(i int, selection *goquery.Selection) {
			title := selection.Find("span").First().Text()

			hrefe, _ := selection.Attr("href")
			var showCurrentPageUrl url.URL
			showCurrentPageUrl.Scheme = parsedUrl.Scheme
			showCurrentPageUrl.Host = parsedUrl.Host
			showCurrentPageUrl.Path = hrefe
			response, err := http.Get(showCurrentPageUrl.String())
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()
			showCurrentPageDoc, _ := goquery.NewDocumentFromReader(response.Body)
			texto := showCurrentPageDoc.Find(".fullText").Text()
			mapa[title] = texto

		})

	}
	fmt.Println(mapa)

}
