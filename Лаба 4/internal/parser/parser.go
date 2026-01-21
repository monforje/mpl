package parser

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"lab4/internal/core"
	"lab4/pkg/cfg"

	"github.com/PuerkitoBio/goquery"
)

func ProcessVesselLink(ctx context.Context, client *http.Client, url string) (core.Ship, error) {
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return core.Ship{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return core.Ship{}, fmt.Errorf("ошибка HTTP запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return core.Ship{}, fmt.Errorf("статус код: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return core.Ship{}, fmt.Errorf("ошибка парсинга HTML: %v", err)
	}

	var vesselLinks []string
	doc.Find(cfg.ShipLinkClass).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, cfg.VesselPath) {
			fullURL := cfg.BaseURL + href
			vesselLinks = append(vesselLinks, fullURL)
		}
	})

	if len(vesselLinks) == 0 {
		return core.Ship{URL: url}, fmt.Errorf("суда не найдены")
	}

	vesselURL := vesselLinks[0]
	if len(vesselLinks) > 1 {
		fmt.Printf("  Найдено %d судов, берем первое: %s\n", len(vesselLinks), vesselURL)
	}

	return GetVesselDetails(ctx, client, vesselURL, url)
}

func GetVesselDetails(ctx context.Context, client *http.Client, vesselURL, originalURL string) (core.Ship, error) {
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, vesselURL, nil)
	if err != nil {
		return core.Ship{}, fmt.Errorf("ошибка создания запроса на страницу судна: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return core.Ship{}, fmt.Errorf("ошибка HTTP запроса на страницу судна: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return core.Ship{}, fmt.Errorf("статус код страницы судна: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return core.Ship{}, fmt.Errorf("ошибка парсинга HTML страницы судна: %v", err)
	}

	name := ""
	doc.Find(cfg.TitleSelector).Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			name = strings.TrimSpace(s.Text())
		}
	})

	imoMMSI := ""
	doc.Find(cfg.IMOMMSISelector).Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Parent().Text(), "IMO / MMSI") {
			imoMMSI = strings.TrimSpace(s.Text())
		}
	})

	imo, mmsi := "", ""
	if imoMMSI != "" {
		parts := strings.Split(imoMMSI, "/")
		if len(parts) >= 2 {
			imo = strings.TrimSpace(parts[0])
			mmsi = strings.TrimSpace(parts[1])
		}
	}

	vesselType := ""
	doc.Find(cfg.TypeSelector).Each(func(i int, s *goquery.Selection) {
		if strings.TrimSpace(s.Text()) == cfg.AISTypeText {
			nextTd := s.Next()
			if nextTd.Length() > 0 {
				vesselType = strings.TrimSpace(nextTd.Text())
			}
		}
	})

	if name == "" || imo == "" || mmsi == "" || vesselType == "" {
		return core.Ship{URL: originalURL}, fmt.Errorf("не все данные найдены: name='%s', imo='%s', mmsi='%s', type='%s'",
			name, imo, mmsi, vesselType)
	}

	return core.Ship{
		Name: name,
		IMO:  imo,
		MMSI: mmsi,
		Type: vesselType,
		URL:  originalURL,
	}, nil
}
