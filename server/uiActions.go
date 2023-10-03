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
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/senseyeio/duration"
)

var menuPost *model.Post

var respHTTP http.Response

var OpURLStr string

var APIKeyStr string

var projectID string

var timeLogID string

var typeID string

var timeEntry string

var wpEntry string

var wpID string

var activityID string

var customFieldForBillableHours string

var timeEntriesSchema map[string]interface{}

var subscribedChannelID string

func OpAuth(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	p.API.LogDebug("Request body from dialog submit: ", jsonBody)
	dialogCancelled := jsonBody["cancelled"].(bool)
	user, err := p.API.GetBot(opBot, false)
	if err != nil {
		p.API.LogError(err.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		var post *model.Post
		if dialogCancelled {
			p.API.LogInfo("Op Auth Dialog cancelled by user.")
			post = getCreatePostMsg(user.UserId, jsonBody["channel_id"].(string), messages.DlgCancelMsg)
		} else {
			submission := jsonBody["submission"].(map[string]interface{})
			mmUserID := jsonBody["user_id"].(string)
			for key, value := range submission {
				p.API.LogInfo("Storing OpenProject auth credentials: " + key + ":" + value.(string))
				_ = p.API.KVSet(key, []byte(value.(string)))
			}
			setOPStr(p)
			opUserID := []byte(APIKeyStr + " " + OpURLStr)

			_ = p.API.KVDelete(OpURLStr)
			_ = p.API.KVDelete(APIKeyStr)

			channelID := jsonBody["channel_id"].(string)
			resp, err := GetUserDetails(OpURLStr, APIKeyStr)
			if err == nil {
				opResBody, _ := io.ReadAll(resp.Body)
				defer resp.Body.Close()
				var opJSONRes map[string]string
				_ = json.Unmarshal(opResBody, &opJSONRes)
				p.API.LogDebug("Response from op-mattermost: ", opJSONRes)
				if opJSONRes["_type"] != "Error" {
					p.API.LogDebug("Setting MM and OP user id pair: ", mmUserID, opUserID)
					_ = p.API.KVSet(mmUserID, opUserID)
					post = getCreatePostMsg(user.UserId, channelID, "Hello "+opJSONRes["firstName"]+" :)")
					var attachmentMap map[string]interface{}
					_ = json.Unmarshal([]byte(GetAttachmentJSON(pluginURL)), &attachmentMap)
					post.SetProps(attachmentMap)
				} else {
					p.API.LogError(opJSONRes["errorIdentifier"] + " " + opJSONRes["message"])
					post = getCreatePostMsg(user.UserId, channelID, opJSONRes["message"])
				}
			} else {
				p.API.LogError("OpenProject login failed: ", err)
				post = getCreatePostMsg(user.UserId, channelID, messages.OpAuthFailMsg)
			}
		}
		menuPost, _ = p.API.CreatePost(post)
	}
}

func ShowSelProject(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string, action string) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	p.API.LogInfo("apikey: " + APIKeyStr + " opURL: " + OpURLStr)
	user, e := p.API.GetBot(opBot, false)
	if e != nil {
		p.API.LogError(e.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		channelID := jsonBody["channel_id"].(string)
		resp, err := GetProjects(OpURLStr, APIKeyStr)
		var post *model.Post
		if err == nil {
			opResBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			var opJSONRes Projects
			_ = json.Unmarshal(opResBody, &opJSONRes)
			p.API.LogInfo("Projects response from op-mattermost: ", opJSONRes)
			if opJSONRes.Type != "Error" {
				p.API.LogInfo("Projects obtained from OP: ", opJSONRes.Embedded.Elements)
				var options = getOptArrayForProjectElements(opJSONRes.Embedded.Elements)
				post = getUpdatePostMsg(user.UserId, channelID, messages.ProjectSelMsg)
				var attachmentMap map[string]interface{}
				_ = json.Unmarshal(getProjectOptAttachmentJSON(pluginURL, action, options), &attachmentMap)
				post.SetProps(attachmentMap)
			} else {
				p.API.LogError(messages.ProjectFailMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.ProjectFailMsg)
			}
		} else {
			p.API.LogError(messages.ProjectFailMsg, err)
			post = getUpdatePostMsg(user.UserId, channelID, messages.ProjectFailMsg)
		}

		_, _ = p.API.UpdatePost(post)
	}
}

func WPHandler(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	p.API.LogDebug("Request body from project select: ", jsonBody["context"])
	submission := jsonBody["context"].(map[string]interface{})
	var action string
	var selectedOption []string
	for key, value := range submission {
		switch key {
		case "action":
			action = value.(string)
			p.API.LogInfo("action: " + action)
		case "selected_option":
			selectedOption = strings.Split(value.(string), "opt")
			projectID = selectedOption[1]
			p.API.LogInfo("selected option: " + projectID)
		}
	}
	switch action {
	case "showSelWP":
		showSelWP(p, w, jsonBody)
	case "createWP":
		createWP(p, w, jsonBody)
	default:
		http.NotFound(w, r)
	}
}

func showSelWP(p plugin.MattermostPlugin, w http.ResponseWriter, jsonBody map[string]interface{}) {
	user, e := p.API.GetBot(opBot, false)
	if e != nil {
		p.API.LogError(e.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		channelID := jsonBody["channel_id"].(string)
		resp, err := GetWPsForProject(projectID, OpURLStr, APIKeyStr)
		var post *model.Post
		if err == nil {
			opResBody, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			var opJSONRes WorkPackages
			_ = json.Unmarshal(opResBody, &opJSONRes)
			p.API.LogInfo("Work packages response from op-mattermost: ", opJSONRes)
			if opJSONRes.Type != "Error" {
				p.API.LogInfo("Work packages obtained from OP: ", opJSONRes.Embedded.Elements)
				var options = getOptArrayForWPElements(opJSONRes.Embedded.Elements)
				post = getUpdatePostMsg(user.UserId, channelID, messages.WPSelMsg)
				var attachmentMap map[string]interface{}
				_ = json.Unmarshal(getWPOptAttachmentJSON(pluginURL, "showTimeLogDlg", options), &attachmentMap)
				post.SetProps(attachmentMap)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(200)
				_ = respHTTP.Write(w)
			} else {
				p.API.LogError(messages.WPFailMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.WPFailMsg)
			}
		} else {
			p.API.LogError("Failed to fetch work packages from OpenProject: ", err)
			post = getUpdatePostMsg(user.UserId, channelID, messages.WPFailMsg)
		}
		_, _ = p.API.UpdatePost(post)
	}
}

func createWP(p plugin.MattermostPlugin, w http.ResponseWriter, jsonBody map[string]interface{}) {
	user, e := p.API.GetBot(opBot, false)
	if e != nil {
		p.API.LogError(e.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		channelID := jsonBody["channel_id"].(string)
		resp, err := GetTypes(OpURLStr, APIKeyStr)
		var post *model.Post
		if err == nil {
			opResBody, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			var opJSONRes Types
			_ = json.Unmarshal(opResBody, &opJSONRes)
			p.API.LogInfo("Types response from OP: ", opJSONRes)
			if opJSONRes.Type != "Error" {
				var typesOptions = getOptArrayForTypes(opJSONRes.Embedded.Elements)
				resp1, err1 := GetAvailableAssignees(OpURLStr, APIKeyStr, projectID)
				if err1 == nil {
					opResBody1, _ := io.ReadAll(resp1.Body)
					defer resp1.Body.Close()
					var opJSONRes1 AvailableAssignees
					_ = json.Unmarshal(opResBody1, &opJSONRes1)
					p.API.LogInfo("Available assignees response from OP: ", opJSONRes1)
					var availableAssigneesOptions = getOptArrayForAvailableAssignees(opJSONRes1.Embedded.Elements)
					triggerID := jsonBody["trigger_id"].(string)
					openCreateWPDialog(p, triggerID, pluginURL, typesOptions, availableAssigneesOptions)
					post = getUpdatePostMsg(user.UserId, channelID, "Opening WP create dialog...")
					_, _ = p.API.UpdatePost(post)
				} else {
					p.API.LogError(messages.AssigneeFailMsg)
					post = getUpdatePostMsg(user.UserId, channelID, messages.AssigneeFailMsg)
					_, _ = p.API.UpdatePost(post)
				}
			} else {
				p.API.LogError(messages.TypesFailMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.TypesFailMsg)
				_, _ = p.API.UpdatePost(post)
			}
		}
	}
}

func SaveWP(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	channelID := jsonBody["channel_id"].(string)
	p.API.LogInfo("Submission data: ", jsonBody)
	dialogCancelled := jsonBody["cancelled"].(bool)
	user, err := p.API.GetBot(opBot, false)
	if err != nil {
		p.API.LogError(err.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		var post *model.Post
		if dialogCancelled {
			p.API.LogInfo("Save WP dialog cancelled by user.")
			post = getUpdatePostMsg(user.UserId, channelID, messages.DlgCancelMsg)
		} else {
			submission := jsonBody["submission"].(map[string]interface{})
			p.API.LogInfo("WP submission data: ", submission)
			wpJSON, _ := GetWPBodyJSON(submission)
			p.API.LogDebug("WP body JSON: ", string(wpJSON))
			notify := strconv.FormatBool(submission["notify"].(bool))
			resp, err := PostWP(wpJSON, OpURLStr, APIKeyStr, notify)
			if err == nil {
				switch resp.StatusCode {
				case 201:
					p.API.LogInfo(messages.SaveWPSuccessMsg)
					post = getUpdatePostMsg(user.UserId, channelID, messages.SaveWPSuccessMsg)
				case 403:
					p.API.LogError(messages.WPCreateForbiddenMsg)
					post = getUpdatePostMsg(user.UserId, channelID, messages.WPCreateForbiddenMsg)
				case 404:
					p.API.LogError(messages.WPNotExist)
					post = getUpdatePostMsg(user.UserId, channelID, messages.WPNotExist)
				case 422:
					p.API.LogError(messages.WPTypeErrMsg)
					post = getUpdatePostMsg(user.UserId, channelID, messages.WPTypeErrMsg)
				default:
					p.API.LogError(messages.UnknownStatusCode)
					post = getUpdatePostMsg(user.UserId, channelID, messages.UnknownStatusCode)
				}
			} else {
				p.API.LogError(messages.GenericErrMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.GenericErrMsg)
			}
			defer resp.Body.Close()
		}
		_, _ = p.API.UpdatePost(post)
	}
}

func LoadTimeLogDlg(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	triggerID := jsonBody["trigger_id"].(string)
	submission := jsonBody["context"].(map[string]interface{})
	channelID := jsonBody["channel_id"].(string)
	var action string
	var selectedOption []string
	for key, value := range submission {
		switch key {
		case "action":
			action = value.(string)
			p.API.LogInfo("action: " + action)
		case "selected_option":
			selectedOption = strings.Split(value.(string), "|:-")
			timeLogID = selectedOption[1]
			p.API.LogInfo("selected option: " + timeLogID)
		}
	}
	switch action {
	case "showTimeLogDlg":
		user, e := p.API.GetBot(opBot, false)
		if e != nil {
			p.API.LogError(e.Error())
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(502)
			_ = respHTTP.Write(w)
		} else {
			var timeEntriesBody TimeEntriesBody
			timeEntriesBody.Links.WorkPackage.Href = "/api/v3/work_packages/" + timeLogID
			p.API.LogDebug("Time entries body: ", timeEntriesBody)
			timeEntriesBodyJSON, _ := json.Marshal(timeEntriesBody)
			resp, err := PostTimeEntriesForm(timeEntriesBodyJSON, OpURLStr, APIKeyStr)
			var post *model.Post
			if err == nil {
				opResBody, _ := io.ReadAll(resp.Body)
				defer resp.Body.Close()
				var opJSONRes TimeEntries
				_ = json.Unmarshal(opResBody, &opJSONRes)
				p.API.LogDebug("Time entries response from OpenProject: ", opJSONRes)
				if opJSONRes.Type != "Error" {
					var options = getOptArrayForAllowedValues(opJSONRes.Embedded.Schema.Activity.Embedded.AllowedValues)
					openLogTimeDialog(p, triggerID, pluginURL, options)
					post = getUpdatePostMsg(user.UserId, channelID, "Opening time log dialog...")
					_, _ = p.API.UpdatePost(post)
				} else {
					p.API.LogError(messages.ActivityFailMsg)
					post = getUpdatePostMsg(user.UserId, channelID, messages.ActivityFailMsg)
					_, _ = p.API.UpdatePost(post)
				}
			} else {
				p.API.LogError(messages.ActivityFailMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.ActivityFailMsg)
				_, _ = p.API.UpdatePost(post)
			}
		}
	case "cnfDelWP":
		cnfDelWP(p, w, channelID)
	default:
		http.NotFound(w, r)
	}
}

func openCreateWPDialog(p plugin.MattermostPlugin, triggerID string, pluginURL string, types []*model.PostActionOptions, assignees []*model.PostActionOptions) {
	p.API.LogInfo("Types from op-mattermost: ", types)
	p.API.LogInfo("Assignees from op-mattermost: ", assignees)
	p.API.LogDebug("Trigger ID for log time dialog: ", triggerID)
	err := p.API.OpenInteractiveDialog(GetWPCreateDlg(pluginURL, triggerID, types, assignees))
	if err != nil {
		p.API.LogError("Error creating create WP dialog", err)
	}
}

func openLogTimeDialog(p plugin.MattermostPlugin, triggerID string, pluginURL string, options []*model.PostActionOptions) {
	p.API.LogInfo("Activities from op-mattermost: ", options)
	p.API.LogDebug("Trigger ID for log time dialog: ", triggerID)
	err := p.API.OpenInteractiveDialog(GetLogTimeDlg(pluginURL, triggerID, options))
	if err != nil {
		p.API.LogError("Error creating log time dialog", err)
	}
}

func GetTimeLog(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	user, e := p.API.GetBot(opBot, false)
	if e != nil {
		p.API.LogError(e.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		channelID := jsonBody["channel_id"].(string)
		resp, err := GetTimeEntries(OpURLStr, APIKeyStr)
		var post *model.Post
		if err == nil {
			opResBody, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			var opJSONRes TimeEntryList
			_ = json.Unmarshal(opResBody, &opJSONRes)
			p.API.LogDebug("Time entries response from OpenProject: ", opJSONRes)
			if opJSONRes.Type != "Error" {
				var timeLogs = getOptArrayForTimeEntries(opJSONRes.Embedded.Elements)
				p.API.LogInfo("Time entries from op-mattermost: ", timeLogs)
				post = getUpdatePostMsg(user.UserId, channelID, timeLogs)
				_, _ = p.API.UpdatePost(post)
			} else {
				p.API.LogError(messages.TimeEntryFailMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.TimeEntryFailMsg)
				_, _ = p.API.UpdatePost(post)
			}
		} else {
			p.API.LogError(messages.TimeEntryFailMsg)
			post = getUpdatePostMsg(user.UserId, channelID, messages.TimeEntryFailMsg)
			_, _ = p.API.UpdatePost(post)
		}
	}
}

func getOptArrayForTimeEntries(elements []TimeElement) string {
	var tableTxt string
	if len(elements) != 0 {
		tableTxt = "#### Time entries logged\n"
		tableTxt += "| User | Spent On | Project | Work Package | Activity | Logged Time | Billed Time | Comment |\n"
		tableTxt += "|:---------|:---------|:--------|:-------------|:---------|:------------|:------------|:--------|\n"
		for _, element := range elements {
			tableTxt += "| " + element.Links.User.Title + " | "
			d, _ := duration.ParseISO8601(element.Hours)
			loggedTime := convDurationToHoursMin(d)
			billedHours := convHoursToHoursMin(element.CustomField)
			tableTxt += element.SpentOn + " | "
			tableTxt += element.Links.Project.Title + " | "
			tableTxt += element.Links.WorkPackage.Title + " | "
			tableTxt += element.Links.Activity.Title + " | "
			tableTxt += loggedTime + " | "
			tableTxt += billedHours + " | "
			tableTxt += strings.ReplaceAll(element.Comment.Raw, "/\n/g", " ") + " |\n"
		}
	} else {
		tableTxt = "Couldn't find time entries logged by you :confused: Try logging time using `/op`"
	}
	return tableTxt
}

func OpenAuthDialog(p plugin.MattermostPlugin, triggerID string, pluginURL string, logoURL string) {
	_ = p.API.OpenInteractiveDialog(GetOpAuthDlg(pluginURL, triggerID, logoURL))
}

func Logout(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	mmUserID := jsonBody["user_id"].(string)
	p.API.LogInfo("Deleting op login for mm user id: " + mmUserID)
	err := p.API.KVDelete(mmUserID)
	if err == nil {
		user, e := p.API.GetBot(opBot, false)
		if e != nil {
			p.API.LogError(e.Error())
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(502)
			_ = respHTTP.Write(w)
		} else {
			post := getUpdatePostMsg(user.UserId, jsonBody["channel_id"].(string), messages.ByeMsg)
			_, _ = p.API.UpdatePost(post)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
		}
	} else {
		p.API.LogError(" Error deleting mmUserID", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = respHTTP.Write(w)
}

func HandleSubmission(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	p.API.LogInfo("Submission data: ", jsonBody)
	dialogCancelled := jsonBody["cancelled"].(bool)
	user, err := p.API.GetBot(opBot, false)
	if err != nil {
		p.API.LogError(err.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		channelID := jsonBody["channel_id"].(string)
		var post *model.Post
		if dialogCancelled {
			p.API.LogInfo("Log time dialog cancelled by user.")
			post = getUpdatePostMsg(user.UserId, jsonBody["channel_id"].(string), messages.DlgCancelMsg)
		} else {
			submission := jsonBody["submission"].(map[string]interface{})
			if checkDate(submission["spent_on"].(string)) {
				billableHours := submission["billable_hours"].(string)
				loggedHours := submission["spent_hours"].(string)
				if checkHours(billableHours, loggedHours) {
					timeEntriesBodyJSON, _ := GetTimeEntriesBodyJSON(submission, loggedHours, billableHours)
					resp, err := PostTimeEntry(timeEntriesBodyJSON, OpURLStr, APIKeyStr)
					p.API.LogDebug("Time entries body JSON: ", string(timeEntriesBodyJSON))
					if err == nil {
						opResBody, _ := io.ReadAll(resp.Body)
						p.API.LogDebug("Time entries response body from OpenProject: ", string(opResBody))
						defer resp.Body.Close()
						var opJSONRes TimeElement
						_ = json.Unmarshal(opResBody, &opJSONRes)
						p.API.LogDebug("Time entries response from OpenProject: ", opJSONRes)
						if opJSONRes.Type != "Error" {
							p.API.LogInfo("Time logged. Save response: ")
							post = getUpdatePostMsg(user.UserId, channelID, "Time entry ID - "+strconv.Itoa(opJSONRes.ID)+messages.LogTimeSuccessMsg)
						} else {
							p.API.LogError(messages.TimeEntrySaveFailMsg)
							post = getUpdatePostMsg(user.UserId, channelID, messages.TimeEntrySaveFailMsg)
						}
					} else {
						p.API.LogError(messages.TimeEntrySaveFailMsg)
						post = getUpdatePostMsg(user.UserId, channelID, messages.TimeEntrySaveFailMsg)
					}
				} else {
					p.API.LogInfo("Billable hours incorrect: ", billableHours)
					post = getUpdatePostMsg(user.UserId, channelID, messages.BillableHourMsg)
				}
			} else {
				p.API.LogInfo("Date incorrect: ", jsonBody["spent_on"])
				post = getUpdatePostMsg(user.UserId, channelID, messages.DateIncorrectMsg)
			}
		}
		_, _ = p.API.UpdatePost(post)
	}
}

func DeleteWorkPackage(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	submission := jsonBody["context"].(map[string]interface{})
	channelID := jsonBody["channel_id"].(string)
	var action string
	var selectedOption []string
	for key, value := range submission {
		switch key {
		case "action":
			action = value.(string)
			p.API.LogInfo("action: " + action)
		case "selected_option":
			selectedOption = strings.Split(value.(string), "|:-")
			wpEntry = selectedOption[0]
			wpID = selectedOption[1]
			p.API.LogInfo("selected option: " + wpID)
		}
	}
	switch action {
	case "delWP":
		delWP(p, w, wpID, channelID)
	case "cnfDelWP":
		cnfDelWP(p, w, channelID)
	default:
		showDelWPSel(p, w, channelID)
	}
}

func delWP(p plugin.MattermostPlugin, w http.ResponseWriter, wpID string, channelID string) {
	p.API.LogDebug("Deleting WP with ID: " + wpID)
	resp, err := DelWP(OpURLStr, APIKeyStr, wpID)
	user, e := p.API.GetBot(opBot, false)
	if e != nil {
		p.API.LogError(e.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		var post *model.Post
		if err == nil {
			p.API.LogDebug("Work package delete response from OpenProject: ", resp.StatusCode)
			switch resp.StatusCode {
			case 204:
				p.API.LogInfo(messages.WPLogDelMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.WPLogDelMsg)
			case 403:
				p.API.LogError(messages.InsufficientPrivMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.InsufficientPrivMsg)
			case 404:
				p.API.LogError(messages.TimeEntryNotExist)
				post = getUpdatePostMsg(user.UserId, channelID, messages.TimeEntryNotExist)
			default:
				p.API.LogError(messages.TimeLogDelErrMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.TimeLogDelErrMsg)
			}
		} else {
			p.API.LogError(messages.WPDelErrMsg)
			post = getUpdatePostMsg(user.UserId, channelID, messages.WPDelErrMsg)
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(resp.StatusCode)
		_ = respHTTP.Write(w)
		_, _ = p.API.UpdatePost(post)
	}
}

func showDelWPSel(p plugin.MattermostPlugin, w http.ResponseWriter, channelID string) {
	resp, err := GetWorkPackages(OpURLStr, APIKeyStr)
	user, e := p.API.GetBot(opBot, false)
	if e != nil {
		p.API.LogError(e.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		var post *model.Post
		if err == nil {
			opResBody, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			var opJSONRes WorkPackages
			_ = json.Unmarshal(opResBody, &opJSONRes)
			p.API.LogDebug("Work package list response from OpenProject: ", opJSONRes)
			if opJSONRes.Type != "Error" {
				var attachmentMap map[string]interface{}
				var options = getOptArrayForWPElements(opJSONRes.Embedded.Elements)
				p.API.LogInfo("Work packages KV : ", options)
				post = getUpdatePostMsg(user.UserId, channelID, messages.WPLogSelMsg)
				_ = json.Unmarshal(getWPOptJSON(pluginURL, "cnfDelWP", options), &attachmentMap)
				post.SetProps(attachmentMap)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(200)
				_ = respHTTP.Write(w)
			} else {
				p.API.LogError(messages.WPFetchFailMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.WPFetchFailMsg)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			p.API.LogError(messages.WPFetchFailMsg)
			post = getUpdatePostMsg(user.UserId, channelID, messages.WPFetchFailMsg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = p.API.UpdatePost(post)
	}
}

func cnfDelWP(p plugin.MattermostPlugin, w http.ResponseWriter, channelID string) {
	var attachmentMap map[string]interface{}
	var post *model.Post
	user, err := p.API.GetBot(opBot, false)
	if err != nil {
		p.API.LogError(err.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		post = getUpdatePostMsg(user.UserId, channelID, messages.CnfWPLogMsg+"\n"+wpEntry)
		_ = json.Unmarshal(getCnfDelBtnJSON(pluginURL+"/delWP", "delWP"), &attachmentMap)
		post.SetProps(attachmentMap)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_ = respHTTP.Write(w)
		_, _ = p.API.UpdatePost(post)
	}
}

func DeleteTimeLog(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	submission := jsonBody["context"].(map[string]interface{})
	channelID := jsonBody["channel_id"].(string)
	var action string
	var selectedOption []string
	for key, value := range submission {
		switch key {
		case "action":
			action = value.(string)
			p.API.LogInfo("action: " + action)
		case "selected_option":
			selectedOption = strings.Split(value.(string), "|:-")
			timeEntry = selectedOption[0]
			timeLogID = selectedOption[1]
			p.API.LogInfo("selected option: " + timeLogID)
		}
	}
	switch action {
	case "delSelTimeLog":
		delTimeLog(p, w, timeLogID, channelID)
	case "cnfDelTimeLog":
		cnfDelTimeLog(p, w, channelID)
	default:
		showTimeLogSel(p, w, channelID)
	}
}

func cnfDelTimeLog(p plugin.MattermostPlugin, w http.ResponseWriter, channelID string) {
	user, err := p.API.GetBot(opBot, false)
	if err != nil {
		p.API.LogError(err.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		var attachmentMap map[string]interface{}
		post := getUpdatePostMsg(user.UserId, channelID, messages.CnfDelTimeLogMsg+"\n"+timeEntry)
		_ = json.Unmarshal(getCnfDelBtnJSON(pluginURL+"/delTimeLog", "delSelTimeLog"), &attachmentMap)
		post.SetProps(attachmentMap)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_ = respHTTP.Write(w)
		_, _ = p.API.UpdatePost(post)
	}
}

func delTimeLog(p plugin.MattermostPlugin, w http.ResponseWriter, timeLogID string, channelID string) {
	p.API.LogDebug("Deleting time log with ID: " + timeLogID)
	resp, err := DelTimeLog(OpURLStr, APIKeyStr, timeLogID)
	user, e := p.API.GetBot(opBot, false)
	if e != nil {
		p.API.LogError(e.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		var post *model.Post
		if err == nil {
			p.API.LogDebug("Time log delete response from OpenProject: ", resp.StatusCode)
			switch resp.StatusCode {
			case 204:
				var attachmentMap map[string]interface{}
				post = getUpdatePostMsg(user.UserId, channelID, messages.TimeLogDelMsg)
				_ = json.Unmarshal(getTimeLogDelMsgJSON(pluginURL), &attachmentMap)
				post.SetProps(attachmentMap)
			case 403:
				p.API.LogError(messages.InsufficientPrivMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.InsufficientPrivMsg)
			case 404:
				p.API.LogError(messages.TimeEntryNotExist)
				post = getUpdatePostMsg(user.UserId, channelID, messages.TimeEntryNotExist)
			default:
				p.API.LogError(messages.TimeLogDelErrMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.TimeLogDelErrMsg)
			}
		} else {
			p.API.LogError(messages.TimeLogDelErrMsg)
			post = getUpdatePostMsg(user.UserId, channelID, messages.TimeLogDelErrMsg)
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(resp.StatusCode)
		_ = respHTTP.Write(w)
		_, _ = p.API.UpdatePost(post)
	}
}

func showTimeLogSel(p plugin.MattermostPlugin, w http.ResponseWriter, channelID string) {
	resp, err := GetTimeEntries(OpURLStr, APIKeyStr)
	user, e := p.API.GetBot(opBot, false)
	if e != nil {
		p.API.LogError(e.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		var post *model.Post
		if err == nil {
			opResBody, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			var opJSONRes TimeEntryList
			_ = json.Unmarshal(opResBody, &opJSONRes)
			p.API.LogDebug("Time entries response from OpenProject: ", opJSONRes)
			if opJSONRes.Type != "Error" {
				var attachmentMap map[string]interface{}
				var options = getOptArrayForTimeLogElements(opJSONRes.Embedded.Elements)
				p.API.LogInfo("Time entries KV : ", options)
				post = getUpdatePostMsg(user.UserId, channelID, messages.TimeLogSelMsg)
				_ = json.Unmarshal(getTimeLogOptJSON(pluginURL, "cnfDelTimeLog", options), &attachmentMap)
				post.SetProps(attachmentMap)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(200)
				_ = respHTTP.Write(w)
			} else {
				p.API.LogError(messages.TimeEntryFailMsg)
				post = getUpdatePostMsg(user.UserId, channelID, messages.TimeEntryFailMsg)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			p.API.LogError(messages.TimeEntryFailMsg)
			post = getUpdatePostMsg(user.UserId, channelID, messages.TimeEntryFailMsg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = p.API.UpdatePost(post)
	}
}

func NotificationSubscribe(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	subscribedChannelID = jsonBody["channel_id"].(string)
	user, err := p.API.GetBot(opBot, false)
	if err != nil {
		p.API.LogError(err.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		post := getUpdatePostMsg(user.UserId, subscribedChannelID, messages.SubscribeMsg)
		_, _ = p.API.UpdatePost(post)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write(body)
		_ = respHTTP.Write(w)
	}
}

func NotifyChannel(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	p.API.LogDebug("Request body from OpenProject Webhooks: ", string(body))
	defer r.Body.Close()
	var opJSONRes Notification
	user, err := p.API.GetBot(opBot, false)
	if err != nil {
		p.API.LogError(err.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(502)
		_ = respHTTP.Write(w)
	} else {
		_ = json.Unmarshal(body, &opJSONRes)
		if subscribedChannelID == "" {
			p.API.LogDebug("Not subscribed to OpenProject notification webhook")
		} else {
			notificationType := strings.Split(opJSONRes.Action, ":")[0]
			var post *model.Post
			notificationMsg := "**OpenProject notification** "
			switch notificationType {
			case "time_entry":
				notificationMsg += opJSONRes.Action + " by " + opJSONRes.TimeEntry.Links.User.Title + " for " + opJSONRes.TimeEntry.Links.WorkPackage.Title + " in " + opJSONRes.TimeEntry.Links.Project.Title
			case "project":
				notificationMsg += opJSONRes.Action + " : " + opJSONRes.Project.Links.Self.Title
			case "work_package":
				notificationMsg += opJSONRes.Action + " : " + opJSONRes.WorkPackage.Links.Self.Title
			case "attachment":
				notificationMsg += opJSONRes.Action + " : " + opJSONRes.Attachment.Links.Self.Title
			default:
				notificationMsg += opJSONRes.Action
			}
			post = getUpdatePostMsg(user.UserId, subscribedChannelID, notificationMsg)
			_, _ = p.API.UpdatePost(post)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write(body)
		_ = respHTTP.Write(w)
	}
}
