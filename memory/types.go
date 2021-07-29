package memory

import "syscall"

type JSONOut struct {
	New    []Process `json:"new"`
	Clsoed []Process `json:"closed"`
}

//There is a small chance that the same process ID will be assigned to the different process
type UniqueProcess struct {
	PID          uint32
	CreationTime string
}

type Process struct {
	Name            string `json:"name"`
	PID             uint32 `json:"pid"`
	MainModulePath  string `json:"mainModulePath"`
	CreationTime    string `json:"openTime"`
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
