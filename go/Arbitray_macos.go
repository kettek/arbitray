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

func getAppDir() (loc string, isApp bool, err error) {
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
    loc = filepath.Join(dir, "../../")
    isApp = true
  } else {
    loc = exe
    isApp = false
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
func openDir(path string) error {
  return exec.Command("open", []string{path+"/"}...).Start()
}

func restart() {
    exe, isApp, _ := getAppDir()
    if isApp == true {
      cmd := exec.Command("open", append([]string{"-n", exe}, os.Args[1:]...)...)
      cmd.Output()
    } else {
      cmd := exec.Command(exe, os.Args[1:]...)
      cmd.Start()
      cmd.Process.Release()
    }
}
