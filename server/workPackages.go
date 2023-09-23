package main

type WorkPackages struct {
	Type     string `json:"_type"`
	Links    Links  `json:"_links"`
	Total    int    `json:"total"`
	Count    int    `json:"count"`
	Embedded struct {
		Elements []Element `json:"elements"`
	} `json:"_embedded"`
}

type WorkPackagePostBody struct {
	Links struct {
		Type struct {
			Href string `json:"href"`
		} `json:"type"`
		Project struct {
			Href string `json:"href"`
		} `json:"project"`
	} `json:"_links"`
	Subject  string `json:"subject"`
	Assignee struct {
		Href string `json:"href"`
	}
}
