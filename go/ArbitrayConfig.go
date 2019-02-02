package main

import (
  "os"
  "path/filepath"
  "strings"
  "log"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "github.com/gen2brain/dlgs"
)

type ArbitrayConfig struct {
  Programs []*ArbitrayProgram `json:"programs,omitempty"`
  HideItems map[string]bool `json:"hideItems,omitempty"`
}

func (c *ArbitrayConfig) Load() (err error) {
  var file *os.File
  var bytes []byte

  // Open arbitray.json or create it if it does not exist.
  if _, err := os.Stat("arbitray.json"); err != nil {
    if os.IsNotExist(err) {
      dlgs.Warning("Arbitray", "No arbitray.json found, creating it with defaults.\nPlease configure it via the menu bar icon.")
      if err = c.generateDefault(); err != nil {
        dlgs.Error("Arbitray", fmt.Sprintf("Issue generating default arbitray.json: %s", err.Error()))
        log.Fatal(err)
      }
      if err = c.Save(); err != nil {
        dlgs.Error("Arbitray", fmt.Sprintf("Issue saving arbitray.json: %s", err.Error()))
        log.Fatal(err)
      }
    } else {
      dlgs.Error("Arbitray", fmt.Sprintf("Issue loading arbitray.json: %s", err.Error()))
      log.Fatal(err)
    }
  }

  if file, err = os.OpenFile("arbitray.json", os.O_RDWR|os.O_CREATE, 0644); err != nil {
    dlgs.Error("Arbitray", err.Error())
    log.Fatal(err)
  }
  defer file.Close()

  // Read JSON into Config.
  if bytes, err = ioutil.ReadAll(file); err != nil {
    dlgs.Error("Arbitray", err.Error())
    log.Fatal(err)
  }
  json.Unmarshal([]byte(bytes), &c)
  // Ensure some sanity.
  c.Ensure()
  return
}

func (c *ArbitrayConfig) Save() (err error) {
  var out []byte

  // Create our indented JSON string.
  if out, err = json.MarshalIndent(c, "", "\t"); err != nil {
    return err
  }

  // Write it out.
  if err = ioutil.WriteFile("arbitray.json", out, 0644); err != nil {
    return err
  }
  return
}

func (c *ArbitrayConfig) Ensure() (err error) {
  for _, program := range c.Programs {
    if program.Program == "" {
      continue
    }
    if program.Title == "" {
      _, exe := filepath.Split(program.Program)

      if n := strings.LastIndexByte(exe, '.'); n > 0 {
        program.Title = exe[:n]
      } else {
        program.Title = exe
      }
    }
    if program.Tooltip == "" {
      program.Tooltip = fmt.Sprintf("Run %s", program.Title)
    }
  }
  return
}
