package main

type Types struct {
	Links    Links  `json:"_links"`
	Type     string `json:"_type"`
	Total    int64  `json:"total"`
	Count    int64  `json:"count"`
	Embedded struct {
		Elements []TypeElement `json:"elements"`
	} `json:"_embedded"`
}

type TypeElement struct {
	Type        string `json:"_type"`
	Links       Links  `json:"_links"`
	Color       string `json:"color"`
	IsDefault   bool   `json:"isDefault"`
	IsMilestone bool   `json:"isMilestone"`
	Position    int64  `json:"position"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
