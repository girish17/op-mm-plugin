package main

var messages = struct {
	OpAuthFailMsg        string
	ProjectSelMsg        string
	ProjectFailMsg       string
	WPSelMsg             string
	WPFailMsg            string
	ActivityFailMsg      string
	TimeEntryFailMsg     string
	ByeMsg               string
	BillableHourMsg      string
	TimeEntrySaveFailMsg string
	DateIncorrectMsg     string
	DlgCancelMsg         string
	NotImplemented       string
	LogTimeSuccessMsg    string
	TimeLogSelMsg        string
	CnfDelTimeLogMsg     string
	TimeLogDelMsg        string
	TimeLogDelErrMsg     string
	InsufficientPrivMsg  string
	TimeEntryNotExist    string
}{
	OpAuthFailMsg:        "OpenProject authentication failed. Please try again.",
	ProjectSelMsg:        "*Please select a project*",
	ProjectFailMsg:       "Failed to fetch projects from OpenProject",
	WPSelMsg:             "*Please select a work package*",
	WPFailMsg:            "Failed to fetch work packages from OpenProject",
	ActivityFailMsg:      "Failed to fetch activities from OpenProject",
	TimeEntryFailMsg:     "Failed to fetch time entries from OpenProject",
	ByeMsg:               "Donate at https://www.paypal.me/girishm17 :wave:",
	BillableHourMsg:      "**It seems that billable hours was incorrect :thinking: Please note billable hours should be less than or equal to logged hours. **",
	TimeEntrySaveFailMsg: "Failed to save time entry in OpenProject",
	DateIncorrectMsg:     "**It seems that date was incorrect :thinking: Please enter a date within last one year and in YYYY-MM-DD format. **",
	DlgCancelMsg:         "** If you would like to try again then, `/op` **",
	NotImplemented:       "Not yet implemented. Donate at https://www.paypal.me/girishm17 :wave:",
	LogTimeSuccessMsg:    "\n**Time logged! You are awesome :sunglasses: **\n To view time logged try `/op",
	TimeLogSelMsg:        "*Please select a time log*",
	CnfDelTimeLogMsg:     "**Confirm time log deletion?**",
	TimeLogDelMsg:        "**Time log deleted!**",
	TimeLogDelErrMsg:     "**That didn't work :pensive: Couldn't delete time log\n Please try again...`/op`**",
	InsufficientPrivMsg:  "**You don't have sufficient privileges to do that :pensive: **",
	TimeEntryNotExist:    "** Time entry does not exist or you don't have sufficient privileges to see it :pensive: **",
}
