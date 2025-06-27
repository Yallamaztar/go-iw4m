package iw4m

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
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
	path := fmt.Sprintf("%s/Console/Execute?serverId=%s&command=%s",
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

func (s *Server) Help() (HelpModel, error) {
	help := make(HelpModel)

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
	doc.Find("div.col-12.align-self-center.text-center.text-lg-left.col-lg-4").Each(func(i int, s *goquery.Selection) {
		spans := s.Find("span")
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
		func(i int, s *goquery.Selection) {
			spans := s.Find("span")
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
		func(i int, s *goquery.Selection) {
			if span := s.Find("span.text-primary"); span.Length() > 0 {
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

func (s *Server) Reports() ([]ReportModel, error) {
	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/Action/RecentReportsForm/", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	var reports []ReportModel
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

			reports = append(reports, ReportModel{
				Origin:    origin,
				Reason:    reason,
				Target:    target,
				Timestamp: timestamp,
			})
			i++
		})

	return reports, nil
}

func (s *Server) ServerIDs() ([]ServerID, error) {
	html := s.Wrapper.DoRequest(fmt.Sprintf("%s/Console", s.Wrapper.BaseURL))

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var serverIDs []ServerID
	doc.Find("select#console_server_select option").Each(
		func(i int, s *goquery.Selection) {
			name := strings.TrimSpace(s.Text())
			id, exists := s.Attr("value")
			if exists {
				serverIDs = append(serverIDs, ServerID{
					Server: name, ID: id,
				})
			}
		})

	return serverIDs, nil
}

func (s *Server) SendCommand(command string) string {
	encodedCommand := url.QueryEscape(command)
	path := fmt.Sprintf("%s/Console/Execute?serverId=%s&command=%s",
		s.Wrapper.BaseURL, s.Wrapper.ServerID, encodedCommand)

	r := s.Wrapper.DoRequest(path)
	return r
}

func (s *Server) ReadChat() ([]ChatModel, error) {
	html := s.Wrapper.DoRequest(s.Wrapper.BaseURL)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var chat []ChatModel
	doc.Find("div.text-truncate").Each(
		func(i int, s *goquery.Selection) {
			var origin string
			var message string

			if span := s.Find("span").First(); span.Length() > 0 {
				if tag := span.Find("colorcode"); tag.Length() > 0 {
					origin = tag.Text()
				}
			}

			if spans := s.Find("span"); spans.Length() > 1 {
				if messageTag := spans.Eq(1).Find("colorcode"); messageTag.Length() > 0 {
					message = messageTag.Text()
				}
			}

			if origin != "" && message != "" {
				chat = append(chat, ChatModel{Origin: origin, Message: message})
			}
		})

	return chat, nil
}

func (s *Server) FindPlayer(name, xuid string, count, offset, direction int) (string, error) {
	if name == "" {
		return "", nil
	}

	params := url.Values{}
	params.Set("name", name)
	params.Set("xuid", xuid)
	params.Set("count", strconv.Itoa(count))
	params.Set("offset", strconv.Itoa(offset))
	params.Set("direction", strconv.Itoa(direction))
	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/api/client/find?%s",
		s.Wrapper.BaseURL, params.Encode()))

	return r, nil
}

func (s *Server) GetPlayers() ([]PlayerModel, error) {
	var players []PlayerModel

	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	selectors := map[string]string{
		"creator":   "level-color-7.no-decoration.text-truncate.ml-5.mr-5",
		"owner":     "level-color-6.no-decoration.text-truncate.ml-5.mr-5",
		"moderator": "level-color-5.no-decoration.text-truncate.ml-5.mr-5",
		"senior":    "level-color-4.no-decoration.text-truncate.ml-5.mr-5",
		"admin":     "level-color-3.no-decoration.text-truncate.ml-5.mr-5",
		"trusted":   "level-color-2.no-decoration.text-truncate.ml-5.mr-5",
		"user":      "text-light-dm.text-dark-lm.no-decoration.text-truncate.ml-5.mr-5",
		"flagged":   "level-color-1.no-decoration.text-truncate.ml-5.mr-5",
		"banned":    "level-color--1.no-decoration.text-truncate.ml-5.mr-5",
	}
	for role, selector := range selectors {
		doc.Find("a." + selector).Each(
			func(i int, s *goquery.Selection) {
				colorcode := s.Find("colorcode")
				if colorcode.Length() > 0 {
					name := strings.TrimSpace(colorcode.Text())
					href, exists := s.Attr("href")
					if exists && len(href) >= 17 {
						xuid := href[16:]
						players = append(players, PlayerModel{
							Role: role,
							Name: name,
							XUID: xuid,
							URL:  strings.TrimSpace(href),
						})
					}
				}
			})
	}

	return players, nil
}

func (s *Server) AdminRoles() ([]string, error) {
	var roles []string

	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/Action/editForm/?id=2&meta=\"\"", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	doc.Find(`select[name="level"] option`).Each(func(i int, s *goquery.Selection) {
		if val, exists := s.Attr("value"); exists && val != "" {
			roles = append(roles, val)
		}
	})

	return roles, nil
}

func (s *Server) GetRoles() ([]string, error) {
	var roles []string

	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/Action/editForm/?id=2&meta=\"\"", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	doc.Find("select option").Each(func(i int, s *goquery.Selection) {
		colorcode := s.Find("colorcode")
		if colorcode.Length() == 0 {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				roles = append(roles, text)
			}
			return
		}

		text := strings.TrimSpace(colorcode.Text())
		if text != "" {
			roles = append(roles, text)
		}
	})

	return roles, nil
}

func (s *Server) RecentClients(offset int) ([]RecentClientModel, error) {
	var recentClients []RecentClientModel

	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/Action/RecentClientsForm?offset=%d&count=20", s.Wrapper.BaseURL, offset))

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	doc.Find("div.bg-very-dark-dm.bg-light-ex-lm.p-15.rounded.mb-10").Each(func(i int, entry *goquery.Selection) {
		var client RecentClientModel

		user := entry.Find("div.d-flex.flex-row").First()
		if user.Length() > 0 {
			a := user.Find("a.h4.mr-auto")
			if a.Length() > 0 {
				client.Name = strings.TrimSpace(a.Find("colorcode").Text())
				if link, exists := a.Attr("href"); exists {
					client.Link = link
				}
			}

			tooltip := user.Find("div[data-toggle='tooltip']")
			if tooltip.Length() > 0 {
				if country, exists := tooltip.Attr("data-title"); exists {
					client.Country = country
				}
			}
		}

		client.IPAddress = strings.TrimSpace(entry.Find("div.align-self-center.mr-auto").Text())
		client.LastSeen = strings.TrimSpace(entry.Find("div.align-self-center.text-muted.font-size-12").Text())

		if client.Name != "" {
			recentClients = append(recentClients, client)
		}
	})

	return recentClients, nil
}

func (s *Server) RecentAuditLog() (*AuditLogModel, error) {
	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/Admin/AuditLog", s.Wrapper.BaseURL))

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	tbody := doc.Find("#audit_log_table_body")
	if tbody.Length() == 0 {
		return nil, nil
	}

	tr := tbody.Find("tr.d-none.d-lg-table-row.bg-dark-dm.bg-light-lm").First()
	if tr.Length() == 0 {
		return nil, nil
	}

	columns := tr.Find("td")
	originAnchor := columns.Eq(1).Find("a").First()
	targetAnchor := columns.Eq(2).Find("a").First()

	originName := originAnchor.Text()
	href, _ := originAnchor.Attr("href")

	target := ""
	if targetAnchor.Length() > 0 {
		target = targetAnchor.Text()
	} else {
		target = columns.Eq(2).Text()
	}

	return &AuditLogModel{
		Type:   strings.TrimSpace(columns.Eq(0).Text()),
		Origin: strings.TrimSpace(originName),
		Href:   strings.TrimSpace(href),
		Target: strings.TrimSpace(target),
		Data:   strings.TrimSpace(columns.Eq(4).Text()),
		Time:   strings.TrimSpace(columns.Eq(5).Text()),
	}, nil
}

func (s *Server) AuditLogs(count int) ([]AuditLogModel, error) {
	var auditLogs []AuditLogModel

	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/Admin/AuditLog", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	tbody := doc.Find("#audit_log_table_body")
	if tbody.Length() == 0 {
		return auditLogs, nil
	}

	rows := tbody.Find("tr.d-none.d-lg-table-row.bg-dark-dm.bg-light-lm")
	rows.EachWithBreak(
		func(i int, tr *goquery.Selection) bool {
			if i >= count {
				return false
			}

			columns := tr.Find("td")
			if columns.Length() < 6 {
				return true
			}

			originAnchor := columns.Eq(1).Find("a").First()
			targetAnchor := columns.Eq(2).Find("a").First()

			originName := originAnchor.Text()
			href, _ := originAnchor.Attr("href")

			target := ""
			if targetAnchor.Length() > 0 {
				target = targetAnchor.Text()
			} else {
				target = columns.Eq(2).Text()
			}

			auditLogs = append(auditLogs, AuditLogModel{
				Type:   strings.TrimSpace(columns.Eq(0).Text()),
				Origin: strings.TrimSpace(originName),
				Href:   strings.TrimSpace(href),
				Target: strings.TrimSpace(target),
				Data:   strings.TrimSpace(columns.Eq(4).Text()),
				Time:   strings.TrimSpace(columns.Eq(5).Text()),
			})

			return true
		})

	return auditLogs, nil
}

func (s *Server) Admins(role string, count int) ([]AdminModel, error) {
	if role == "" {
		role = "all"
	}

	var admins []AdminModel

	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/Client/Privileged", s.Wrapper.BaseURL))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	doc.Find("table.table.mb-20").EachWithBreak(
		func(i int, table *goquery.Selection) bool {
			header := table.Find("thead tr th").First()
			if header.Length() == 0 {
				return true
			}

			_role := strings.TrimSpace(header.Text())
			if strings.ToLower(role) == "all" || strings.EqualFold(_role, role) {
				tbody := table.Find("tbody")
				if tbody.Length() == 0 {
					return true
				}

				tbody.Find("tr").EachWithBreak(
					func(j int, row *goquery.Selection) bool {
						tag := row.Find("a.text-force-break").First()
						if tag.Length() == 0 {
							return true // next row
						}
						name := strings.TrimSpace(tag.Text())

						badge := row.Find("div.badge").First()
						game := "N/A"
						if badge.Length() > 0 {
							game = strings.TrimSpace(badge.Text())
						}

						tds := row.Find("td")
						lastConnected := "N/A"
						if tds.Length() > 0 {
							lastConnected = strings.TrimSpace(tds.Eq(tds.Length() - 1).Text())
						}

						admins = append(admins, AdminModel{
							Name:          name,
							Role:          role,
							Game:          game,
							LastConnected: lastConnected,
						})

						if count > 0 && len(admins) >= count {
							return false
						}
						return true
					})
				if count > 0 && len(admins) >= count {
					return false
				}
			}
			return true
		})

	return admins, nil
}

func (s *Server) TopPlayers(count int) ([]TopPlayerModel, error) {
	var topPlayers []TopPlayerModel

	r := s.Wrapper.DoRequest(fmt.Sprintf("%s/Stats/GetTopPlayersAsync?offset=0&count=%d&serverId=0", s.Wrapper.BaseURL, count))
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(r))
	if err != nil {
		return nil, err
	}

	doc.Find("div.card.m-0.mt-15.p-20.d-flex.flex-column.flex-md-row.justify-content-between").Each(func(i int, entry *goquery.Selection) {
		rankDiv := entry.Find("div.d-flex.flex-column.w-full.w-md-quarter")
		if rankDiv.Length() == 0 {
			return
		}

		rank := strings.TrimSpace(rankDiv.Find("div.d-flex.text-muted > div").Text())
		name := strings.TrimSpace(rankDiv.Find("div.d-flex.flex-row colorcode").Text())
		link, _ := rankDiv.Find("div.d-flex.flex-row a").Attr("href")
		rating := strings.TrimSpace(rankDiv.Find("div.font-size-14 span").Text())

		stats := make(map[string]string)
		rankDiv.Find("div.d-flex.flex-column.font-size-12.text-right.text-md-left div").Each(func(i int, div *goquery.Selection) {
			primary := strings.TrimSpace(div.Find("span.text-primary").Text())
			secondary := strings.TrimSpace(div.Find("span.text-muted").Text())
			if primary != "" && secondary != "" {
				stats[secondary] = primary
			}
		})

		topPlayers = append(topPlayers, TopPlayerModel{
			Rank:   "#" + rank,
			Name:   name,
			Link:   link,
			Rating: rating,
			Stats:  stats,
		})
	})

	return topPlayers, nil
}
