const Module    = require('module')
const path      = require('path')
const fs        = require('fs')
const Arbitray  = require('./arbitray')
const wasw      = require('windows-api-show-window')

// Hide console window on Windows ASAP.
if (process.platform === 'win32' && process.env.NODE_ENV !== 'production') {
  wasw.hideCurrentProcessWindow()
}

// Add handlers for ico and png to base64.
Module._extensions['.ico'] = Module._extensions['.png'] = (module, filename) => {
  return module._compile(
    'module.exports="'
    + fs.readFileSync(filename).toString('base64')
    + '"',
    filename
  )
}

// Adjust our CWD if an argument is provided.
if (process.argv.length == 3) {
  process.chdir(path.resolve(process.argv[2]))
}

// Start up Arbitray.
const arbitray = new Arbitray()
