package batdongsan_com_vn

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

type Product struct {
	Title          string
	Link           string
	Price          string
	Area           string
	Location       string
	Image          string
	ProductContent string
}

func Start() {
	fName := "batdongsan_com_vn.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	// Instantiate default collector
	c := colly.NewCollector()

	products := make([]Product, 0)

	c.OnHTML("div .product-item", func(e *colly.HTMLElement) {
		link, _    := e.DOM.Find("a").Attr("href")
		image, _    := e.DOM.Find("img").Attr("src")

		log.Printf(image)
		product := Product{
			Title:          e.ChildText(".product-title"),
			Price:          e.ChildText(".price"),
			Area:           e.ChildText(".area"),
			Link:           "https://batdongsan.com.vn/" + link,
			Location:       e.ChildText(".location"),
			Image:          image,
			ProductContent: e.ChildText(".product-content"),
		}
		products = append(products, product)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://batdongsan.com.vn/nha-dat-ban")

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(products)

	log.Printf("Scraping finished, check file %q for results\n", fName)
}