package main

import (
	"net/http"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

var webhook = os.Getenv("WEBHOOK_URL")
var pattern = regexp.MustCompile(`^grouped-item  product-purchase-wrapper-404(69|75)`)

// Alert does a GET request to the webhook URL
func Alert() {
	response, err := http.Get(webhook)
	if err != nil {
		log.WithField("err", err).Fatal()
	}
	log.WithField("status", response.StatusCode).Warn("Alerted via Telegram")
}

// IsAvailable checks if the item is available
func IsAvailable(i int, s *goquery.Selection) bool {
	class, exists := s.Attr("class")
	itemName := s.Find(".item-name").Text()
	if exists && pattern.MatchString(class) {
		log.WithField("item", itemName).Warn("Product is available, exiting loop")
		Alert()
		return false
	}
	return true
}

func init() {
	if !strings.HasPrefix(webhook, "https://") {
		log.WithField("webhook_url", webhook).Fatal("WEBHOOK_URL environment variable seems uninitialized")
	}
}

// HandleError handles common errors
func HandleError(err error) {
	if err != nil {
		log.WithField("err", err).Fatal()
	}
}

func main() {
	log.SetOutput(os.Stdout)

	c := colly.NewCollector(
		colly.AllowedDomains("rogueeurope.eu", "www.rogueeurope.eu"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36"),
	)
	err := c.SetProxy("socks5://localhost:9050")
	HandleError(err)

	c.OnHTML("div.product-purchase", func(e *colly.HTMLElement) {
		e.DOM.Find(".grouped-item").Has(".grouped-item-row .item-qty").EachWithBreak(IsAvailable)
	})

	c.OnResponse(func(r *colly.Response) {
		log.WithFields(log.Fields{
			"status": r.StatusCode,
			"url":    r.Request.URL,
		}).Info("Checking availability")
	})

	err = c.Visit("https://www.rogueeurope.eu/rogue-calibrated-kg-steel-plates-eu")
	HandleError(err)
}
