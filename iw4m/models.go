package iw4m

type CommandHelp struct {
	Alias          string `json:"alias"`
	Description    string `json:"description"`
	RequiresTarget string `json:"requires_target"`
	Syntax         string `json:"syntax"`
	MinLevel       string `json:"min_level"`
}

type HelpCategory struct {
	Commands map[string]CommandHelp `json:"commands"`
}

type HelpModel map[string]HelpCategory

type ReportModel struct {
	Origin    string
	Reason    string
	Target    string
	Timestamp string
}

type ServerID struct {
	Server string
	ID     string
}

type ChatModel struct {
	Origin  string
	Message string
}

type PlayerModel struct {
	Role string
	Name string
	XUID string
	URL  string
}

type RecentClientModel struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
	Country   string `json:"country,omitempty"`
	IPAddress string `json:"ip_address"`
	LastSeen  string `json:"last_seen"`
}

type AuditLogModel struct {
	Type   string
	Origin string
	Href   string
	Target string
	Data   string
	Time   string
}

type AdminModel struct {
	Name          string
	Role          string
	Game          string
	LastConnected string
}

type TopPlayerModel struct {
	Rank   string            `json:"rank"`
	Name   string            `json:"name"`
	Link   string            `json:"link"`
	Rating string            `json:"rating"`
	Stats  map[string]string `json:"stats"`
}

type AdvancedStatsModel struct {
	Name         string                   `json:"name"`
	Link         string                   `json:"link"`
	IconURL      string                   `json:"icon_url"`
	Summary      string                   `json:"summary"`
	PlayerStats  []StatEntry              `json:"player_stats"`
	HitLocations map[string][]HitLocation `json:"hit_locations"`
	WeaponUsages map[string][]WeaponUsage `json:"weapon_usages"`
}

type StatEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type HitLocation struct {
	Location   string `json:"location"`
	Hits       string `json:"hits"`
	Percentage string `json:"percentage"`
	Damage     string `json:"damage"`
}

type WeaponUsage struct {
	Weapon              string `json:"weapon"`
	FavoriteAttachments string `json:"favorite_attachments"`
	Kills               string `json:"kills"`
	Hits                string `json:"hits"`
	Damage              string `json:"damage"`
	Usage               string `json:"usage"`
}

type PlayerResponse struct {
	Clients []struct {
		XUID string `json:"xuid"`
	} `json:"clients"`
}
