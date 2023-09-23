package main

import (
	"bytes"
	"net/http"
	"net/url"
)

const apiVersionStr string = "/api/v3/"

var client = &http.Client{}

func GetUserDetails(opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "users/me"
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func GetProjects(opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "projects"
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
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
	if err != nil {
		defer resp.Body.Close()
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
	if err != nil {
		defer resp.Body.Close()
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
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func DelWP(opURLStr string, apiKeyStr string, wpID string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "work_packages/" + wpID
	req, _ := http.NewRequest("DELETE", fullURL, nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func PostTimeEntriesForm(timeEntriesBodyJSON []byte, opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "time_entries/form"
	req, _ := http.NewRequest("POST", fullURL, bytes.NewBuffer(timeEntriesBodyJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
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
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func PostTimeEntry(timeEntryJSON []byte, opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "time_entries"
	req, _ := http.NewRequest("POST", fullURL, bytes.NewBuffer(timeEntryJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func GetTimeEntriesSchema(opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "time_entries/schema"
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func DelTimeLog(opURLStr string, apiKeyStr string, timeLogID string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "time_entries/" + timeLogID
	req, _ := http.NewRequest("DELETE", fullURL, nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func GetTypes(opURLStr string, apiKeyStr string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "types"
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func GetAvailableAssignees(opURLStr string, apiKeyStr string, projectID string) (*http.Response, error) {
	fullURL := opURLStr + apiVersionStr + "projects/" + projectID + "/available_assignees"
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}
