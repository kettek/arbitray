// +build linux

package main

import (
  "os"
  "os/exec"
  "syscall"
  "path/filepath"
)

const (
  CONFIG_STRING = "Edit"
  RELOAD_STRING = "Reload"
  LOGS_STRING = "Logs"
  QUIT_STRING = "Quit"
)

func (c *ArbitrayConfig) generateDefault() (err error) {
  c.Programs = make([]*ArbitrayProgram, 1)
  c.Programs[0] = &ArbitrayProgram{
    Title: "Screenshot",
    Program: "scrot",
  }
  return
}

func getAppDir() (loc string, err error) {
  var exe string

  if exe, err = os.Executable(); err != nil {
    return
  }
  loc = exe
  return
}

func (a *Arbitray) platformInit() (err error) {
  var exe string
  var dir string

  if exe, err = os.Executable(); err != nil {
    return
  }

  dir = filepath.Dir(exe)
  a.workingDir = dir

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
  }
  return
}

func (p *ArbitrayProgram) Kill() (err error) {
  err = p.Cmd.Process.Signal(syscall.SIGINT)
  return
}

func open(path string) error {
  return exec.Command("open", []string{path}...).Start()
}
func openDir(path string) error {
  return exec.Command("open", []string{path+"/"}...).Start()
}

func restart() {
  exe, _ := getAppDir()
  cmd := exec.Command(exe, os.Args[1:]...)
  cmd.Start()
  cmd.Process.Release()
}
