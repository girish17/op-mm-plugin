package main

import (
	"net/http"

	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/girish17/op-mm-plugin/server/util"
)

// ServeHTTP demonstrates a plugin that handles HTTP requests.
func (p *Plugin) ServeHTTP(_ *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/opAuth":
		util.OpAuth(p.MattermostPlugin, r, pluginURL)
	case "/createTimeLog":
		util.ShowSelProject(p.MattermostPlugin, r, pluginURL)
	case "/projSel":
		util.WPHandler(p.MattermostPlugin, w, r, pluginURL)
	case "/wpSel":
		util.LoadTimeLogDlg(p.MattermostPlugin, w, r, pluginURL)
	case "/logTime":
		util.HandleSubmission(p.MattermostPlugin, w, r, pluginURL)
	case "/getTimeLog":
		util.GetTimeLog(p.MattermostPlugin, r)
	case "/delTimeLog":
		http.NotFound(w, r)
	case "/createWP":
		http.NotFound(w, r)
	case "/saveWP":
		http.NotFound(w, r)
	case "/delWP":
		util.ShowDelWPSel()
	case "/bye":
		util.Logout(p.MattermostPlugin, w, r)
	default:
		http.NotFound(w, r)
	}
}
