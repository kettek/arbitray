package main

import (
  "github.com/getlantern/systray"
  "os/exec"
)

type ArbitrayProgram struct {
  Title string `json:"title,omitempty"`
  Tooltip string `json:"tooltip,omitempty"`
  Program string `json:"program,omitempty"`
  Arguments []string `json:"arguments,omitempty"`
  Options ArbitrayProgramOptions `json:"options,omitempty"`
  MenuItem *systray.MenuItem
  Cmd *exec.Cmd
  CloseChan chan bool
  KillChan chan bool
}
