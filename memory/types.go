package memory

import "syscall"

type Process struct {
	Name            string `json:"name"`
	PID             uint32 `json:"pid"`
	MainModulePath  string `json:"mainModulePath"`
	StartingTime    string `json:"openTime"`
	SessionID       uint32 `json:"sessionID"`
	SessionUserName string `json:"sessionUserName"`
	UserSID         string `json:"sessionUserSID"`
	UserLastLogin   string `json:"sessionLoginTime"`
}

type ProcessTime struct {
	CreationTime syscall.Filetime
	ExitTime     syscall.Filetime
	KernelTime   syscall.Filetime
	UserTime     syscall.Filetime
}
