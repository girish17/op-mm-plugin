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
	"net/http"

	"github.com/mattermost/mattermost/server/public/plugin"
)

// ServeHTTP demonstrates a plugin that handles HTTP requests.
func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/opAuth":
		OpAuth(p.MattermostPlugin, r, pluginURL)
	case "/createTimeLog":
		ShowSelProject(p.MattermostPlugin, r, pluginURL, "showSelWP")
	case "/projSel":
		WPHandler(p.MattermostPlugin, w, r)
	case "/wpSel":
		LoadTimeLogDlg(p.MattermostPlugin, w, r, pluginURL)
	case "/logTime":
		HandleSubmission(p.MattermostPlugin, w, r)
	case "/getTimeLog":
		GetTimeLog(p.MattermostPlugin, r)
	case "/delTimeLog":
		DeleteTimeLog(p.MattermostPlugin, w, r)
	case "/createWP":
		ShowSelProject(p.MattermostPlugin, r, pluginURL, "createWP")
	case "/saveWP":
		SaveWP(p.MattermostPlugin, r)
	case "/delWP":
		DeleteWorkPackage(p.MattermostPlugin, w, r)
	case "/subscribe":
		NotificationSubscribe(p.MattermostPlugin, w, r)
	case "/notifyChannel":
		NotifyChannel(p.MattermostPlugin, w, r)
	case "/bye":
		Logout(p.MattermostPlugin, w, r)
	default:
		http.NotFound(w, r)
	}
}
