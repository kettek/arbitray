package main

import (
  "os"
  "path/filepath"
  "strings"
  "path"
  "log"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "github.com/gen2brain/dlgs"
)

type ArbitrayConfig struct {
  Programs []*ArbitrayProgram `json:"programs,omitempty"`
}

func (c *ArbitrayConfig) Load() (err error) {
  var dir string
  var file *os.File

  // Get absolute location of arbitray. FIXME: This should be CWD.
  if dir, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
    log.Fatal(err)
  }
  filepath := path.Join(dir, "arbitray.json")

  // Open arbitray.json or create it if it does not exist.
  if _, err := os.Stat(filepath); err == nil {
    log.Print("Loading arbitray.json")
    if file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644); err != nil {
      log.Fatal(err)
    }
    defer file.Close()
  } else if os.IsNotExist(err) {
    dlgs.Warning("Arbitray", "No arbitray.json found, creating it with defaults.")
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

  // Read JSON into Config.
  bytes, _ := ioutil.ReadAll(file)
  json.Unmarshal([]byte(bytes), &c)
  // Ensure some sanity.
  c.Ensure()
  return
}

func (c *ArbitrayConfig) Save() (err error) {
  var dir string
  var out []byte

  // Get absolute location of arbitray. FIXME: This should be CWD.
  if dir, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
    return err
  }
  filepath := path.Join(dir, "arbitray.json")

  // Create our indented JSON string.
  if out, err = json.MarshalIndent(c, "", "\t"); err != nil {
    return err
  }

  // Write it out.
  if err = ioutil.WriteFile(filepath, out, 0644); err != nil {
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
