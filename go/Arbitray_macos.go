// +build darwin

package main

import (
  "os"
  "os/exec"
  "path/filepath"
)

func (c *ArbitrayConfig) generateDefault() (err error) {
  c.Programs = make([]*ArbitrayProgram, 1)
  c.Programs[0] = &ArbitrayProgram{
    Title: "Screenshot",
    Program: "/Applications/Utilities/Screenshot.app/Contents/MacOS/Screenshot",
  }
  return
}

func (a *Arbitray) platformInit() (err error) {
  var exe string
  var dir string
  var base string

  if exe, err = os.Executable(); err != nil {
    return
  }
  // Is an app
  dir = filepath.Dir(exe)
  base = filepath.Base(dir)
  if base == "MacOS" {
    a.workingDir = filepath.Join(dir, "../../../")
  } else {
    a.workingDir = dir
  }

  return
}

func open(path string) error {
  return exec.Command("open", []string{path}...).Start()
}

func restart() {
  cmd := exec.Command(os.Args[:1][0], os.Args[1:]...)
  cmd.Output()
}
