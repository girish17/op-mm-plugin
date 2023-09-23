package main

type Notification struct {
	Action      string       `json:"action"`
	TimeEntry   TimeEntries  `json:"time_entry"`
	Project     Projects     `json:"project"`
	WorkPackage WorkPackages `json:"work_package"`
	Attachment  Attachment   `json:"attachment"`
}
