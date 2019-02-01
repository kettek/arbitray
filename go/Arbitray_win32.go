// +build windows

package main

import (
  "os"
  "os/exec"
  "path/filepath"
  "syscall"
)

func (c *ArbitrayConfig) generateDefault() (err error) {
  c.Programs = make([]*ArbitrayProgram, 1)
  c.Programs[0] = &ArbitrayProgram{
    Title: "MSConfig",
    Tooltip: "Open MSConfig. Needs Admin",
    Program: "C:\\Windows\\System32\\msconfig.exe",
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
  if dir := filepath.Dir(p.Program); dir != "." {
    p.Cmd.Dir = dir
  }
  p.Cmd.Args = p.Arguments

  if p.Options.Hide {
    p.Cmd.SysProcAttr = &syscall.SysProcAttr{ HideWindow: true }
  }
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
