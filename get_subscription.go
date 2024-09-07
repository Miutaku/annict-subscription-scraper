package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	animeID := r.URL.Query().Get("id")
	if animeID == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://annict.com/works/%s", animeID)
	log.Printf("Fetching page: %s", url)

	res, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch page", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		http.Error(w, fmt.Sprintf("Failed to fetch page: %s", res.Status), http.StatusInternalServerError)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		http.Error(w, "Failed to parse page", http.StatusInternalServerError)
		return
	}

	// Fetch the list of streaming services
	services, err := fetchStreamingServices()
	if err != nil {
		http.Error(w, "Failed to fetch streaming services", http.StatusInternalServerError)
		return
	}

	log.Printf("Fetched Services: %+v", services)

	// Check each service's availability
	var serviceResponses []ServiceResponse
	for _, service := range services {
		found := doc.Find(fmt.Sprintf("a:contains('%s')", service)).Length() > 0
		serviceResponses = append(serviceResponses, ServiceResponse{Name: service, Available: found})
	}

	response := Response{Services: serviceResponses}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func fetchStreamingServices() ([]string, error) {
	// Create a collector
	c := colly.NewCollector()

	// Define callback for visiting the channel list page
	channels := []map[string]string{}
	c.OnHTML("table", func(e *colly.HTMLElement) {
		// Find all rows within the table (excluding the header row)
		e.ForEach("tbody tr", func(_ int, el *colly.HTMLElement) {
			channel := map[string]string{}

			// Extract data from each table cell
			el.ForEach("td", func(_ int, td *colly.HTMLElement) {
				text := strings.TrimSpace(td.Text)
				switch td.Index {
				case 0:
					channel["ID"] = text
				case 1:
					channel["名前"] = text
				case 2:
					channel["チャンネルグループ"] = text
				case 3:
					// SVG要素を検索し、data-icon属性の値をチェック
					if td.Index == 3 {
						text := strings.TrimSpace(td.Text)
						channel["(Annictがサポートしている) 動画サービス"] = "" // 初期化
						// "-"の場合、そのまま"-"をセット
						if text == "-" {
							channel["(Annictがサポートしている) 動画サービス"] = "-"
						} else {
							channel["(Annictがサポートしている) 動画サービス"] = "○"
						}
					}
				case 4:
					channel["ソート番号"] = text
				case 5:
					channel["状態"] = text
				}
			})

			// Append the extracted channel data to the channels slice
			channels = append(channels, channel)
		})
	})

	// Visit the channel list page
	url := "https://annict.com/db/channels"
	err := c.Visit(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var supportedServices []string
	for _, channel := range channels {
		if channel["(Annictがサポートしている) 動画サービス"] == "○" {
			supportedServices = append(supportedServices, channel["名前"])
		}
	}

	return supportedServices, nil
}
