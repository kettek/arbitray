const SysTray           = require('systray').default
const fs                = require('fs')
const path              = require('path')
const child_process     = require('child_process')
const bunyan            = require('bunyan')
const bunyanDebugStream = require('bunyan-debug-stream')
const dialog            = require('dialog')
const opn               = require('opn')

const uuidv4 = function(a){return a?(a^Math.random()*16>>a/4).toString(16):([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g,uuidv4)}

class Arbitray {
  constructor(uniqueName = uuidv4()) {
    this._uniqueName = uniqueName
    this._processes  = []
    // Create our logs directory.
    try {
      fs.mkdirSync(path.join(process.cwd(), 'logs'))
    } catch(err) {
      if (err.code !== 'EEXIST') {
        throw err
      }
    }
    // Create our logger.
    this.log = bunyan.createLogger({
      name: "Arbitray",
      streams: [
        {
          type: 'raw',
          stream: bunyanDebugStream({
            forceColor: true
          })
        },
        // Popup dialog for warning logs and up.
        {
          level: 'warn',
          stream: {
            write: data => {
              data = JSON.parse(data)
              if (data.level == bunyan.WARN) {
                dialog.warn(data.msg, data.name)
              } else {
                dialog.err(data.msg, data.name)
              }
            }
          }
        },
        {
          type: 'rotating-file',
          path: path.join(process.cwd(), 'logs', 'arbitray.log'),
          period: '1d',
          count: 3
        }
      ],
      serializers: bunyanDebugStream.serializers
    })
    this._logs       = []
    this.loadConfig()
    this.createTray()
    this.log.info("Arbitray started.")
  }

  loadConfig() {
    try {
      this.Config = require(path.join(process.cwd(), 'arbitray.json'))
    } catch(err) {
      this.log.warn(`Failed to open ${path.join(process.cwd(), 'arbitray.json')}, creating default configuration.`)
      this.Config = this.createDefaultConfig()
    }
  }

  createDefaultConfig() {
    let Config = {
      programs: [
      ]
    }
    if (process.platform === 'win32') {
      Config.programs.push({
        "title": "System Configuration",
        "program": "C:\\Windows\\System32\\msconfig.exe",
        "tooltip": "Open the System Configuration(msconfig)",
      })
      Config.programs.push({
        "title": "CMD",
        "program": "C:\\Windows\\System32\\cmd.exe",
        "tooltip": "Open a command prompt.",
        "arguments": [],
        "options": {}
      })
    } else {
      Config.programs.push({
        "title": "top",
        "program": "/usr/bin/top",
        "tooltip": "Show processes.",
        "arguments": [],
        "options": {}
      })
    }
    try {
      fs.writeFileSync(
        path.join(process.cwd(), 'arbitray.json'),
        JSON.stringify(Config, null, '\t')
      )
    } catch(err) {
      this.log.fatal(err)
    }
    return Config
  }

  /**
   * Getters and setters for arbitray's configuration.
   */
  get Config() {
    return this._config
  }
  set Config(config) {
    // Ensure sane defaults.
    this._config = {
      ...{
        programs: []
      },
      ...config
    }
    // Iterate through programs to ensure needed fields are populated.
    this._config.programs = this._config.programs.map((entry, index) => {
      if (!entry.program) {
        return
      }
      if (!entry.title) {
        entry.title = path.basename(entry.program, path.extname(entry.program))
      }
      if (!entry.tooltip) {
        entry.tooltip = "Run " + entry.title;
      }
      if (entry.options) {
        if (entry.options.cwd) {
          entry.options.cwd = path.resolve(entry.options.cwd)
        }
      }
      if (!entry.arguments) {
        entry.arguments = []
      }
      return entry
    })
  }

  get Icon() {
    return this._icon
  }
  set Icon(icon) {
    this._icon = icon
  }

  /**
   *
   */
  isProcessRunning(index) {
    return this._processes[index] ? true : false
  }
  startProcess(index) {
    if (index >= this.Config.programs.length) return false
    if (this._processes[index]) return false
    if (this.Config.programs[index].useSpawn) {
      this._processes[index] = child_process.spawn(
        this.Config.programs[index].program,
        this.Config.programs[index].arguments,
        {
          ...{ detached: false, windowsHide: true, },
          ...(this.Config.programs[index].options ? this.Config.programs[index].options : {})
        }
      )
    } else {
      this._processes[index] = child_process.exec(
        '"' + this.Config.programs[index].program + '" ' + this.Config.programs[index].arguments.join(' '),
        {
          ...{ detached: false, windowsHide: true, },
          ...(this.Config.programs[index].options ? this.Config.programs[index].options : {})
        }
      )
    }

    this.log.info(`Spawning "${this.Config.programs[index].title}"`)
    // Create logger for program
    this._logs[index] = bunyan.createLogger({
      name: this.Config.programs[index].title,
      pid:  this._processes[index].pid,
      streams: [
        {
          type: 'raw',
          stream: bunyanDebugStream({
            forceColor: true
          })
        },
        {
          level: 'warn',
          stream: {
            write: data => {
              data = JSON.parse(data)
              if (data.level == bunyan.WARN) {
                dialog.warn(data.msg, data.name)
              } else {
                dialog.err(data.msg, data.name)
              }
            }
          }
        },
        {
          type: 'rotating-file',
          path: path.join(process.cwd(), 'logs', this.Config.programs[index].title+'.log'),
          period: '1d',
          count: 3
        }

      ],
      serializers: bunyanDebugStream.serializers
    })
    // Set up handlers.
    this._processes[index].on('error', err => {
      this.log.error(`Failed to spawn ${this.Config.programs[index].title}: ${err}`)
    })
    if (this._processes[index].stdout) this._processes[index].stdout.on('data', data => this._logs[index].info(data.toString()))
    if (this._processes[index].stderr) this._processes[index].stderr.on('data', data => this._logs[index].error(data.toString()))
    this._processes[index].on('close', code => {
      this.log.info(`Closed "${this.Config.programs[index].title}" with code ${code}`)
      this.updateTrayProcess(index, {
        enabled: true,
        checked: false
      })
      delete this._processes[index]
      delete this._logs[index]
    })
    // Update tray icon.
    this.updateTrayProcess(index, {
      enabled: true,
      checked: true
    })
  }
  stopProcess(index) {
    if (index >= this.Config.programs.length) return false
    if (!this._processes[index]) return false
    if (process.platform === 'win32') {
      child_process.spawn('taskkill', ["/pid", this._processes[index].pid, '/f', '/t'])
    } else {
      this._processes[index].kill()
    }
    // 'close' handler will clean up our state.
  }
  updateTrayProcess(index, opts) {
    this._systray.sendAction({
      type: 'update-item',
      item: {
        ...this._systray._conf.menu.items[index],
        ...opts
      },
      seq_id: index,
    })
  }

  /**
   * Kill and recreate the system tray.
   */
  updateTray() {
    this._systray.kill(false)
    this.createTray()
  }
  /**
   * Create the system tray.
   */
  createTray() {
    let iconSrc = this.Config.icon ? this.Config.icon : '../resources/icon'
    this.Icon = require(iconSrc + (process.platform === 'win32' ? '.ico' : '.png'))

    this._systray = new SysTray({
      menu: {
        icon: this.Icon,
        title: "Arbitray",
        tooltip: "A list of processes to manage from `config.json`",
        items: [
          ...this.Config.programs.map((entry, index) => {
            return {
              title: entry.title,
              tooltip: entry.tooltip,
              checked: this.isProcessRunning(index),
              enabled: true,
            }
          }),
          ...[
            {
              title: "✎ Config",
              tooltip: "Open the configuration.",
              enabled: true,
            },
            {
              title: "📜 Logs",
              tooltip: "Open the logs directory.",
              enabled: true,
            },
            {
              title: "💀 Quit",
              tooltip: "Quit arbitray, closing all processes.",
              enabled: true,
            }
          ]
        ],
      },
      debug: true,
      copyDir: true,
    })
    this._systray.onClick(action => {
      // Program Action
      if (action.seq_id < this.Config.programs.length) {
        if (!this.isProcessRunning(action.seq_id)) {
          this.startProcess(action.seq_id)
        } else {
          this.stopProcess(action.seq_id)
        }
      }
      // Built-in Action
      if (action.seq_id >= this.Config.programs.length) {
        let id = action.seq_id - this.Config.programs.length
        console.log(id)
        if (id == 0) {
          opn( path.join(process.cwd(), 'arbitray.json') )
        } else if (id == 1) {
          opn( path.join(process.cwd(), 'logs') )
        } else if (id == 2) {
          this.Quit()
        }
      }
    })
  }

  /**
   * Cleanup then Desroy.
   */
  Quit() {
    this.Cleanup(() => {
      this.Destroy()
    })
  }
  /**
   * Destroy.
   */
  Destroy() {
    this._systray.kill(true)
    this.log.info("Arbitray closed.")
  }
  /**
   * Request each process to close, optionally calling cb once done.
   */
  Cleanup(cb) {
    let processes = this._processes.filter(item => item)
    if (processes.length == 0) {
      cb()
      return
    }
    processes.forEach((process, index) => {
      if (cb) {
        processes[index].on('close', () => {
          if (this._processes.filter(item => item).length == 0) {
            cb()
          }
        })
      }
      this.stopProcess(index)
    })
  }
}

module.exports = Arbitray
