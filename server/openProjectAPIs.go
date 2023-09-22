package main

import (
	"bytes"
	"net/http"
)

const apiVersionStr string = "/api/v3/"

var client = &http.Client{}

func GetUserDetails(opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", opURLStr+apiVersionStr+"users/me", nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func GetProjects(opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", opURLStr+apiVersionStr+"projects", nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func GetWPsForProject(projectID string, opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", opURLStr+apiVersionStr+"projects/"+projectID+"/work_packages", nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func PostTimeEntriesForm(timeEntriesBodyJSON []byte, opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("POST", opURLStr+apiVersionStr+"time_entries/form", bytes.NewBuffer(timeEntriesBodyJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func GetTimeEntries(opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", opURLStr+apiVersionStr+"time_entries", nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func PostTimeEntry(timeEntryJSON []byte, opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("POST", opURLStr+apiVersionStr+"time_entries", bytes.NewBuffer(timeEntryJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func GetTimeEntriesSchema(opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", opURLStr+apiVersionStr+"time_entries/schema", nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}

func DelTimeLog(opURLStr string, apiKeyStr string, timeLogID string) (*http.Response, error) {
	url := opURLStr + apiVersionStr + "time_entries/" + timeLogID
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
	}
	return resp, err
}
