package util

import (
	"bytes"
	"net/http"
)

const apiVersionStr string = "/api/v3/"

func GetUserDetails(opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", opURLStr+apiVersionStr+"users/me", nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	return client.Do(req)
}

func GetProjects(opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", opURLStr+apiVersionStr+"projects", nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	return client.Do(req)
}

func GetWPsForProject(projectID string, opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", opURLStr+apiVersionStr+"projects/"+projectID+"/work_packages", nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	return client.Do(req)
}

func PostTimeEntriesForm(timeEntriesBodyJSON []byte, opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("POST", opURLStr+apiVersionStr+"time_entries/form", bytes.NewBuffer(timeEntriesBodyJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	return client.Do(req)
}

func GetTimeEntries(opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", opURLStr+apiVersionStr+"time_entries", nil)
	req.SetBasicAuth("apikey", apiKeyStr)
	return client.Do(req)
}

func PostTimeEntry(timeEntryJSON []byte, opURLStr string, apiKeyStr string) (*http.Response, error) {
	req, _ := http.NewRequest("POST", opURLStr+apiVersionStr+"time_entries", bytes.NewBuffer(timeEntryJSON))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth("apikey", apiKeyStr)
	return client.Do(req)
}
