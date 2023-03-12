package privilege

import (
	"golang.org/x/sys/windows"
)

func EnablePrivilege(name string) (err error) {
	p := windows.CurrentProcess()
	var token windows.Token
	err = windows.OpenProcessToken(p, windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY, &token)
	if err != nil {
		return
	}
	defer token.Close()

	newState := windows.Tokenprivileges{
		PrivilegeCount: 1,
	}

	systemName, err := windows.UTF16PtrFromString("")
	if err != nil {
		return
	}
	privilegeName, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return
	}
	err = windows.LookupPrivilegeValue(systemName, privilegeName, &newState.Privileges[0].Luid)
	if err != nil {
		return
	}
	newState.Privileges[0].Attributes = windows.SE_PRIVILEGE_ENABLED

	err = windows.AdjustTokenPrivileges(token, false, &newState, 0, nil, nil)
	if err != nil {
		return
	}

	return nil
}
