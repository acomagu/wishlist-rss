package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Listening on port %s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func handler() http.Handler {
	router := httprouter.New()
	router.GET("/:id", func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if err := func() error {
			resp, err := hitWishLister(params.ByName("id"))
			if err != nil {
				return err
			}

			var items []*struct {
				Name      string `json:"name",xml:"name"`
				Link      string `json:"link"`
				DateAdded string `json:"date-added"`
				date      time.Time
			}
			if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
				return err
			}

			for _, item := range items {
				item.date, err = parseDate(item.DateAdded)
				if err != nil {
					return err
				}
			}

			var latestDate time.Time
			for _, item := range items {
				if latestDate.Before(item.date) {
					latestDate = item.date
				}
			}

			fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
			fmt.Fprintf(w, "<rss version=\"2.0\"\n")
			fmt.Fprintf(w, "	xmlns:dc=\"http://purl.org/dc/elements/1.1/\"\n")
			fmt.Fprintf(w, "	xmlns:sy=\"http://purl.org/rss/1.0/modules/syndication/\"\n")
			fmt.Fprintf(w, "	xmlns:admin=\"http://webns.net/mvcb/\"\n")
			fmt.Fprintf(w, "	xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\">\n")
			fmt.Fprintf(w, "	<channel>\n")
			fmt.Fprintf(w, "		<title>Wishlist</title>\n")
			fmt.Fprintf(w, "		<dc:language>ja</dc:language>\n")
			fmt.Fprintf(w, "		<dc:date>%s</dc:date>\n", latestDate.Format(time.RFC822))
			for _, item := range items {
				name, link := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
				xml.EscapeText(name, []byte(item.Name))
				xml.EscapeText(link, []byte(item.Link))

				fmt.Fprintf(w, "		<item>\n")
				fmt.Fprintf(w, "			<title>%s</title>\n", name)
				fmt.Fprintf(w, "			<link>%s</link>\n", link)
				fmt.Fprintf(w, "			<dc:date>%s</dc:date>\n", item.date.Format(time.RFC822))
				fmt.Fprintf(w, "		</item>\n")
			}
			fmt.Fprintf(w, "	</channel>\n")
			fmt.Fprintf(w, "</rss>\n")

			return nil
		}(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	return router
}

func parseDate(str string) (time.Time, error) {
	return time.ParseInLocation("2006年1月2日に追加された商品", str, time.FixedZone("JST", 9*60*60))
}

var hitWishLister = func(id string) (*http.Response, error) {
	url := fmt.Sprintf("http://www.justinscarpetti.com/projects/amazon-wish-lister/api/?id=%s&tld=co.jp&format=json", id)
	return http.Get(url)
}
