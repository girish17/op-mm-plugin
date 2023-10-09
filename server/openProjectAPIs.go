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
	"bytes"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const apiVersionStr string = "/api/v3/"

var client = &http.Client{}

func GetUserDetails(opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "users/me"
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot fetch user details from OpenProject for %s and %s", opURLStr, apiKeyStr)
	}
	return resp, err
}

func GetProjects(opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "projects"
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot fetch Projects from OpenProject for %s and %s", opURLStr, apiKeyStr)
	}
	return resp, err
}

func GetWPsForProject(projectID string, opURLStr string, apiKeyStr string) (*http.Response, error) {
	queryParameters := url.Values{}
	queryParameters.Add("sortBy", "[[\"created_at\",\"desc\"]]")
	fullURL := opURLStr + apiVersionStr + "projects/" + projectID + "/work_packages?" + queryParameters.Encode()
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot fetch WPs from OpenProject for Project ID# %s", projectID)
	}
	return resp, err
}

func GetWorkPackages(opURLStr string, apiKeyStr string) (*http.Response, error) {
	queryParameters := url.Values{}
	queryParameters.Add("sortBy", "[[\"created_at\",\"desc\"]]")
	fullURL := opURLStr + apiVersionStr + "/work_packages?" + queryParameters.Encode()
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot fetch WPs from OpenProject for %s and %s", opURLStr, apiKeyStr)
	}
	return resp, err
}

func PostWP(wpJSON []byte, opURLStr string, apiKeyStr string, notify string) (*http.Response, error) {
	queryParameters := url.Values{}
	queryParameters.Add("notify", notify)
	fullURL := opURLStr + apiVersionStr + "work_packages?" + queryParameters.Encode()
	req, _ := http.NewRequest("POST", fullURL, bytes.NewBuffer(wpJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot save WP in OpenProject")
	}
	return resp, err
}

func DelWP(opURLStr string, apiKeyStr string, wpID string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "work_packages/" + wpID
	req, _ := http.NewRequest("DELETE", fullURL, nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot delete WP# %s in OpenProject ", wpID)
	}
	return resp, err
}

func PostTimeEntriesForm(timeEntriesBodyJSON []byte, opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "time_entries/form"
	req, _ := http.NewRequest("POST", fullURL, bytes.NewBuffer(timeEntriesBodyJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot save Time Entry form in OpenProject ")
	}
	return resp, err
}

func GetTimeEntries(opURLStr string, apiKeyStr string) (*http.Response, error) {
	queryParameters := url.Values{}
	queryParameters.Add("sortBy", "[[\"created_at\",\"desc\"]]")
	fullURL := opURLStr + apiVersionStr + "time_entries?" + queryParameters.Encode()
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot fetch time entries from OpenProject for %s and %s", opURLStr, apiKeyStr)
	}
	return resp, err
}

func PostTimeEntry(timeEntryJSON []byte, opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "time_entries"
	req, _ := http.NewRequest("POST", fullURL, bytes.NewBuffer(timeEntryJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot save Time entry in OpenProject for %s and %s", opURLStr, apiKeyStr)
	}
	return resp, err
}

func GetTimeEntriesSchema(opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "time_entries/schema"
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot fetch Time entries schema from OpenProject for %s and %s", opURLStr, apiKeyStr)
	}
	return resp, err
}

func DelTimeLog(opURLStr string, apiKeyStr string, timeLogID string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "time_entries/" + timeLogID
	req, _ := http.NewRequest("DELETE", fullURL, nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot delete Time entry# %s in OpenProject", timeLogID)
	}
	return resp, err
}

func GetTypes(opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "types"
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot fetch Types from OpenProject for %s and %s", opURLStr, apiKeyStr)
	}
	return resp, err
}

func GetAvailableAssignees(opURLStr string, apiKeyStr string, projectID string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "projects/" + projectID + "/available_assignees"
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil && resp != nil {
		defer resp.Body.Close()
		err = errors.Wrapf(err, "Cannot fetch Assignees from OpenProject for Project ID# %s", projectID)
	}
	return resp, err
}
