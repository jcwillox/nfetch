// +build !windows

package sysinfo

func shell(name string, exe string) (string, string, error) {
	return name, "", nil
}
