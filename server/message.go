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

const (
	OpAuthFailMsg        = "OpenProject authentication failed. Please try again."
	ProjectSelMsg        = "*Please select a project*"
	ProjectFailMsg       = "Failed to fetch projects from OpenProject"
	WPSelMsg             = "*Please select a work package*"
	WPFailMsg            = "Failed to fetch work packages from OpenProject"
	ActivityFailMsg      = "Failed to fetch activities from OpenProject"
	TimeEntryFailMsg     = "Failed to fetch time entries from OpenProject"
	ByeMsg               = "Logged out.  :wave:"
	BillableHourMsg      = "**It seems that billable hours was incorrect :thinking: Please note billable hours should be less than or equal to logged hours. **"
	TimeEntrySaveFailMsg = "Failed to save time entry in OpenProject"
	DateIncorrectMsg     = "**It seems that date was incorrect :thinking: Please enter a date within last one year and in YYYY-MM-DD format. **"
	DlgCancelMsg         = "** If you would like to try again then, `/op` **"
	LogTimeSuccessMsg    = "\n**Time logged! You are awesome :sunglasses: **\n To view time logged try `/op"
	TimeLogSelMsg        = "*Please select a time log*"
	CnfDelTimeLogMsg     = "**Confirm time log deletion?**"
	TimeLogDelMsg        = "**Time log deleted!**"
	TimeLogDelErrMsg     = "**That didn't work :pensive: Couldn't delete time log\n Please try again...`/op`**"
	InsufficientMsg      = "**You don't have sufficient privileges to do that :pensive: **"
	TimeEntryNotExist    = "** Time entry does not exist or you don't have sufficient privileges to see it :pensive: **"
	AssigneeFailMsg      = "**That didn't work :pensive: Couldn't to fetch available assignees from OP**"
	TypesFailMsg         = "**That didn't work :pensive: Couldn't to fetch types from OP**"
	WPCreateForbiddenMsg = "**It seems that you don't have permission to create work package for this project :confused: **"
	WPTypeErrMsg         = "**Work package type is not set to one of the allowed values. Couldn't create work package :pensive: **"
	GenericErrMsg        = "** Unknown error occurred :pensive: Can you please try again? **"
	SaveWPSuccessMsg     = "\n**Work package created! You are awesome :sunglasses: **\n To log time for a work package try `/op`"
	WPNotExist           = "**Work package does not exist or you don't have sufficient privileges to see it :pensive: **"
	WPLogDelMsg          = "\n**Work package deleted!**"
	WPDelErrMsg          = "**That didn't work :pensive: Couldn't delete work package\n Please try again... `/op`**"
	WPLogSelMsg          = "*Please select a work package*"
	WPFetchFailMsg       = "**That didn't work :pensive: Couldn't fetch work packages from OP**"
	CnfWPLogMsg          = "**Confirm work package deletion?**"
	UnknownStatusCode    = "**Unknown status code - error occurred** :pensive:"
	SubscribeMsg         = "**Subscribed to OpenProject notifications :sunglasses: **"
)
