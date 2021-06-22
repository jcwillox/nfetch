package sysinfo

import (
	"bufio"
	"nfetch/pkg/utils"
	"path/filepath"
)

func Pkgs(show []string) map[string]int {
	if len(show) == 0 {
		show = AllPkgManagers
	}

	// is the use of a channel & struct more efficient than a concurrent map & wait group?
	type Result struct {
		Manager string
		Pkgs    int
	}

	results := make(chan Result)
	pkgMap := make(map[string]int)

	for _, manager := range show {
		go func(manager string) {
			pkgs, _ := CountPkgs(manager)
			results <- Result{
				Manager: manager,
				Pkgs:    pkgs,
			}
		}(manager)
	}

	for i := 0; i < len(show); i++ {
		result := <-results

		if result.Pkgs > 0 {
			pkgMap[result.Manager] = result.Pkgs
		}
	}

	return pkgMap
}

func countCmdLines(offset int, include func(line string) bool, cmdArgs ...string) (int, error) {
	reader, err := utils.StartCommand(cmdArgs[0], cmdArgs[1:]...)
	if err != nil {
		return 0, err
	}
	return countReaderLines(reader, offset, include), nil
}

func countDirEntries(pattern string, offset int) (int, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0, err
	}
	return len(matches) + offset, nil
}

func countReaderLines(reader *bufio.Reader, offset int, include func(line string) bool) int {
	count := offset
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if include == nil || include(line) {
			count++
		}
	}
	return count
}
