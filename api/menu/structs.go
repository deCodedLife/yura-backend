package api

type MenuItem struct {
	Title    string    `json:"title"`
	Icon     string    `json:"icon"`
	URL      string    `json:"url"`
	Children *MenuItem `json:"children"`
}

type Menu struct {
	Items []MenuItem `json:"items"`
}
