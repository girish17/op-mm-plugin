package main

type AvailableAssignees struct {
	Links    Links  `json:"_links"`
	Type     string `json:"_type"`
	Total    int64  `json:"total"`
	Count    int64  `json:"count"`
	Embedded struct {
		Elements []AvailableAssigneesElement `json:"elements"`
	} `json:"_embedded"`
}

type AvailableAssigneesElement struct {
	Links     AvailableAssigneeLinks `json:"_links"`
	Type      string                 `json:"_type"`
	ID        int                    `json:"id"`
	Avatar    string                 `json:"avatar"`
	Email     string                 `json:"email"`
	FirstName string                 `json:"firstName"`
	LastName  string                 `json:"lastName"`
	Login     string                 `json:"login"`
	Status    string                 `json:"status"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
}

type AvailableAssigneeLinks struct {
	Delete Lock `json:"delete"`
	Self   Self `json:"self"`
	Lock   Lock `json:"lock"`
}

type Lock struct {
	Href   string `json:"href"`
	Method string `json:"method"`
	Title  string `json:"title"`
}
