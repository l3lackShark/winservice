package memory

import "syscall"

type JSONChanges struct {
	New    []Process `json:"new"`
	Clsoed []Process `json:"closed"`
}

//There is a small chance that the same process ID will be assigned to the different process
type UniqueProcess struct {
	PID          uint32
	CreationTime string
}

type Process struct {
	Name           string `json:"name"`
	PID            uint32 `json:"pid"`
	MainModulePath string `json:"mainModulePath"`
	CreationTime   string `json:"openTime"`
	User           User   `json:"owningUser"`
}

type User struct {
	SessionID uint32 `json:"sessionID"`
	Name      string `json:"name"`
	SID       string `json:"sid"`
	LastLogin string `json:"loginTime"`
}

type ProcessTime struct {
	CreationTime syscall.Filetime
	ExitTime     syscall.Filetime
	KernelTime   syscall.Filetime
	UserTime     syscall.Filetime
}
