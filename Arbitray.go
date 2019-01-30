package main

import (
  "fmt"
  "github.com/kettek/arbitray-go/icon"
  "github.com/getlantern/systray"
  "github.com/gen2brain/dlgs"
  "sync"
  "io"
  "os"
  "os/exec"
  "bufio"
  "log"
  "path"
  "path/filepath"
)

type Arbitray struct {
  config ArbitrayConfig
  waitGroup sync.WaitGroup
}

func (a *Arbitray) Init() (err error) {
  a.config.Load()
  os.Mkdir("logs", os.ModeDir)
  systray.Run(a.onReady, a.onQuit)
  return
}

func (a *Arbitray) onReady() {
  systray.SetIcon(icon.Data)
  systray.SetTitle("Arbitray")
  systray.SetTooltip("Arbitrary Process Launcher")

  // Add our processes as menu items.
  for _, program := range a.config.Programs {
    program.MenuItem = systray.AddMenuItem(program.Title, program.Tooltip)
    program.CloseChan = make(chan bool)
    program.KillChan = make(chan bool)
    program.Log = log.New(nil, "", log.LstdFlags)
    // This seems heavy to have go routines handling each entry's input...
    go func(program *ArbitrayProgram) {
      for {
        select {
        case <-program.MenuItem.ClickedCh:
          if !program.MenuItem.Checked() {
            a.waitGroup.Add(1)
            go a.startProgram(program)
          } else {
            program.KillChan <- true
          }
        }
      }
    }(program)
  }
  systray.AddSeparator()
  // Add our base items.
  mConfig := systray.AddMenuItem("✎ Config", "Config Arbitray")
  go func() {
    <-mConfig.ClickedCh
    // Get absolute location of arbitray. FIXME: This should be CWD.
    var dir string
    var err error
    if dir, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
      log.Fatal(err)
    }
    open(path.Join(dir, "arbitray.json"))
  }()

  mLogs := systray.AddMenuItem("📜 Logs", "Logs Arbitray")
  go func() {
    <-mLogs.ClickedCh
    var dir string
    var err error
    // Get absolute location of arbitray. FIXME: This should be CWD.
    if dir, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
      log.Fatal(err)
    }
    open(path.Join(dir, "logs"))
  }()

  mQuit := systray.AddMenuItem("💀 Quit", "Quit Arbitray")
  go func() {
    <-mQuit.ClickedCh
    a.Quit()
  }()

}
func (a *Arbitray) onQuit() {
  fmt.Println("Should do cleanup here.")
}
func (a *Arbitray) Quit() {
  fmt.Println("Quit, apparently.")
  for index, _ := range a.config.Programs {
    if a.config.Programs[index].MenuItem.Checked() == true {
      a.config.Programs[index].KillChan <- true
    }
  }
  a.waitGroup.Wait()
  systray.Quit()
}

func (a *Arbitray) startProgram(p *ArbitrayProgram) {
  defer func() {
    p.MenuItem.Uncheck()
    a.waitGroup.Done()
  }()
  p.MenuItem.Check()

  // Create our loggers.
  logFile, err := os.OpenFile(fmt.Sprintf("%s/%s.log", "logs", p.Title), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
  if err != nil {
    log.Fatal(err)
  }
  defer logFile.Close()
  p.Log.SetOutput(logFile)

  // Set up our command.
  p.Cmd = exec.Command(p.Program)

  // stdout
  stdoutChan := make(chan string)
  go func() {
    stdout, err := p.Cmd.StdoutPipe()
    if err != nil {
      fmt.Printf("Uhoh, error getting stdout: %v\n", err)
    }
    reader := bufio.NewReader(stdout)
    for {
      in, err := reader.ReadString('\n')
      if err != nil {
        if err != io.EOF {
          fmt.Printf("Uhoh, stdout error: %v\n", err)
        }
        return
      }
      stdoutChan <- in
    }
  }()
  // stderr
  stderrChan := make(chan string)
  go func() {
    stderr, err := p.Cmd.StderrPipe()
    if err != nil {
      fmt.Printf("Uhoh, error getting stderr: %v\n", err)
    }

    reader := bufio.NewReader(stderr)
    for {
      in, err := reader.ReadString('\n')
      if err != nil {
        if err != io.EOF {
          fmt.Printf("Uhoh, stderr error: %v\n", err)
        }
        return
      }
      stderrChan <- in
    }
  }()
  // Run our command.
  fmt.Printf("[%s]: %s starting.\n", "Arbitray", p.Title)
  if err := p.Cmd.Start(); err != nil {
    dlgs.Error("Arbitray", err.Error())
    fmt.Printf("[%s] %s\n", p.Title, err.Error())
    p.Log.Printf("Error: %s\n", err.Error())
  }
  go func() {
    if err := p.Cmd.Wait(); err != nil {
      fmt.Printf("[%s]: %s\n", p.Title, err.Error())
      p.Log.Printf("Error: %s\n", err.Error())
      //dlgs.Error("Arbitray", err.Error())
    }
    p.CloseChan <- true
  }()
  //
  ListenLoop:
    for {
      select {
      case msg := <-stdoutChan:
        fmt.Printf("[%s] %s", p.Title, msg)
        p.Log.Println(msg)
      case msg := <-stderrChan:
        fmt.Printf("[%s] %s", p.Title, msg)
        p.Log.Printf("Error: %s", msg)
      case <-p.KillChan:
        p.Cmd.Process.Signal(os.Kill)
      case <-p.CloseChan:
        break ListenLoop
      }
    }
  fmt.Printf("[%s]: %s finished.\n", "Arbitray", p.Title)
}
