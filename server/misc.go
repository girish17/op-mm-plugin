package main

import (
	"encoding/json"
	"strconv"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

func getProjectOptAttachmentJSON(pluginURL string, action string, options []Option) []byte {
	attachments := OptAttachments{Attachments: []Attachment{
		{
			Actions: []Action{
				{
					Name: "Type to search for a project...",
					Integration: Integration{
						URL: pluginURL + "/projSel",
						Context: Context{
							Action: action,
						},
					},
					Type:    "select",
					Options: options,
				},
				{
					Name: "Cancel Project search",
					Integration: Integration{
						URL: pluginURL + "/bye",
					},
				},
			},
		},
	},
	}
	attachmentsJSON, _ := json.Marshal(attachments)
	return attachmentsJSON
}

func getWPOptAttachmentJSON(pluginURL string, action string, options []Option) []byte {
	attachments := OptAttachments{Attachments: []Attachment{
		{
			Actions: []Action{
				{
					Name: "Type to search for a work package...",
					Integration: Integration{
						URL: pluginURL + "/wpSel",
						Context: Context{
							Action: action,
						},
					},
					Type:    "select",
					Options: options,
				},
				{
					Name: "Cancel WP search",
					Integration: Integration{
						URL: pluginURL + "/createTimeLog",
					},
				},
			},
		},
	},
	}
	attachmentsJSON, _ := json.Marshal(attachments)
	return attachmentsJSON
}

func getCreatePostMsg(userID string, channelID string, msg string) *model.Post {
	var post = &model.Post{
		UserId:    userID,
		ChannelId: channelID,
		Message:   msg,
	}
	return post
}

func getOptArrayForProjectElements(elements []Element) []Option {
	var options []Option
	for _, element := range elements {
		id := strconv.Itoa(element.ID)
		options = append(options, Option{
			Text:  element.Name,
			Value: "opt" + id,
		})
	}
	return options
}

func getOptArrayForWPElements(elements []Element) []Option {
	var options []Option
	for _, element := range elements {
		id := strconv.Itoa(element.ID)
		options = append(options, Option{
			Text:  element.Subject,
			Value: "opt" + id,
		})
	}
	return options
}

func getOptArrayForAllowedValues(allowedValues []AllowedValues) []*model.PostActionOptions {
	var postActionOptions []*model.PostActionOptions
	for _, value := range allowedValues {
		id := strconv.Itoa(value.ID)
		postActionOptions = append(postActionOptions, &model.PostActionOptions{
			Text:  value.Name,
			Value: "opt" + id,
		})
	}
	return postActionOptions
}

func GetAttachmentJSON(pluginURL string) string {
	return `{
			"attachments": [
				  {
					"text": "What would you like me to do?",
					"actions": [
					  {
						"name": "Log time",
						"integration": {
						  "url": "` + pluginURL + `/createTimeLog",
						  "context": {
							"action": "showSelWP"
						  }
						}
					  },
					  {
						"name": "Create Work Package",
						"integration": {
                          "url": "` + pluginURL + `/createWP",
						  "context": {
							"action": "createWP"
						  }
						}
					  },
					  {
						"name": "View time logs",
						"integration": {
					      "url": "` + pluginURL + `/getTimeLog",
						  "context": {
							"action": "getTimeLog"
						  }
						}
					  },
					  {
						"name": "Delete time log",
						"integration": {
                          "url": "` + pluginURL + `/delTimeLog",
						  "context": {
							"action": "delTimeLog"
						  }
						}
					  },
					  {
						"name": "Delete Work Package",
						"integration": {
                         "url": "` + pluginURL + `/delWP",
						  "context": {
							"action": ""
						  }
						}
					  },
					  {
						"name": "Bye :wave:",
						"integration": {
                          "url": "` + pluginURL + `/bye",
						  "context": {
							"action": "bye"
						  }
						}
					  }
					]
				  }
			]
		}`
}

func getUpdatePostMsg(userID string, channelID string, msg string) *model.Post {
	var post = &model.Post{
		Id:        menuPost.Id,
		UserId:    userID,
		ChannelId: channelID,
		Message:   msg,
	}
	return post
}

func setOPStr(p plugin.MattermostPlugin) {
	opURL, _ := p.API.KVGet("opURL")
	apiKey, _ := p.API.KVGet("apiKey")
	OpURLStr = string(opURL)
	APIKeyStr = string(apiKey)
	p.API.LogInfo("opURLStr: " + OpURLStr + " apiKeyStr: " + APIKeyStr)
}

func checkDate() bool {
	return true
}

func checkHours(_ string) bool {
	return true
}
