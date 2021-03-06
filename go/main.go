/*
Copyright 2018-2019 Ketchetwahmeegwun T. Southall

This file is part of arbitray.

arbitray is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

arbitray is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with arbitray.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
  "github.com/getlantern/systray"
)

var arbitray Arbitray

func main() {
  systray.Run(arbitray.onReady, arbitray.onQuit)
}
