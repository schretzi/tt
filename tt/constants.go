package tt

const (
	ErrorString string = "%s %+v\n"
	CharTrack   string = " ▶"
	CharFinish  string = " ■"
	CharErase   string = " ◀"
	CharError   string = " ▲"
	CharInfo    string = " ●"
	CharMore    string = " ◆"
)

const (
	FlagDebug string = "debug"
)

// Erros Messages Common

const (
	errorUnambiguously string = "filter does not resolve to exactly one Row, list the records and use ID"
)

// Error Messages Client
const (
	errorClientId   string = "missing Client ID"
	errorClientName string = "missing Client Name"
)

const (
	errorProjectId            string = "missing Client ID"
	errorProjectName          string = "missing Project Name"
	errorProjectClientMissing string = "Client is needed for a Project to create"
)

const (
	msgCreateClient    string = "Created Client:"
	msgModifiedClient  string = "Modified Client:"
	msgActive          string = "Active"
	msgInactive        string = "Inactive"
	msgCreateProject   string = "Created Project:"
	msgModifiedProject string = "Create Project:"
)

const (
	cfgDefaultClient      string = "client.default"
	cfgDefaultProject     string = "project.default"
	cfgDefaultTask        string = "task.default"
	cfgReportDefaultRange string = "report.default"

	cfgDebug              string = "debug"
	cfgFirstWeekDayMonday string = "firstWeekDayMonday"
	cfgSketchbarPath      string = "sketchybar.path"
)
