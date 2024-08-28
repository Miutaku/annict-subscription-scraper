package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Response structure for each service
type ServiceResponse struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
}

// Full response structure
type Response struct {
	Services []ServiceResponse `json:"services"`
}

func main() {
	http.HandleFunc("/", handleRequest)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	animeID := r.URL.Query().Get("id")
	if animeID == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://annict.com/works/%s/info", animeID)

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
	var services []string
	res, err := http.Get("https://annict.com/db/channels")
	if err != nil {
		return services, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return services, fmt.Errorf("failed to fetch channels: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return services, err
	}

	// Extract streaming services from the page
	doc.Find(".db-channel-group").Each(func(i int, s *goquery.Selection) {
		group := s.Find(".db-channel-group__name").Text()
		if strings.Contains(group, "動画配信サービス") {
			s.Find(".db-channel-group__channels a").Each(func(j int, c *goquery.Selection) {
				service := strings.TrimSpace(c.Text())
				if service != "" {
					services = append(services, service)
				}
			})
		}
	})

	return services, nil
}
