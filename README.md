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

Although more will inevitably be added in the future, at the moment the only important section is `programs`.

### programs
The programs section is an array of objects that define a runnable entry in the system tray. Each entry can have the following properties:

  * `title` -- The name to use in the tray.
  * `tooltip` -- The tooltip to use in the tray.
  * `program` -- The location to the program to run.
  * `arguments` -- An array of arguments to pass to the program.
  * `options` -- I forget.
