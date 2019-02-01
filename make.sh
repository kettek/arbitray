#/bin/sh


## Acquire tools
if ! [ -x "$(command -v go-apper)" ]; then
  echo "Installing go-apper..."
  go get github.com/kettek/go-apper
  if [ $? -ne 0 ]; then
    echo Failure executing go get github.com/kettek/go-apper
    exit
  fi
fi

if ! [ -x "$(command -v 2goarray)" ]; then
  echo "Installing 2goarray..."
  go get github.com/cratonica/2goarray
  if [ $? -ne 0 ]; then
    echo Failure executing go get github.com/cratonica/2goarray
    exit
  fi
fi

## Build icon
OUTPUT=go/ArbitrayIcon.go
echo Generating $OUTPUT
echo "//+build linux darwin" > $OUTPUT
echo >> $OUTPUT
cat "resources/arbitray.png" | 2goarray iconData main >> $OUTPUT
if [ $? -ne 0 ]; then
  echo Failure generating $OUTPUT
  exit
fi

## Build arbitray
cd go && go build -o ../arbitray && cd ..
if [ $? -ne 0 ]; then
  echo Failure building Arbitray
  exit
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
  ## Package into .app
  go-apper -bin arbitray -icon resources/arbitray.png -identifier net.kettek.arbitray -name "Arbitray" -o ./
fi

echo Finished
