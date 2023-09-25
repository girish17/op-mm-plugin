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

type AvailableAssignees struct {
	Links    Links  `json:"_links"`
	Type     string `json:"_type"`
	Total    int64  `json:"total"`
	Count    int64  `json:"count"`
	Embedded struct {
		Elements []AvailableAssigneesElement `json:"elements"`
	} `json:"_embedded"`
}

type AvailableAssigneesElement struct {
	Links     AvailableAssigneeLinks `json:"_links"`
	Type      string                 `json:"_type"`
	ID        int                    `json:"id"`
	Avatar    string                 `json:"avatar"`
	Email     string                 `json:"email"`
	FirstName string                 `json:"firstName"`
	LastName  string                 `json:"lastName"`
	Login     string                 `json:"login"`
	Status    string                 `json:"status"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
}

type AvailableAssigneeLinks struct {
	Delete Lock `json:"delete"`
	Self   Self `json:"self"`
	Lock   Lock `json:"lock"`
}

type Lock struct {
	Href   string `json:"href"`
	Method string `json:"method"`
	Title  string `json:"title"`
}
