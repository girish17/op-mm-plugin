package util

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/senseyeio/duration"

	"github.com/girish17/op-mm-plugin/server/types"
)

const opBot = "op-mattermost"

var menuPost *model.Post

var respHTTP http.Response

var client = &http.Client{}

var opURLStr string

var apiKeyStr string

var projectID string

var wpID string

var activityID string

func OpAuth(p plugin.MattermostPlugin, r *http.Request, pluginURL string) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	p.API.LogDebug("Request body from dialog submit: ", jsonBody)
	dialogCancelled := jsonBody["cancelled"].(bool)
	user, _ := p.API.GetUserByUsername(opBot)
	var post *model.Post
	if dialogCancelled {
		p.API.LogInfo("Op Auth Dialog cancelled by user.")
		post = getCreatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.DlgCancelMsg)
	} else {
		submission := jsonBody["submission"].(map[string]interface{})
		mmUserID := jsonBody["user_id"].(string)
		for key, value := range submission {
			p.API.LogInfo("Storing OpenProject auth credentials: " + key + ":" + value.(string))
			_ = p.API.KVSet(key, []byte(value.(string)))
		}
		setOPStr(p)
		opUserID := []byte(apiKeyStr + " " + opURLStr)

		_ = p.API.KVDelete(opURLStr)
		_ = p.API.KVDelete(apiKeyStr)

		resp, err := GetUserDetails(opURLStr, apiKeyStr)
		if err == nil {
			opResBody, _ := io.ReadAll(resp.Body)
			var opJSONRes map[string]string
			_ = json.Unmarshal(opResBody, &opJSONRes)
			p.API.LogDebug("Response from op-mattermost: ", opJSONRes)
			if opJSONRes["_type"] != "Error" {
				p.API.LogDebug("Setting MM and OP user id pair: ", mmUserID, opUserID)
				_ = p.API.KVSet(mmUserID, opUserID)
				post = getCreatePostMsg(user.Id, jsonBody["channel_id"].(string), "Hello "+opJSONRes["firstName"]+" :)")
				var attachmentMap map[string]interface{}
				_ = json.Unmarshal([]byte(GetAttachmentJSON(pluginURL)), &attachmentMap)
				post.SetProps(attachmentMap)
			} else {
				p.API.LogError(opJSONRes["errorIdentifier"] + " " + opJSONRes["message"])
				post = getCreatePostMsg(user.Id, jsonBody["channel_id"].(string), opJSONRes["message"])
			}
		} else {
			p.API.LogError("OpenProject login failed: ", err)
			post = getCreatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.OpAuthFailMsg)
		}
		defer resp.Body.Close()
	}
	menuPost, _ = p.API.CreatePost(post)
}

func ShowSelProject(p plugin.MattermostPlugin, r *http.Request, pluginURL string) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	p.API.LogInfo("apikey: " + apiKeyStr + " opURL: " + opURLStr)
	user, _ := p.API.GetUserByUsername(opBot)
	resp, err := GetProjects(opURLStr, apiKeyStr)
	var post *model.Post
	if err == nil {
		opResBody, _ := io.ReadAll(resp.Body)
		var opJSONRes types.Projects
		_ = json.Unmarshal(opResBody, &opJSONRes)
		p.API.LogInfo("Projects response from op-mattermost: ", opJSONRes)
		if opJSONRes.Type != "Error" {
			p.API.LogInfo("Projects obtained from OP: ", opJSONRes.Embedded.Elements)
			var options = getOptArrayForProjectElements(opJSONRes.Embedded.Elements)
			post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.ProjectSelMsg)
			var attachmentMap map[string]interface{}
			_ = json.Unmarshal(getProjectOptAttachmentJSON(pluginURL, "showSelWP", options), &attachmentMap)
			post.SetProps(attachmentMap)
		} else {
			p.API.LogError(messages.ProjectFailMsg)
			post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.ProjectFailMsg)
		}
	} else {
		p.API.LogError(messages.ProjectFailMsg, err)
		post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.ProjectFailMsg)
	}
	defer resp.Body.Close()
	_, _ = p.API.UpdatePost(post)
}

func WPHandler(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string) {
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
		user, _ := p.API.GetUserByUsername(opBot)
		resp, err := GetWPsForProject(projectID, opURLStr, apiKeyStr)
		var post *model.Post
		if err == nil {
			opResBody, _ := io.ReadAll(resp.Body)
			var opJSONRes types.WorkPackages
			_ = json.Unmarshal(opResBody, &opJSONRes)
			p.API.LogInfo("Work packages response from op-mattermost: ", opJSONRes)
			if opJSONRes.Type != "Error" {
				p.API.LogInfo("Work packages obtained from OP: ", opJSONRes.Embedded.Elements)
				var options = getOptArrayForWPElements(opJSONRes.Embedded.Elements)
				post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.WPSelMsg)
				var attachmentMap map[string]interface{}
				_ = json.Unmarshal(getWPOptAttachmentJSON(pluginURL, "showTimeLogDlg", options), &attachmentMap)
				post.SetProps(attachmentMap)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(200)
				_ = respHTTP.Write(w)
			} else {
				p.API.LogError(messages.WPFailMsg)
				post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.WPFailMsg)
			}
		} else {
			p.API.LogError("Failed to fetch work packages from OpenProject: ", err)
			post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.WPFailMsg)
		}
		defer resp.Body.Close()
		_, _ = p.API.UpdatePost(post)
	case "createWP":
		http.NotFound(w, r)
	default:
		http.NotFound(w, r)
	}
}

func LoadTimeLogDlg(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request, pluginURL string) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	triggerID := jsonBody["trigger_id"].(string)
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
			wpID = selectedOption[1]
			p.API.LogInfo("selected option: " + wpID)
		}
	}
	switch action {
	case "showTimeLogDlg":
		user, _ := p.API.GetUserByUsername(opBot)
		var timeEntriesBody types.TimeEntriesBody
		timeEntriesBody.Links.WorkPackage.Href = "/api/v3/work_packages/" + wpID
		p.API.LogDebug("Time entries body: ", timeEntriesBody)
		timeEntriesBodyJSON, _ := json.Marshal(timeEntriesBody)
		resp, err := PostTimeEntriesForm(timeEntriesBodyJSON, opURLStr, apiKeyStr)
		var post *model.Post
		if err == nil {
			opResBody, _ := io.ReadAll(resp.Body)
			var opJSONRes types.TimeEntries
			_ = json.Unmarshal(opResBody, &opJSONRes)
			p.API.LogDebug("Time entries response from OpenProject: ", opJSONRes)
			if opJSONRes.Type != "Error" {
				var options = getOptArrayForAllowedValues(opJSONRes.Embedded.Schema.Activity.Embedded.AllowedValues)
				openLogTimeDialog(p, triggerID, pluginURL, options)
				post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "Opening time log dialog...")
				_, _ = p.API.UpdatePost(post)
			} else {
				p.API.LogError(messages.ActivityFailMsg)
				post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.ActivityFailMsg)
				_, _ = p.API.UpdatePost(post)
			}
		} else {
			p.API.LogError(messages.ActivityFailMsg)
			post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.ActivityFailMsg)
			_, _ = p.API.UpdatePost(post)
		}
		defer resp.Body.Close()
	case "cnfDelWP":
		http.NotFound(w, r)
	default:
		http.NotFound(w, r)
	}
}

func openLogTimeDialog(p plugin.MattermostPlugin, triggerID string, pluginURL string, options []*model.PostActionOptions) {
	p.API.LogInfo("Activities from op-mattermost: ", options)
	p.API.LogDebug("Trigger ID for log time dialog: ", triggerID)
	dialog := p.API.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: triggerID,
		URL:       pluginURL + "/logTime",
		Dialog: model.Dialog{
			CallbackId:       "log_time_dlg",
			Title:            "Log time for work package",
			IntroductionText: "Please enter details to log time",
			IconURL:          pluginURL + "/public/op_logo.jpg",
			Elements: []model.DialogElement{{
				DisplayName: "Date",
				Name:        "spent_on",
				Type:        "text",
				Default:     time.Now().Format("2006-01-02"),
				Placeholder: "YYYY-MM-DD",
				HelpText:    "Please enter date within last one year and in YYYY-MM-DD format",
			}, {
				DisplayName: "Comment",
				Name:        "comments",
				Type:        "textarea",
				Placeholder: "Please mention comments if any",
				Optional:    true,
			}, {
				DisplayName: "Select Activity",
				Name:        "activity",
				Type:        "select",
				Default:     options[0].Value,
				Placeholder: "Type to search for activity",
				Options:     options,
			}, {
				DisplayName: "Spent hours",
				Name:        "spent_hours",
				Type:        "text",
				Default:     "0.5",
				Placeholder: "hours like 0.5, 1, 3 ...",
				HelpText:    "Please enter spent hours to be logged",
			}, {
				DisplayName: "Billable hours",
				Name:        "billable_hours",
				Type:        "text",
				Default:     "0.0",
				Placeholder: "hours like 0.5, 1, 3 ...",
				HelpText:    "Please ensure billable hours is less than or equal to spent hours",
			}},
			SubmitLabel:    "Log time",
			NotifyOnCancel: true,
		},
	})
	p.API.LogDebug("Dialog object returned: ", dialog)
}

func GetTimeLog(p plugin.MattermostPlugin, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	user, _ := p.API.GetUserByUsername(opBot)
	resp, err := GetTimeEntries(opURLStr, apiKeyStr)
	var post *model.Post
	if err == nil {
		opResBody, _ := io.ReadAll(resp.Body)
		var opJSONRes types.TimeEntryList
		_ = json.Unmarshal(opResBody, &opJSONRes)
		p.API.LogDebug("Time entries response from OpenProject: ", opJSONRes)
		if opJSONRes.Type != "Error" {
			var timeLogs = getOptArrayForTimeEntries(opJSONRes.Embedded.Elements)
			p.API.LogInfo("Time entries from op-mattermost: ", timeLogs)
			post = getCreatePostMsg(user.Id, jsonBody["channel_id"].(string), timeLogs)
			_, _ = p.API.CreatePost(post)
		} else {
			p.API.LogError(messages.TimeEntryFailMsg)
			post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.TimeEntryFailMsg)
			_, _ = p.API.UpdatePost(post)
		}
	} else {
		p.API.LogError(messages.TimeEntryFailMsg)
		post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.TimeEntryFailMsg)
		_, _ = p.API.UpdatePost(post)
	}
	defer resp.Body.Close()
}

func getOptArrayForTimeEntries(elements []types.TimeElement) string {
	var tableTxt string
	if len(elements) != 0 {
		tableTxt = "#### Time entries logged by you\n"
		tableTxt += "| Spent On | Project | Work Package | Activity | Logged Time | Billed Time | Comment |\n"
		tableTxt += "|:---------|:--------|:-------------|:---------|:------------|:------------|:--------|\n"
		for _, element := range elements {
			d, _ := duration.ParseISO8601(element.Hours)
			var loggedTime = ""
			if d.TH != 0 {
				hours := strconv.Itoa(d.TH)
				if d.TH > 1 {
					loggedTime = hours + " hours "
				} else {
					loggedTime = hours + " hour "
				}
			}
			if d.TM != 0 {
				minutes := strconv.Itoa(d.TM)
				if d.TM > 1 {
					loggedTime = loggedTime + minutes + " minutes"
				} else {
					loggedTime = loggedTime + minutes + " minute"
				}
			}
			billedHours := strconv.FormatFloat(element.CustomField1, 'f', 2, 64)
			tableTxt += "| " + element.SpentOn + " | " + element.Links.Project.Title + " | " + element.Links.WorkPackage.Title + " | " + element.Links.Activity.Title + " | " + loggedTime + " | " + billedHours + " hours" + " | " + strings.ReplaceAll(element.Comment.Raw, "/\n/g", " ") + " |\n"
		}
	} else {
		tableTxt = "Couldn't find time entries logged by you :confused: Try logging time using `/op`"
	}
	return tableTxt
}

func ShowDelWPSel() {

}

func OpenAuthDialog(p plugin.MattermostPlugin, triggerID string, pluginURL string, logoURL string) {
	_ = p.API.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: triggerID,
		URL:       pluginURL + "/opAuth",
		Dialog: model.Dialog{
			CallbackId:       "op_auth_dlg",
			Title:            "OpenProject Authentication",
			IntroductionText: "Please enter credentials to log in",
			IconURL:          logoURL,
			Elements: []model.DialogElement{{
				DisplayName: "OpenProject URL",
				Name:        "opURL",
				Type:        "text",
				SubType:     "url",
				Default:     "http://localhost:8080",
				Placeholder: "http://localhost:8080",
				Optional:    false,
				HelpText:    "Please enter the URL of OpenProject server",
			}, {
				DisplayName: "OpenProject api-key",
				Name:        "apiKey",
				Type:        "text",
				SubType:     "password",
				Placeholder: "api-key generated from your account page in OpenProject",
				Optional:    false,
				HelpText:    "api-key can be generated within 'My account' section of OpenProject",
			}},
			SubmitLabel:    "Log in",
			NotifyOnCancel: true,
		},
	})
}

func Logout(p plugin.MattermostPlugin, w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	mmUserID := jsonBody["user_id"].(string)
	p.API.LogInfo("Deleting op login for mm user id: " + mmUserID)
	_ = p.API.KVDelete(mmUserID)

	user, _ := p.API.GetUserByUsername(opBot)
	post := getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.ByeMsg)
	_, _ = p.API.UpdatePost(post)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = respHTTP.Write(w)
}

func HandleSubmission(p plugin.MattermostPlugin, _ http.ResponseWriter, r *http.Request, _ string) {
	body, _ := io.ReadAll(r.Body)
	var jsonBody map[string]interface{}
	_ = json.Unmarshal(body, &jsonBody)
	p.API.LogInfo("Submission data: ", jsonBody)
	dialogCancelled := jsonBody["cancelled"].(bool)
	user, _ := p.API.GetUserByUsername(opBot)
	var post *model.Post
	if dialogCancelled {
		p.API.LogInfo("Log time dialog cancelled by user.")
		post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.DlgCancelMsg)
	} else {
		if checkDate() {
			if checkHours(jsonBody["billable_hours"].(string)) {
				var timeEntriesBody types.TimeEntryPostBody
				timeEntriesBody.Links.Project.Href = "/api/v3/projects/" + projectID
				timeEntriesBody.Links.WorkPackage.Href = "/api/v3/work_packages/" + wpID
				activityID = strings.Split(jsonBody["activity"].(string), "opt")[1]
				timeEntriesBody.Links.Activity.Href = "/api/v3/time_entries/activities/" + activityID
				p.API.LogDebug("Time entries body: ", timeEntriesBody)
				timeEntriesBodyJSON, _ := json.Marshal(timeEntriesBody)
				resp, err := PostTimeEntry(timeEntriesBodyJSON, opURLStr, apiKeyStr)
				if err == nil {
					opResBody, _ := io.ReadAll(resp.Body)
					var opJSONRes types.TimeEntries
					_ = json.Unmarshal(opResBody, &opJSONRes)
					p.API.LogDebug("Time entries response from OpenProject: ", opJSONRes)
					if opJSONRes.Type != "Error" {
						p.API.LogInfo("Time logged. Save response: ")
						post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), "Time entry ID - ")
					} else {
						p.API.LogError(messages.TimeEntrySaveFailMsg)
						post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.TimeEntrySaveFailMsg)
					}
				} else {
					p.API.LogError(messages.TimeEntrySaveFailMsg)
					post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.TimeEntrySaveFailMsg)
				}
				defer resp.Body.Close()
			} else {
				p.API.LogInfo("Billable hours incorrect: ", jsonBody["billable_hours"])
				post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.BillableHourMsg)
			}
		} else {
			p.API.LogInfo("Date incorrect: ", jsonBody["spent_on"])
			post = getUpdatePostMsg(user.Id, jsonBody["channel_id"].(string), messages.DateIncorrectMsg)
		}
	}
	_, _ = p.API.UpdatePost(post)
}
