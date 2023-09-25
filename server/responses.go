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

type Option struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

type Context struct {
	Action string `json:"action"`
}

type Integration struct {
	URL     string  `json:"url"`
	Context Context `json:"context"`
}

type Action struct {
	Name        string      `json:"name"`
	Integration Integration `json:"integration"`
	Type        string      `json:"type"`
	Options     []Option    `json:"options"`
}

type Attachment struct {
	Actions []Action `json:"actions"`
}

type OptAttachments struct {
	Attachments []Attachment `json:"attachments"`
}
