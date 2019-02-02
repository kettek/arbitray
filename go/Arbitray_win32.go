// +build windows

package main

import (
  "os"
  "os/exec"
  "path/filepath"
  "syscall"
)

const (
  CONFIG_STRING = "Edit"
  RELOAD_STRING = "Reload"
  LOGS_STRING = "Logs"
  QUIT_STRING = "Exit"
)

func (c *ArbitrayConfig) generateDefault() (err error) {
  c.Programs = make([]*ArbitrayProgram, 1)
  c.Programs[0] = &ArbitrayProgram{
    Program: "explorer.exe",
    Arguments: []string{"ms-screenclip:"},
    Title: "Take Screenshot",
  }
  return
}

func (a *Arbitray) platformInit() (err error) {
  var exe string

  if exe, err = os.Executable(); err != nil {
    return
  }
  a.workingDir = filepath.Dir(exe)
  return
}

func (p *ArbitrayProgram) CreateCommand() (err error) {
  p.Cmd = exec.Command(p.Program)
  if p.Options.CWD != "" {
    p.Cmd.Dir = p.Options.CWD
  } else {
    if dir := filepath.Dir(p.Program); dir != "." {
      p.Cmd.Dir = dir
    }
  }
  p.Cmd.Args = append([]string{p.Program}, p.Arguments...)

  if p.Options.Hide {
    p.Cmd.SysProcAttr = &syscall.SysProcAttr{ HideWindow: true }
  }
  return
}

func (p *ArbitrayProgram) Kill() (err error) {
  err = p.Cmd.Process.Signal(os.Kill)
  return
}

func open(path string) error {
  cmd := exec.Command("cmd", []string{"/c", "start", path}...)
  cmd.SysProcAttr = &syscall.SysProcAttr{ HideWindow: true }
  _, err := cmd.Output()
  return err
}

func openDir(path string) error {
  cmd := exec.Command("cmd", []string{"/c", "start", path+"\\"}...)
  cmd.SysProcAttr = &syscall.SysProcAttr{ HideWindow: true }
  _, err := cmd.Output()
  return err
}

func restart() {
  cmd := exec.Command("cmd", append([]string{"/c"}, os.Args...)...)
  cmd.SysProcAttr = &syscall.SysProcAttr{ HideWindow: true }
  cmd.Output()
}
