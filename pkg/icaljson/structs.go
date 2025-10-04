package icaljson

type Calendar struct {
	ProdID  string  `json:"prodid,omitempty"`
	Version string  `json:"version,omitempty"`
	Events  []Event `json:"events,omitempty"`
}

type Event struct {
	UID         string   `json:"uid,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	Description string   `json:"description,omitempty"`
	Start       string   `json:"start,omitempty"`
	End         string   `json:"end,omitempty"`
	Location    string   `json:"location,omitempty"`
	URL         string   `json:"url,omitempty"`
	Status      string   `json:"status,omitempty"`
	Categories  []string `json:"categories,omitempty"`
}
