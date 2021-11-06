package sysinfo

import (
	"github.com/jcwillox/nfetch/pkg/utils"
	"github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/process"
	"os"
	"strings"
)

func Terminal() (string, error) {
	// identify terminal from environment variables
	if _, present := os.LookupEnv("WT_SESSION"); present {
		return "Windows Terminal", nil
	}

	termProgram := os.Getenv("TERM_PROGRAM")
	switch termProgram {
	case "iTerm.app":
		return "iTerm2", nil
	case "Terminal.app":
		return "Apple Terminal", nil
	case "Hyper":
		return "HyperTerm", nil
	case "FluentTerminal":
		return "Fluent Terminal", nil
	case "vscode":
		return "Visual Studio Code", nil
	default:
		if termProgram != "" {
			return strings.TrimSuffix(termProgram, ".app"), nil
		}
	}

	term := os.Getenv("TERM")
	if term == "tw52" || term == "tw100" {
		return "TosWin2", nil
	}

	if term, present := os.LookupEnv("TERMINAL_EMULATOR"); present {
		return term, nil
	}

	if _, present := os.LookupEnv("SSH_CONNECTION"); present {
		return os.Getenv("SSH_TTY"), nil
	}

	// identify terminal from process tree
	processes, err := ps.Processes()
	if err != nil {
		// unable to list processes
		return getTTY(), err
	}

	proc := findProcess(processes, os.Getppid())
	if proc == nil {
		// unable to find parent process
		return getTTY(), err
	}

	// it should never be the direct parent, as that's the 'shell'
	// so we immediately fetch the next parent

	for {
		// fetch next parent
		parent := findProcess(processes, proc.PPid())
		if parent == nil {
			return getTTY(), err
		}

		proc = parent
		name := utils.GetFileName(proc.Executable())

		switch {
		case strings.HasSuffix(name, "sh"), name == "screen", strings.HasPrefix(name, "su"),
			name == "newgrp", name == os.Getenv("SHELL"):
			break
		case name == "init", name == "(init)", strings.HasPrefix(name, "login"), strings.Contains(name, "Login"):
			return getTTY(), nil
		case name == "explorer", name == "conhost":
			return "Windows Console", nil
		case name == "Console":
			return "Console2/Z", nil
		case name == "ConEmuC64":
			return "ConEmu", nil
		case name == "gnome-terminal-":
			return "gnome-terminal", err
		case name == "urxvtd":
			return "urxvt", err
		case strings.HasSuffix(name, "nvim"):
			return "Neovim Terminal", err
		case strings.Contains(name, "NeoVimServer"):
			return "VimR Terminal", err
		case name == "ruby", name == "1", name == "systemd", name == "bwrap", name == "go",
			utils.HasPrefixMulti(name, "tmux", "sshd", "python", "USER", "kdeinit", "launchd"):
			return "(unknown)", err
		default:
			if strings.HasSuffix("-wrapped", name) {
				name = strings.TrimPrefix(name, ".")
				return strings.TrimSuffix(name, "-wrapped"), nil
			}
			return name, nil
		}
	}
}

func getTTY() string {
	terminal, _ := (&process.Process{Pid: int32(os.Getpid())}).Terminal()
	return terminal
}

func findProcess(process []ps.Process, pid int) ps.Process {
	for _, p := range process {
		if p.Pid() == pid {
			return p
		}
	}
	return nil
}
