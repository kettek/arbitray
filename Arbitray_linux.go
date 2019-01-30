// +build linux

package main

import (
  "os/exec"
)

func open(path string) error {
  return exec.Command("xdg-open", [path]).Start()
}
