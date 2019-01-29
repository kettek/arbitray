package main

import (
  "fmt"
  "github.com/getlantern/systray"
  "github.com/getlantern/systray/example/icon"
  "github.com/gen2brain/dlgs"
  "sync"
  "io"
  "os"
  "os/exec"
  "bufio"
)

type Arbitray struct {
  config ArbitrayConfig
  waitGroup sync.WaitGroup
}

func (a *Arbitray) Init() (err error) {
  a.config.Load()
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
    // This seems heavy to have go routines for each entry...
    go func() {
      for {
        select {
        case <-program.MenuItem.ClickedCh:
          if !program.MenuItem.Checked() {
            a.waitGroup.Add(1)
            go a.startProgram(program)
          } else {
            program.KillChan <- true
          }
          fmt.Printf("%s clicked.\n", program.Title)
        }
      }
    }()
  }
  systray.AddSeparator()
  // Add our base items.
  mQuit := systray.AddMenuItem("Quit", "Quit Arbitray")
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
  if err := p.Cmd.Start(); err != nil {
    dlgs.Error("Arbitray", err.Error())
  }
  go func() {
    if err := p.Cmd.Wait(); err != nil {
      fmt.Printf("%s closed with: %s\n", p.Title, err.Error())
      //dlgs.Error("Arbitray", err.Error())
    }
    p.CloseChan <- true
  }()
  //
  ListenLoop:
    for {
      select {
      case msg := <-stdoutChan:
        fmt.Printf("[%s] %s\n", p.Title, msg)
      case msg := <-stderrChan:
        fmt.Printf("[%s] %s\n", p.Title, msg)
      case <-p.KillChan:
        p.Cmd.Process.Signal(os.Kill)
      case <-p.CloseChan:
        break ListenLoop
      }
    }
}
