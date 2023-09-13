package main

type ActivityLinks struct {
	Self     Self   `json:"self"`
	Projects []Self `json:"projects"`
}

type AllowedValues struct {
	Type     string        `json:"type"`
	ID       int           `json:"id"`
	Name     string        `json:"name"`
	Position int           `json:"position"`
	Default  bool          `json:"default"`
	Links    ActivityLinks `json:"_links"`
}

type Activity struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	Required   bool   `json:"required"`
	HasDefault bool   `json:"hasDefault"`
	Writable   bool   `json:"writable"`
	Location   string `json:"location"`
	Embedded   struct {
		AllowedValues []AllowedValues `json:"allowedValues"`
	} `json:"_embedded"`
}

type Dependency struct {
}

type Rtl struct {
}

type TimeEntryOption struct {
	Rtl Rtl `json:"rtl"`
}

type AllowedValue struct {
	Href string `json:"href"`
}

type LinksTimeEntryWP struct {
	AllowedValues AllowedValue `json:"allowedValues"`
}

type TimeEntryWP struct {
	Type       string           `json:"type"`
	Name       string           `json:"name"`
	Required   bool             `json:"required"`
	HasDefault bool             `json:"hasDefault"`
	Writable   bool             `json:"writable"`
	Location   string           `json:"location"`
	Links      LinksTimeEntryWP `json:"_links"`
}

type ID struct {
	Type       string          `json:"type"`
	Name       string          `json:"name"`
	Required   bool            `json:"required"`
	HasDefault bool            `json:"hasDefault"`
	Writable   bool            `json:"writable"`
	Options    TimeEntryOption `json:"options"`
}

type SchemaLinks struct {
}

type Schema struct {
	Type         string       `json:"_type"`
	Dependencies []Dependency `json:"_dependencies"`
	ID           ID           `json:"id"`
	CreatedAt    ID           `json:"createdAt"`
	UpdatedAt    ID           `json:"updatedAt"`
	SpentOn      ID           `json:"spentOn"`
	Hours        ID           `json:"hours"`
	User         ID           `json:"user"`
	WorkPackage  TimeEntryWP  `json:"workPackage"`
	Project      TimeEntryWP  `json:"project"`
	Activity     Activity     `json:"activity"`
	CustomField  string       `json:"customField1"`
	Links        SchemaLinks  `json:"_links"`
}

type Link struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type PayloadLinks struct {
	Project     Link `json:"project"`
	Activity    Link `json:"activity"`
	WorkPackage Link `json:"workPackage"`
}

type Comment struct {
	Format string `json:"format"`
	Raw    string `json:"raw"`
	HTML   string `json:"html"`
}

type Payload struct {
	Links       PayloadLinks `json:"_links"`
	Hours       string       `json:"hours"`
	Comment     Comment      `json:"comment"`
	SpentOn     string       `json:"spentOn"`
	CustomField string       `json:"customField1"`
}

type ValidationError struct {
}

type EmbeddedTimeEntry struct {
	Payload          Payload         `json:"payload"`
	Schema           Schema          `json:"schema"`
	ValidationErrors ValidationError `json:"validationErrors"`
}

type LinksTimeEntry struct {
	Self     Self `json:"self"`
	Validate Self `json:"validate"`
	Commit   Self `json:"commit"`
}

type TimeEntries struct {
	Type     string            `json:"_type"`
	Embedded EmbeddedTimeEntry `json:"_embedded"`
	Links    LinksTimeEntry    `json:"_links"`
}

type TimeEntriesBody struct {
	Links struct {
		WorkPackage struct {
			Href string `json:"href"`
		} `json:"workPackage"`
	} `json:"_links"`
}

type TimeEntryPostBody struct {
	Links struct {
		WorkPackage struct {
			Href string `json:"href"`
		} `json:"workPackage"`
		Activity struct {
			Href string `json:"href"`
		} `json:"activity"`
		Project struct {
			Href string `json:"href"`
		} `json:"project"`
	} `json:"_links"`
	Hours       string  `json:"hours"`
	Comment     Comment `json:"comment"`
	SpentOn     string  `json:"spentOn"`
	CustomField string  `json:"customField1"`
}

type UpdateImmediately struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}

type Delete struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}

type TimeLinks struct {
	Self              Self              `json:"self"`
	UpdateImmediately UpdateImmediately `json:"updateImmediately"`
	Delete            Delete            `json:"delete"`
	Project           Link              `json:"project"`
	WorkPackage       Link              `json:"workPackage"`
	User              Link              `json:"user"`
	Activity          Link              `json:"activity"`
	CustomField       string            `json:"customField1"`
}

type TimeElement struct {
	Type        string    `json:"_type"`
	ID          int       `json:"id"`
	Comment     Comment   `json:"comment"`
	SpentOn     string    `json:"spentOn"`
	Hours       string    `json:"hours"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
	Links       TimeLinks `json:"_links"`
	CustomField string    `json:"customField1"`
}

type TimeEntryList struct {
	Type     string `json:"_type"`
	Total    int    `json:"total"`
	Count    int    `json:"count"`
	PageSize int    `json:"pageSize"`
	Offset   int    `json:"offset"`
	Embedded struct {
		Elements []TimeElement `json:"elements"`
	} `json:"_embedded"`
}

/*
func (sl Schema) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["_type"] = sl.Type
	data["_dependencies"] = sl.Dependencies
	data["id"] = sl.ID
	data["createdAt"] = sl.CreatedAt
	data["updatedAt"] = sl.UpdatedAt
	data["spentOn"] = sl.SpentOn
	data["hours"] = sl.Hours
	data["user"] = sl.User
	data["workPackage"] = sl.WorkPackage
	data["project"] = sl.Project
	data["activity"] = sl.Activity
	data["_links"] = sl.Links

	fieldName := getCustomFieldName()
	data[fieldName] = sl.CustomField

	return json.Marshal(data)
}

func (pl Payload) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["_links"] = pl.Links
	data["hours"] = pl.Hours
	data["comment"] = pl.Comment
	data["spentOn"] = pl.SpentOn

	fieldName := getCustomFieldName()
	data[fieldName] = pl.CustomField

	return json.Marshal(data)
}

func (tl TimeLinks) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["self"] = tl.Self
	data["updateImmediately"] = tl.UpdateImmediately
	data["delete"] = tl.Delete
	data["project"] = tl.Project
	data["workPackage"] = tl.WorkPackage
	data["user"] = tl.User
	data["activity"] = tl.Activity

	fieldName := getCustomFieldName()
	data[fieldName] = tl.CustomField

	return json.Marshal(data)
}

func (te TimeElement) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["_type"] = te.Type
	data["id"] = te.ID
	data["comment"] = te.Comment
	data["spentOn"] = te.SpentOn
	data["hours"] = te.Hours
	data["createdAt"] = te.CreatedAt
	data["updatedAt"] = te.UpdatedAt
	data["_links"] = te.Links

	fieldName := getCustomFieldName()
	data[fieldName] = te.CustomField

	return json.Marshal(data)
}

func (tb TimeEntryPostBody) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["_links"] = tb.Links
	data["hours"] = tb.Hours
	data["spentOn"] = tb.SpentOn
	data["comment"] = tb.Comment

	fieldName := getCustomFieldName()
	data[fieldName] = tb.CustomField

	return json.Marshal(data)
}

func getCustomFieldName() string {
	resp, err := GetTimeEntriesSchema(OpURLStr, APIKeyStr)
	if err == nil {
		body, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		var jsonBody map[string]interface{}
		_ = json.Unmarshal(body, &jsonBody)
		return findFirstKeyAndValueContainingBillable(jsonBody)
	}
	return "customField1"
}

func findFirstKeyAndValueContainingBillable(jsonObj map[string]interface{}) string {
	customFieldRegex := regexp.MustCompile(`(?i)customfield\d+`)
	for key, value := range jsonObj {
		if customFieldRegex.MatchString(key) {
			strValue, ok := value.(string)
			if ok && strings.Contains(strings.ToLower(strValue), "billable") {
				return key
			}
		}
		if childObj, ok := value.(map[string]interface{}); ok {
			if matchedKey := findFirstKeyAndValueContainingBillable(childObj); matchedKey != "" {
				return matchedKey
			}
		}
	}
	return "customField1"
}
*/
