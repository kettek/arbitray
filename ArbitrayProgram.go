package main

import (
  "github.com/getlantern/systray"
  "os/exec"
  "log"
)

type ArbitrayProgram struct {
  Title string `json:"title,omitempty"`
  Tooltip string `json:"tooltip,omitempty"`
  Program string `json:"program,omitempty"`
  Arguments []string `json:"arguments,omitempty"`
  Options ArbitrayProgramOptions `json:"options,omitempty"`
  MenuItem *systray.MenuItem `json:"-"`
  Cmd *exec.Cmd `json:"-"`
  CloseChan chan bool `json:"-"`
  KillChan chan bool `json:"-"`
  Log *log.Logger `json:"-"`
}
