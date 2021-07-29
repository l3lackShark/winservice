package memory

import (
	"fmt"
	"os/user"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"github.com/elastic/go-windows"
	"github.com/winlabs/gowin32"
	xsyscall "golang.org/x/sys/windows"
)

const (
	localSystemSID       string = "S-1-5-18" //https://docs.microsoft.com/en-us/windows/security/identity-protection/access-control/security-identifiers
	localSystemSessionID uint32 = 0
)

var (
	modkernel32                = xsyscall.NewLazySystemDLL("kernel32.dll")
	queryFullProcessImageNameW = modkernel32.NewProc("QueryFullProcessImageNameW") //wchar compatible, requires Win XP/2003 or higher https://docs.microsoft.com/en-us/windows/win32/api/psapi/nf-psapi-getprocessimagefilenamew
)

type (
	MemoryApi interface {
		GetAllProcessesAndComputeDiff(oldProcs map[UniqueProcess]Process) (procs map[UniqueProcess]Process, changes JSONChanges, err error)
	}
	memoryApi struct {
		wts *gowin32.WTSServer
	}
)

func New() MemoryApi {
	return &memoryApi{wts: gowin32.OpenWTSServer("")}
}

//GetAllProcessesAndComputeDiff takes a map of all processes(can be empty) and returns the new map and changes that happened compared to the oldProcs
func (api *memoryApi) GetAllProcessesAndComputeDiff(oldProcs map[UniqueProcess]Process) (procs map[UniqueProcess]Process, changes JSONChanges, err error) {

	//get the list of all process IDs
	pids, err := windows.EnumProcesses()
	if err != nil {
		return nil, JSONChanges{}, fmt.Errorf("EnumProcesses(): %w", err)
	}

	procs = make(map[UniqueProcess]Process, len(pids))

	//iterate over them to get the handle
	for _, pid := range pids {
		handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, pid)
		if err != nil {
			continue
		}
		//close the handle after we're done (this will cause a piling up situation, but will close the handles no matter what. Safety measure.)
		defer func() {
			if err := syscall.CloseHandle(handle); err != nil {
				panic(err) //something went wrong horribly
			}
		}()

		//retrieve full executable path on the system. (queryFullProcessImageName won't work here since not all processes use win32 paths (stuff like WSL), it's safer to use native NT paths.  e.x.:  \Device\HarddiskVolume3\Windows\cmd.exe)
		//processPath, err := queryFullProcessImageName(handle)
		processPath, err := windows.GetProcessImageFileName(handle)
		if err != nil {
			return nil, JSONChanges{}, fmt.Errorf("GetProcessImageFileName(): %w", err)
		}
		//processName might have some regexpr checks on it
		processName := filepath.Base(processPath)

		//Get session ID from pid
		var sessionID uint32
		if err := xsyscall.ProcessIdToSessionId(pid, &sessionID); err != nil {
			return nil, JSONChanges{}, fmt.Errorf("ProcessIdToSessionId(): %w", err)
		}

		processTimes, err := getProcessTimes(handle)
		if err != nil {
			return nil, JSONChanges{}, fmt.Errorf("getProcessTimes(): %w", err)
		}

		creationTime := time.Unix(0, processTimes.CreationTime.Nanoseconds()).UTC().Format(time.RFC3339)

		//query the session info
		sessionInfo, err := api.wts.QuerySessionSesionInfo(uint(sessionID))
		if err != nil {
			return nil, JSONChanges{}, fmt.Errorf("QuerySessionSesionInfo(): %w", err)
		}

		//get user SID
		userSID := ""
		userName := ""
		switch sessionID {
		case localSystemSessionID:
			userSID = localSystemSID
			userName = "LocalSystem"
		default:
			user, err := user.Lookup(sessionInfo.UserName)
			if err != nil {
				return nil, JSONChanges{}, fmt.Errorf("Lookup(): %w", err)
			}
			userSID = user.Uid
			userName = sessionInfo.UserName
		}

		uniqueProc := UniqueProcess{PID: pid, CreationTime: creationTime}

		procs[uniqueProc] =
			Process{
				Name:           processName,
				PID:            pid,
				MainModulePath: processPath,
				CreationTime:   creationTime,
				User: User{
					SessionID: sessionID,
					Name:      userName,
					SID:       userSID,
					LastLogin: sessionInfo.LogonTime.UTC().Format(time.RFC3339),
				},
			}

		//check if current process exists in the prev iteration
		if _, exists := oldProcs[uniqueProc]; !exists {
			changes.New = append(changes.New, procs[uniqueProc])
		}
	}

	//check processes that were closed
	for k := range oldProcs {
		if _, exists := procs[k]; !exists {
			changes.Clsoed = append(changes.Clsoed, oldProcs[k])
		}
	}

	return procs, changes, nil
}

//queryFullProcessImageName returns an abosulte path from the given process handle. https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-queryfullprocessimagenamea
func queryFullProcessImageName(handle syscall.Handle) (string, error) {
	var buf [syscall.MAX_PATH]uint16
	n := uint32(len(buf))
	r1, _, e1 := queryFullProcessImageNameW.Call(
		uintptr(handle),
		uintptr(0), //0 - win32 type path, 1 - native path
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&n)))
	if r1 == 0 {
		if e1 != nil {
			return "", e1
		}
		return "", syscall.EINVAL
	}
	return syscall.UTF16ToString(buf[:n]), nil

}

//GetProcessTimes returns winApi times. https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getprocesstimes
func getProcessTimes(handle syscall.Handle) (changes ProcessTime, err error) {
	if err := syscall.GetProcessTimes(handle, &changes.CreationTime, &changes.ExitTime, &changes.KernelTime, &changes.UserTime); err != nil {
		return ProcessTime{}, err
	}
	return changes, nil
}
