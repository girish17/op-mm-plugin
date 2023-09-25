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

type Types struct {
	Links    Links  `json:"_links"`
	Type     string `json:"_type"`
	Total    int64  `json:"total"`
	Count    int64  `json:"count"`
	Embedded struct {
		Elements []TypeElement `json:"elements"`
	} `json:"_embedded"`
}

type TypeElement struct {
	Type        string `json:"_type"`
	Links       Links  `json:"_links"`
	Color       string `json:"color"`
	IsDefault   bool   `json:"isDefault"`
	IsMilestone bool   `json:"isMilestone"`
	Position    int64  `json:"position"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
