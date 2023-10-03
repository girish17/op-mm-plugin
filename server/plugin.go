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
	"strings"
	"sync"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/pkg/errors"
)

const opCommand = "op"
const opBot = "op-mattermost"

var pluginURL string

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// OnActivate See https://developers.mattermost.com/extend/plugins/server/reference/
func (p *Plugin) OnActivate() error {
	if p.MattermostPlugin.API.GetConfig().ServiceSettings.SiteURL == nil {
		p.MattermostPlugin.API.LogError("SiteURL must be set.")
	}

	if err := p.MattermostPlugin.API.RegisterCommand(createOpCommand(p.GetSiteURL())); err != nil {
		return errors.Wrapf(err, "failed to register %s command", opCommand)
	}

	newBot := model.Bot{
		Username:    opBot,
		DisplayName: opBot,
	}

	if _, err := p.MattermostPlugin.API.CreateBot(&newBot); err != nil {
		return errors.Wrapf(err, "failed to register #{opBot} bot")
	}
	p.MattermostPlugin.API.LogInfo("Deleting all KV pairs")
	_ = p.MattermostPlugin.API.KVDeleteAll()
	return nil
}

func (p *Plugin) OnDeactivate() error {
	if e := p.MattermostPlugin.API.PermanentDeleteBot(opBot); e != nil {
		return errors.Wrapf(e, "failed to permanently delete %s bot", opBot)
	}
	p.MattermostPlugin.API.LogInfo("Deleting all KV pairs")
	_ = p.MattermostPlugin.API.KVDeleteAll()
	return nil
}

//goland:noinspection GoDeprecation
func (p *Plugin) ExecuteCommand(_ *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	siteURL := p.GetSiteURL()
	pluginURL = getPluginURL(siteURL)
	logoURL := getLogoURL(siteURL)
	p.MattermostPlugin.API.LogInfo("Plugin URL:" + pluginURL + " Logo URL: " + logoURL)
	opUserID, _ := p.MattermostPlugin.API.KVGet(args.UserId)
	if opUserID == nil {
		p.MattermostPlugin.API.LogDebug("Creating interactive dialog...")
		OpenAuthDialog(p.MattermostPlugin, args.TriggerId, pluginURL, logoURL)
		resp := &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "opening op auth dialog",
			Username:     opBot,
			IconURL:      logoURL,
		}
		return resp, nil
	}
	p.MattermostPlugin.API.LogDebug("opUserID: ", opUserID)
	apiKeyStr := strings.Split(string(opUserID), " ")
	opURLStr := apiKeyStr[1]
	p.MattermostPlugin.API.LogInfo("Retrieving from KV: opURL - " + opURLStr + " apiKey - " + apiKeyStr[0])

	var cmdResp *model.CommandResponse

	resp, _ := GetUserDetails(OpURLStr, APIKeyStr)
	opResBody, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var opJSONRes map[string]string
	_ = json.Unmarshal(opResBody, &opJSONRes)
	p.MattermostPlugin.API.LogInfo("Hello : ", opJSONRes["firstName"])

	var attachmentMap map[string]interface{}
	_ = json.Unmarshal([]byte(GetAttachmentJSON(pluginURL)), &attachmentMap)

	cmdResp = &model.CommandResponse{
		ResponseType: model.CommandResponseTypeInChannel,
		Text:         "Hello " + opJSONRes["name"] + " :)",
		Username:     opBot,
		IconURL:      logoURL,
		Props:        attachmentMap,
	}
	return cmdResp, nil
}

func (p *Plugin) GetSiteURL() string {
	siteURL := ""
	ptr := p.MattermostPlugin.API.GetConfig().ServiceSettings.SiteURL
	if ptr != nil {
		siteURL = *ptr
	}
	return siteURL
}

func getLogoURL(siteURL string) string {
	return getPluginURL(siteURL) + "/public/op_logo.jpg"
}

func getPluginURL(siteURL string) string {
	return siteURL + "/plugins/com.girishm.info.op-mm-plugin"
}

func createOpCommand(siteURL string) *model.Command {
	return &model.Command{
		Trigger:          opCommand,
		Method:           "POST",
		Username:         opBot,
		IconURL:          getLogoURL(siteURL),
		AutoComplete:     true,
		AutoCompleteDesc: "Invoke OpenProject bot for Mattermost",
		AutoCompleteHint: "",
		DisplayName:      opBot,
		Description:      "OpenProject integration for Mattermost",
		URL:              siteURL,
	}
}

func _(opBot string) *model.Bot {
	return &model.Bot{
		Username:    opBot,
		DisplayName: opBot,
		Description: "OpenProject bot",
	}
}

//goland:noinspection GoDeprecation
// func (p *Plugin) setBotIcon() {
//	bundlePath, err := p.MattermostPlugin.API.GetBundlePath()
//	if err != nil {
//		p.MattermostPlugin.API.LogError("failed to get bundle path", err)
//	}
//
//	profileImage, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "op_logo.svg"))
//	if err != nil {
//		p.MattermostPlugin.API.LogError("failed to read profile image", err)
//	}
//
//	user, err := p.MattermostPlugin.API.GetBot(opBot, false)
//	if err != nil {
//		p.MattermostPlugin.API.LogError("failed to fetch bot user", err)
//	}
//
//	if appErr := p.MattermostPlugin.API.SetProfileImage(user.UserId, profileImage); appErr != nil {
//		p.MattermostPlugin.API.LogError("failed to set profile image", appErr)
//	}
// }
