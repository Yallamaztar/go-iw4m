package player

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Yallamaztar/go-iw4m/models"
	"github.com/Yallamaztar/go-iw4m/server"
	"github.com/Yallamaztar/go-iw4m/wrapper"
)

type Player struct {
	Wrapper *wrapper.IW4MWrapper
}

// Constructor to create Player from IW4MWrapper instance
func NewPlayer(w *wrapper.IW4MWrapper) *Player {
	return &Player{Wrapper: w}
}

func (p *Player) PlayerStats(clientID string) (string, error) {
	r := p.Wrapper.DoRequest(fmt.Sprintf("%s/api/stats/%s", p.Wrapper.BaseURL, clientID))

	if r == "" {
		return "", fmt.Errorf("empty response from server")
	}

	return r, nil
}

func (p *Player) AdvancedStats(clientID string) (*models.AdvancedStats, error) {
	r := p.Wrapper.DoRequest(fmt.Sprintf("%s/clientstatistics/%s/advanced", p.Wrapper.BaseURL, clientID))

	if r == "" {
		return nil, fmt.Errorf("empty response from server")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	model := &models.AdvancedStats{
		HitLocations: make(map[string][]models.HitLocation),
		WeaponUsages: make(map[string][]models.WeaponUsage),
	}
	topCard := doc.Find("div.align-self-center.d-flex.flex-column.flex-lg-row.flex-fill.mb-15")
	if topCard.Length() > 0 {
		name := topCard.Find("a.no-decoration").Text()
		href, _ := topCard.Find("a.no-decoration").Attr("href")
		iconURL, _ := topCard.Find("img.img-fluid.align-self-center.w-75").Attr("src")
		summary := topCard.Find("div#client_stats_summary").Text()

		model.Name = strings.TrimSpace(name)
		model.Link = strings.TrimSpace(href)
		model.IconURL = fmt.Sprintf("%s%s", p.Wrapper.BaseURL, strings.TrimSpace(iconURL))
		model.Summary = strings.TrimSpace(summary)
	}

	mainCard := doc.Find("div.flex-fill.flex-xl-grow-1")
	mainCard.Find("div.stat-card").Each(func(i int, stat *goquery.Selection) {
		key := strings.TrimSpace(stat.Find("div.font-size-12.text-muted").Text())
		value := strings.TrimSpace(stat.Find("div.m-0.font-size-16.text-primary").Text())
		if key != "" && value != "" {
			model.PlayerStats = append(model.PlayerStats, models.StatEntry{Key: key, Value: value})
		}
	})

	bottom := doc.Find("div.d-flex.flex-wrap.flex-column-reverse.flex-xl-row")

	bottom.Find("div.mr-0.mr-xl-20.flex-fill.flex-xl-grow-1").Each(func(i int, hit *goquery.Selection) {
		title := strings.TrimSpace(hit.Find("h4.colorcode").Text())
		var entries []models.HitLocation

		hit.Find("tbody tr.bg-dark-dm.bg-light-lm.d-none.d-lg-table-row").Each(func(j int, row *goquery.Selection) {
			spans := row.Find("span")
			if spans.Length() >= 4 {
				entries = append(entries, models.HitLocation{
					Location:   strings.TrimSpace(spans.Eq(0).Text()),
					Hits:       strings.TrimSpace(spans.Eq(1).Text()),
					Percentage: strings.TrimSpace(spans.Eq(2).Text()),
					Damage:     strings.TrimSpace(spans.Eq(3).Text()),
				})
			}
		})
		if len(entries) > 0 {
			model.HitLocations[title] = entries
		}
	})

	bottom.Find("div.flex-fill.flex-xl-grow-1").Each(func(i int, usage *goquery.Selection) {
		title := strings.TrimSpace(usage.Find("h4.colorcode").Text())
		var weapons []models.WeaponUsage

		usage.Find("tbody tr.bg-dark-dm.bg-light-lm.d-none.d-lg-table-row").Each(
			func(j int, row *goquery.Selection) {
				spans := row.Find("span")
				if spans.Length() >= 6 {
					weapons = append(weapons, models.WeaponUsage{
						Weapon:              strings.TrimSpace(spans.Eq(0).Text()),
						FavoriteAttachments: strings.TrimSpace(spans.Eq(1).Text()),
						Kills:               strings.TrimSpace(spans.Eq(2).Text()),
						Hits:                strings.TrimSpace(spans.Eq(3).Text()),
						Damage:              strings.TrimSpace(spans.Eq(4).Text()),
						Usage:               strings.TrimSpace(spans.Eq(5).Text()),
					})
				}
			})
		if len(weapons) > 0 {
			model.WeaponUsages[title] = weapons
		}
	})

	return model, nil
}

func (p *Player) ClientInfo(clientID string) (map[string]interface{}, error) {
	r := p.Wrapper.DoRequest(fmt.Sprintf("%s/api/client/%s", p.Wrapper.BaseURL, clientID))
	if r == "" {
		return nil, fmt.Errorf("empty response from server")
	}

	var data map[string]any
	if err := json.Unmarshal([]byte(r), &data); err != nil {
		return nil, err
	}
	return data, nil
}

// func (p *Player) GetPlayerRankFromName(playerName string) (int, error) {
// 	clientID, err := p.GetClientIDFromName(playerName)
// 	if err != nil {
// 		return -1, err
// 	}

// 	info, err := p.ClientInfo(clientID)
// 	if err != nil {
// 		return -1, err
// 	}

// 	level, ok := info["level"].(float64)
// 	if !ok {
// 		return -1, fmt.Errorf("could not parse level for player %s", playerName)
// 	}

// 	return int(level), nil
// }

func (p *Player) GetXUIDFromName(playerName string) (string, error) {
	server := server.NewServer(p.Wrapper)
	data, _ := server.FindPlayer(playerName, "", "", "", "", "")
	var result models.PlayerResponse
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return "", err
	}

	if len(result.Clients) == 0 {
		return "", fmt.Errorf("player %s not found", playerName)
	}

	return result.Clients[0].XUID, nil
}

func (p *Player) GetNameFromXUID(xuid string) (string, error) {
	server := server.NewServer(p.Wrapper)
	data, err := server.FindPlayer("", "", xuid, "", "", "")
	if err != nil {
		return "", fmt.Errorf("error finding player: %w", err)
	}

	if strings.TrimSpace(data) == "" {
		return "", fmt.Errorf("empty response from FindPlayer")
	}

	var result models.PlayerResponse
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w\nRaw: %s", err, data)
	}

	if len(result.Clients) == 0 {
		return "", fmt.Errorf("no player with xuid %s", xuid)
	}

	return result.Clients[0].Name, nil
}

func (p *Player) GetNameFromClientID(clientID string) (string, error) {
	info, err := p.ClientInfo(clientID)
	if err != nil {
		return "", err
	}

	name, ok := info["name"].(string)
	if !ok {
		return "", fmt.Errorf("could not find name in client info")
	}

	return name, nil
}

// func (p *Player) GetClientIDFromName(playerName string) (string, error) {
// 	data := p.Wrapper.Server(p.Wrapper).FindPlayer(playerName)
// 	var result models.PlayerResponse
// 	if err := json.Unmarshal([]byte(data), &result); err != nil {
// 		return "", err
// 	}

// 	if len(result.Clients) == 0 {
// 		return "", fmt.Errorf("player %s not found", playerName)
// 	}

// 	return result.Clients[0].ClientID, nil
// }
