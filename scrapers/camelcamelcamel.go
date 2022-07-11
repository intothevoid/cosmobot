package scrapers

import (
	"errors"
	"strings"
	"time"

	"github.com/intothevoid/kramerbot/models"
	"github.com/intothevoid/kramerbot/util"
	"go.uber.org/zap"
)

var SID_CCC_AMAZON ScraperID = 1

// Camel Camel Camel - Amazon scraper
type CamCamCamScraper struct {
	BaseUrl         []string               // Urls to scrape
	Logger          *zap.Logger            // Reference to main logger
	SID             ScraperID              // Scraper ID
	ScrapeInterval  int                    // Scrape interval
	MaxDealsToStore int                    // Max. no. of deals to have in memory
	Deals           []models.CamCamCamDeal // List of deals
}

// Check initialisation
func (s *CamCamCamScraper) CheckInit() bool {
	if s.ScrapeInterval == 0 || s.MaxDealsToStore == 0 || len(s.BaseUrl) <= 0 || s.Logger == nil {
		return false
	}
	return true
}

// Scrape the url
func (s *CamCamCamScraper) Scrape() error {
	if !s.CheckInit() {
		return errors.New("Scraper not initialized correctly. Ensure all fields are set")
	}

	for _, currUrl := range s.BaseUrl {
		// Scrape RSS feed
		parser := util.RssParser{
			Url:    currUrl,
			Logger: s.Logger,
		}

		feed, err := parser.ParseFeed()
		if err != nil {
			return err
		}

		// Loop through deals
		for _, deal := range feed.Items {

			// populate the deal
			deal := models.CamCamCamDeal{
				Id:        deal.GUID,
				Title:     deal.Title,
				Url:       deal.Link,
				Published: deal.Published,
				Image:     deal.Image.URL,
				DealType:  feed.Title,
			}
			s.Logger.Debug("Found deal", zap.String("title", deal.Title), zap.String("url", deal.Url), zap.String("time", deal.Published))

			// create item list
			s.Deals = append(s.Deals, deal)
		}
	}

	// Keep deals length under 'MaxDeals'
	if len(s.Deals) > s.MaxDealsToStore {
		s.Deals = s.Deals[len(s.Deals)-s.MaxDealsToStore:]
	}

	return nil
}

// Filter list of deals by keywords
func (s *CamCamCamScraper) FilterByKeywords(keywords []string) []models.CamCamCamDeal {
	filteredDeals := []models.CamCamCamDeal{}
	for _, deal := range s.Deals {
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(deal.Title), strings.ToLower(keyword)) {
				filteredDeals = append(filteredDeals, deal)
			}
		}
	}
	return filteredDeals
}

// Get 'count' deals from the list of deals
func (s *CamCamCamScraper) GetLatestDeals(count int) []models.CamCamCamDeal {
	if len(s.Deals) <= count {
		return s.Deals
	}
	return s.Deals[:count]
}

// go routine to auto scrape every X minutes
func (s *CamCamCamScraper) AutoScrape() {
	// Scrape once before interval
	err := s.Scrape()
	if err != nil {
		s.Logger.Error("Error scraping", zap.Error(err))
	}

	// use timer to run every 'ScrapeInterval' minutes
	t := time.NewTicker(time.Minute * time.Duration(s.ScrapeInterval))
	go func() {
		for range t.C {
			err := s.Scrape()
			if err != nil {
				s.Logger.Error("Error scraping", zap.Error(err))
			}
		}
	}()
}

// Get scraper data
func (s *CamCamCamScraper) GetData() []models.CamCamCamDeal {
	return s.Deals
}
