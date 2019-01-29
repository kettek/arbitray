// +build windows

package main

func (c *ArbitrayConfig) generateDefault() (err error) {
  c.Programs = make([]*ArbitrayProgram, 1)
  c.Programs[0] = &ArbitrayProgram{
    Title: "MSConfig",
    Tooltip: "Open MSConfig. Needs Admin",
    Program: "C:\\Windows\\System32\\msconfig.exe",
  }
  return
}
