# Arbitray -- The Arbitrary Process Manager

![Screenshot](screenshot-win10.PNG?raw=true)

**Arbitray** is a wonderfully named system tray application that you, the almighty user, can use to manage applications.

The truly fantastic things are that:

  * It's really simple. Just a system tray icon with checkmarks to show whether or not a process is running.
  * Configuration is done via a JSON file.
  * Arbitray starts in whatever location it is invoked from. Or by passing it a location as its first argument.
  * Processes started by Arbitray are closed along with it.
  
So, as is undeniably obvious, it is the perfect solution for managing random crap services on your computer, such a Minecraft server.

## Configuration File
Arbitray will automatically create an `arbitray.json` file in the working directory that the arbitray executable is started from. This can be used as a starting point for your own configuration.

### useSpawn
This boolean value specifies if [spawn](https://nodejs.org/api/child_process.html#child_process_child_process_spawn_command_args_options) should be used. If not specified, [exec](https://nodejs.org/api/child_process.html#child_process_child_process_exec_command_options_callback) will be used.

### programs
The programs section is an array of objects that define a runnable entry in the system tray. Each entry can have the following properties:

  * `title` -- The name to use in the tray.
  * `tooltip` -- The tooltip to use in the tray.
  * `program` -- The location to the program to run.
  * `arguments` -- An array of arguments to pass to the program.
  * `options` -- Any of the options listed [here](https://nodejs.org/api/child_process.html#child_process_child_process_spawn_command_args_options).

## FAQ

  * 1. Wouldn't this be better in Go?
    * _Yes._
  * 2. What platforms are supported?
    * Tested is Mac OS Mojave and Windows 10. Both 64-bit. Linux might work but I don't think so due to `pkg` not being able to pack native binaries. Additionally, when I ran arbitray from source, nothing showed up in stalonetray.
