package batdongsan_com_vn

import (
	"fmt"
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

func processSinglePage(url string) []Product {

	fmt.Println("url", url)
	// Instantiate default collector
	c := colly.NewCollector()

	products := make([]Product, 0)

	c.OnHTML("div .product-item", func(e *colly.HTMLElement) {
		link, _    := e.DOM.Find("a").Attr("href")
		image, _    := e.DOM.Find("img").Attr("src")

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

	c.Visit(url)
	return products
}

func Start() []Product {

	products := make([]Product, 0)

	url := "https://batdongsan.com.vn/nha-dat-ban"
	products = processSinglePage(url)

	for i := 0; i < 10; i++ {
		productsSub := processSinglePage(fmt.Sprintf("%s%s%d", url, "/p", i+1))
		products = append(products, productsSub ...)
	}
	return products
}