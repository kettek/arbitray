@ECHO OFF

IF "%GOPATH%"=="" GOTO NOGO

IF NOT EXIST %GOPATH%\bin\2goarray.exe GOTO INSTALL2GOARRAY
:POST2GOARRAY
IF NOT EXIST %GOPATH%\bin\goversioninfo.exe GOTO INSTALLGOVERSIONINFO
:POSTGOVERSIONINFO

:BUILDICON
IF "resources\arbitray.ico"=="" GOTO NOICO
IF NOT EXIST resources\arbitray.ico GOTO BADFILE
ECHO Creating go\ArbitrayIcon.go
ECHO //+build windows > go\ArbitrayIcon.go
ECHO. >> go\ArbitrayIcon.go
TYPE resources\arbitray.ico | %GOPATH%\bin\2goarray iconData main >> go\ArbitrayIcon.go

:BUILDVERSIONINFO
ECHO Building version info and resources
cd go
%GOPATH%\bin\goversioninfo
cd ..

:BUILDARBITRAY
ECHO Building arbitray
cd go
go build -o ..\arbitray.exe -ldflags "-H=windowsgui"
cd ..
GOTO DONE

:CREATEFAIL
ECHO Unable to create output file
GOTO DONE

:INSTALL2GOARRAY
ECHO Installing 2goarray...
go get github.com/cratonica/2goarray
IF ERRORLEVEL 1 GOTO GET2GOARRAYFAIL
GOTO POST2GOARRAY

:GET2GOARRAYFAIL
ECHO Failure running go get github.com/cratonica/2goarray.  Ensure that go and git are in PATH
GOTO DONE

:INSTALLGOVERSIONINFO
ECHO Installing goversioninfo...
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
IF ERRORLEVEL 1 GOTO GETGOVERSIONINFOFAIL
GOTO POSTGOVERSIONINFO

:GETGOVERSIONINFOFAIL
ECHO Failure running go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo. Ensure that go and git are in PATH
GOTO DONE

:NOGO
ECHO GOPATH environment variable not set
GOTO DONE

:NOICO
ECHO Please specify a .ico file
GOTO DONE

:BADFILE
ECHO %1 is not a valid file
GOTO DONE

:DONE

