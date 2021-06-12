// +build windows

package sysinfo

import "golang.org/x/sys/windows/registry"

func ReadRegistryString(k registry.Key, path string, name string) (string, error) {
	key, err := registry.OpenKey(k, path, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	s, _, err := key.GetStringValue(name)
	if err != nil {
		return "", err
	}
	return s, nil
}
