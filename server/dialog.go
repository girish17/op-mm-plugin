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
	"time"

	"github.com/mattermost/mattermost/server/public/model"
)

func GetLogTimeDlg(pluginURL string, triggerID string, options []*model.PostActionOptions) model.OpenDialogRequest {
	return model.OpenDialogRequest{
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
				Default:     "none",
				Placeholder: "Please mention comments if any",
				Optional:    true,
			}, {
				DisplayName: "Select Activity",
				Name:        "activity",
				Type:        "select",
				Placeholder: "Type to search for activity",
				Options:     options,
				Optional:    true,
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
	}
}

func GetOpAuthDlg(pluginURL string, triggerID string, logoURL string) model.OpenDialogRequest {
	return model.OpenDialogRequest{
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
	}
}

func GetWPCreateDlg(pluginURL string, triggerID string, types []*model.PostActionOptions, assignees []*model.PostActionOptions) model.OpenDialogRequest {
	return model.OpenDialogRequest{
		TriggerId: triggerID,
		URL:       pluginURL + "/saveWP",
		Dialog: model.Dialog{
			CallbackId:       "create_wp_dlg",
			Title:            "Create a work package",
			IntroductionText: "Create a work package by providing following details",
			IconURL:          pluginURL + "/public/op_logo.jpg",
			Elements: []model.DialogElement{{
				DisplayName: "Subject",
				Name:        "subject",
				Type:        "text",
				Placeholder: "Name of work package",
				HelpText:    "Please enter date within last one year and in YYYY-MM-DD format",
			}, {
				DisplayName: "Select Type",
				Name:        "type",
				Type:        "select",
				Default:     types[0].Value,
				Placeholder: "Type to search for type",
				Options:     types,
			}, {
				DisplayName: "Assignee",
				Name:        "assignee",
				Type:        "select",
				Placeholder: "Type to search for users",
				Options:     assignees,
				Optional:    true,
			}, {
				DisplayName: "Notify interested users?",
				Name:        "notify",
				Type:        "bool",
				Placeholder: "Send email",
				HelpText:    "Note that this controls notifications for all users interested in changes to the work package (e.g. current user, watchers, author and assignee)",
				Optional:    true,
			}},
			SubmitLabel:    "Create Work Package",
			NotifyOnCancel: true,
		},
	}
}
