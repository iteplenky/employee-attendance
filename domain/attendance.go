package domain

type AttendanceEvent struct {
	ID            int    `json:"id"`
	IIN           string `json:"emp_id"`
	PunchTime     string `json:"punch_time"`
	TerminalAlias string `json:"terminal_alias"`
	Processed     bool   `json:"processed"`
}
