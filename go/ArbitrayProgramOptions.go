package main

import (
)

type ArbitrayProgramOptions struct {
  CWD string `json:"cwd,omitempty"`
  Hide bool `json:"hide,omitempty"`
  CloseCmd string `json:"closeCmd,omitempty"`
}
