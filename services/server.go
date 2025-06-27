package services

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Yallamaztar/go-iw4m/wrapper"
)

type Server struct {
	Wrapper *wrapper.IW4MWrapper
}

// Constructor to create Server from IW4MWrapper instance
func NewServer(w *wrapper.IW4MWrapper) *Server {
	return &Server{Wrapper: w}
}

func (s *Server) ServerUptime() string {
	path := fmt.Sprintf("%s/Console/Execute?serverId=%d&command=%s",
		s.Wrapper.BaseURL, s.Wrapper.ServerID, url.QueryEscape("!uptime"))
	return s.Wrapper.DoRequest(path)
}

func (s *Server) LoginToken() string {
	return s.Wrapper.DoRequest(fmt.Sprintf("%s/Action/GenerateLoginTokenAsync/", s.Wrapper.BaseURL))
}

func (s *Server) Status() string {
	return s.Wrapper.DoRequest(fmt.Sprintf("%s/api/status", s.Wrapper.BaseURL))
}

func (s *Server) Info() string {
	return s.Wrapper.DoRequest(fmt.Sprintf("%s/api/info", s.Wrapper.BaseURL))
}

func (s *Server) Help() (Help, error) {
	help := make(Help)

	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/Home/Help", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	doc.Find("div.command-assembly-container").Each(func(i int, container *goquery.Selection) {
		title := strings.TrimSpace(container.Find("h2.content-title.mb-lg-20.mt-20").Text())
		if title == "" {
			return
		}
		if _, exists := help[title]; !exists {
			help[title] = HelpCategory{Commands: make(map[string]CommandHelp)}
		}

		container.Find("tr.d-none.d-lg-table-row.bg-dark-dm.bg-light-lm").Each(func(_ int, tr *goquery.Selection) {
			tds := tr.Find("td")
			if tds.Length() < 6 {
				return
			}
			name := strings.TrimSpace(tds.Eq(0).Text())
			alias := strings.TrimSpace(tds.Eq(1).Text())
			description := strings.TrimSpace(tds.Eq(2).Text())
			requiresTarget := strings.TrimSpace(tds.Eq(3).Text())
			syntax := strings.TrimSpace(tds.Eq(4).Text())
			minLevel := strings.TrimSpace(tr.Find("td.text-right").Text())

			help[title].Commands[name] = CommandHelp{
				Alias:          alias,
				Description:    description,
				RequiresTarget: requiresTarget,
				Syntax:         syntax,
				MinLevel:       minLevel,
			}
		})
	})

	return help, nil
}

func (s *Server) MapName() (string, error) {
	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return "", err
	}

	var mapName string
	doc.Find("div.col-12.align-self-center.text-center.text-lg-left.col-lg-4").Each(func(i int, sel *goquery.Selection) {
		spans := sel.Find("span")
		if spans.Length() > 0 {
			mapName = strings.TrimSpace(spans.Eq(0).Text())
		}
	})

	if mapName == "" {
		return "", fmt.Errorf("map name not found")
	}
	return mapName, nil
}

func (s *Server) Gamemode() (string, error) {
	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return "", err
	}

	var gameMode string
	doc.Find("div.col-12.align-self-center.text-center.text-lg-left.col-lg-4").Each(
		func(i int, sel *goquery.Selection) {
			spans := sel.Find("span")
			if spans.Length() > 2 {
				gameMode = strings.TrimSpace(spans.Eq(2).Text())
			}
		})

	return gameMode, nil
}

func (s *Server) Iw4mVersion() (string, error) {
	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return "", err
	}

	var version string
	doc.Find("a.sidebar-link").Each(
		func(i int, sel *goquery.Selection) {
			if span := sel.Find("span.text-primary"); span.Length() > 0 {
				version = strings.TrimSpace(span.Text())
				return
			}
		})
	return version, nil
}

func (s *Server) LoggedInAs() (string, error) {
	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return "", err
	}

	var name string
	div := doc.Find("div.sidebar-link.font-size-12.font-weight-light").First()
	if div.Length() > 0 {
		colorcode := div.Find("colorcode")
		if colorcode.Length() > 0 {
			name = strings.TrimSpace(colorcode.Text())
		}
	}
	return name, nil
}

func cleanText(text string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(text, " "))
}

func (s *Server) Rules() []string {
	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/About", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil
	}

	var rules []string
	doc.Find("div.card.m-0.rounded").Each(func(i int, card *goquery.Selection) {
		if card.Find("h5.text-primary.mt-0.mb-0").Length() > 0 {
			card.Find("div.rule").Each(func(j int, ruleDiv *goquery.Selection) {
				rawText := ruleDiv.Text()
				cleaned := cleanText(rawText)
				rules = append(rules, cleaned)
			})
		}
	})
	return rules
}

func (s *Server) GetReports() ([]Report, error) {
	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/Action/RecentReportsForm/", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	var reports []Report

	// Get timestamps from the report blocks
	timestamps := []string{}
	doc.Find("div.rounded.bg-very-dark-dm.bg-light-ex-lm.mt-10.mb-10.p-10").Each(func(i int, block *goquery.Selection) {
		timestamp := strings.TrimSpace(block.Find("div.font-weight-bold").Text())
		timestamps = append(timestamps, timestamp)
	})

	i := 0
	doc.Find("div.font-size-12").Each(
		func(_ int, entry *goquery.Selection) {
			origin := strings.TrimSpace(entry.Find("a").Text())

			reasonTag := entry.Find("span.text-white-dm.text-black-lm colorcode")
			reason := ""
			if reasonTag.Length() > 0 {
				reason = strings.TrimSpace(reasonTag.Text())
			}

			targetTag := entry.Find("span.text-highlight a")
			target := ""
			if targetTag.Length() > 0 {
				target = strings.TrimSpace(targetTag.Text())
			}

			timestamp := ""
			if i < len(timestamps) {
				timestamp = timestamps[i]
			}

			reports = append(reports, Report{
				Origin:    origin,
				Reason:    reason,
				Target:    target,
				Timestamp: timestamp,
			})
			i++
		})

	return reports, nil
}
