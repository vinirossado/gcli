package helper

import (
	"os"
	"runtime"
)

func getUnixDefaultFilePermissions() os.FileMode {
	return 0644
}

func getWindowsDefaultFilePermissions() os.FileMode {
	return 0666
}

func GetDefaultOSPermissionFile() os.FileMode {
	if runtime.GOOS == "windows" {
		return getWindowsDefaultFilePermissions()
	}
	return getUnixDefaultFilePermissions()
}
