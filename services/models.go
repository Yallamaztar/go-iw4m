package services

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

type Help map[string]HelpCategory

type Report struct {
	Origin    string
	Reason    string
	Target    string
	Timestamp string
}

type ServerID struct {
	Server string
	ID     string
}
