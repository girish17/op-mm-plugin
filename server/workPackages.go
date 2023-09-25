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

type WorkPackages struct {
	Type     string `json:"_type"`
	Links    Links  `json:"_links"`
	Total    int    `json:"total"`
	Count    int    `json:"count"`
	Embedded struct {
		Elements []Element `json:"elements"`
	} `json:"_embedded"`
}

type WorkPackagePostBody struct {
	Links struct {
		Type struct {
			Href string `json:"href"`
		} `json:"type"`
		Project struct {
			Href string `json:"href"`
		} `json:"project"`
	} `json:"_links"`
	Subject  string `json:"subject"`
	Assignee struct {
		Href  string `json:"href"`
		Title string `json:"title"`
	} `json:"assignee"`
}
