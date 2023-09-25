/*
 *    Copyright 2023 Girish M
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

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
	CustomField  string       `json:"-"`
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
	CustomField string       `json:"-"`
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
	CustomField string  `json:"-"`
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
	CustomField       string            `json:"-"`
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
	CustomField float64   `json:"-"`
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

func (tel *TimeEntryList) UnmarshalJSON(data []byte) error {
	var jsonData map[string]interface{}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	tel.Type, _ = jsonData["_type"].(string)
	tel.Total, _ = jsonData["total"].(int)
	tel.Count, _ = jsonData["count"].(int)
	tel.PageSize, _ = jsonData["pageSize"].(int)
	tel.Offset, _ = jsonData["offset"].(int)
	embeddedData, ok := jsonData["_embedded"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("missing _embedded field or it's not an object")
	}
	elementsData, ok := embeddedData["elements"].([]interface{})
	if !ok {
		return fmt.Errorf("missing elements field or it's not an array")
	}

	for _, element := range elementsData {
		elementBytes, err := json.Marshal(element)
		if err != nil {
			return err
		}
		var timeElement TimeElement
		if err := json.Unmarshal(elementBytes, &timeElement); err != nil {
			return err
		}
		tel.Embedded.Elements = append(tel.Embedded.Elements, timeElement)
	}
	return nil
}

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

func (te *TimeElement) MarshalJSON() ([]byte, error) {
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

func (te *TimeElement) UnmarshalJSON(data []byte) error {
	var jsonData map[string]interface{}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	te.Type, _ = jsonData["_type"].(string)
	te.ID = int(jsonData["id"].(float64))
	commentBytes, _ := json.Marshal(jsonData["comment"])
	if err := json.Unmarshal(commentBytes, &te.Comment); err != nil {
		return err
	}
	te.SpentOn, _ = jsonData["spentOn"].(string)
	te.Hours, _ = jsonData["hours"].(string)
	te.CreatedAt, _ = jsonData["createdAt"].(string)
	te.UpdatedAt, _ = jsonData["updatedAt"].(string)
	linksBytes, ok := jsonData["_links"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("missing _links field or it's not an object")
	}

	linksData, _ := json.Marshal(linksBytes)
	if err := json.Unmarshal(linksData, &te.Links); err != nil {
		return err
	}

	fieldName := getCustomFieldName()
	te.CustomField, _ = jsonData[fieldName].(float64)
	return nil
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
	if customFieldForBillableHours == "" {
		resp, err := GetTimeEntriesSchema(OpURLStr, APIKeyStr)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			_ = json.Unmarshal(body, &timeEntriesSchema)
			customFieldForBillableHours = findFirstKVContainingBillable(timeEntriesSchema)
			return customFieldForBillableHours
		}
	}
	return customFieldForBillableHours
}

func findFirstKVContainingBillable(jsonObj map[string]interface{}) string {
	customFieldRegex := regexp.MustCompile(`customField\d+`)
	for key, v := range jsonObj {
		if customFieldRegex.MatchString(key) {
			var strValue, ok = v.(map[string]interface{})
			if ok && strings.Contains(strings.ToLower(strValue["name"].(string)), "billable") {
				return key
			}
		}
	}
	return "customField1"
}
