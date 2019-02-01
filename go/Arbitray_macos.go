// +build macos

package main

import (
  "os/exec"
)

func (c *ArbitrayConfig) generateDefault() (err error) {
  c.Programs = make([]*ArbitrayProgram, 1)
  c.Programs[0] = &ArbitrayProgram{
    Title: "XTerm",
    Program: "xterm",
  }
  return
}

func open(path string) error {
  return exec.Command("open", [path]).Start()
}
