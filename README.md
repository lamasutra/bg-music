# bg-music
Background music / sound player based on states and events

## Wails 2

### Preparing Developer Environment

Ubuntu example:

Install go

`sudo add-apt-repository ppa:longsleep/golang-backports`
`sudo apt update`
`suda apt install golang`

Add go bin dir to PATH

`export PATH=$PATH:/usr/local/go/bin`

Check go version

`go version`

Install wails cli

`go install github.com/wailsapp/wails/v2/cmd/wails@latest`

Add bin to PATH if different from /usr/local/go/bin

`export PATH="$PATH:$(go env GOPATH)/bin"`

Check all is set up correctly

`wails doctor`

Install missing dependencies
`sudo apt install nsis`
`sudo apt install libgtk-3-dev`
`sudo apt install libwebkit2gtk-4.1-dev`

Systray requires libayatana

`sudo apt-get install gcc libgtk-3-dev libayatana-appindicator3-dev`

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.