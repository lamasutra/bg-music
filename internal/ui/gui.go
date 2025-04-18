package ui

import (
	"context"
	"embed"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/lamasutra/bg-music/internal/ui/gui"
	"github.com/lamasutra/bg-music/pkg/events"

	bgplogger "github.com/lamasutra/bg-music/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type guiState struct {
	ctx            context.Context // app
	onStartup      func()
	assets         *embed.FS
	icon           []byte
	visible        bool
	playerControls *gui.PlayerControls
	mShow          *systray.MenuItem
}

func (s *guiState) startup(ctx context.Context) {
	s.ctx = ctx
	runtime.Hide(s.ctx)
	events.ListenAll("gui", s.eventDispatcher)
	events.Listen("log", "cli", s.renderMessage)
	s.onStartup()
}

func (s *guiState) renderMessage(args ...any) {
	if len(args) != 1 {
		panic("invalid arguments count, cannot render message")
	}

	msg, ok := args[0].(bgplogger.MessageRenderer)
	if !ok {
		panic("message is not renderer")
	}

	fmt.Println(msg.Render())
}

func (s *guiState) eventDispatcher(event string, values ...any) {
	runtime.EventsEmit(s.ctx, event, values...)
}

// domReady is called after front-end resources have been loaded
func (s guiState) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (s *guiState) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (s *guiState) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func NewGui(assets *embed.FS, icon []byte) *guiState {
	// Create an instance of the app structure
	app := &guiState{
		assets: assets,
		icon:   icon,
	}

	return app
}

func (s *guiState) Run(onStartup func()) {
	s.onStartup = onStartup

	// theme := runtime.Theme(s.ctx)
	// if theme == runtime.ThemeDark {
	// 	println("Dark mode is active")
	// } else {
	// 	println("Light mode is active")
	// }

	app := s.createApplication()

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *guiState) MinimizeWindow() {
	runtime.WindowMinimise(s.ctx)
}

func (s *guiState) CloseApp() {
	runtime.Quit(s.ctx)
}

func (s *guiState) createApplication() *application.Application {
	s.playerControls = gui.NewPlayerControls()

	// Create application with options
	app := application.NewWithOptions(
		&options.App{
			Title:             "Background Music Player",
			Width:             500,
			Height:            304,
			MinWidth:          500,
			MinHeight:         304,
			MaxWidth:          500,
			MaxHeight:         304,
			DisableResize:     false,
			Fullscreen:        false,
			Frameless:         true,
			StartHidden:       true,
			HideWindowOnClose: false,
			BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
			AssetServer: &assetserver.Options{
				Assets: s.assets,
			},
			Menu:             nil,
			Logger:           nil,
			LogLevel:         logger.DEBUG,
			OnStartup:        s.startup,
			OnDomReady:       s.domReady,
			OnBeforeClose:    s.beforeClose,
			OnShutdown:       s.shutdown,
			WindowStartState: options.Normal,
			Bind: []interface{}{
				s,
				s.playerControls,
			},
			Linux: &linux.Options{
				Icon:        s.icon,
				ProgramName: "BG Music Player",
			},
			// Windows platform specific options
			Windows: &windows.Options{
				WebviewIsTransparent: false,
				WindowIsTranslucent:  false,
				DisableWindowIcon:    false,
				// DisableFramelessWindowDecorations: false,
				WebviewUserDataPath: "",
				ZoomFactor:          1.0,
			},
			// Mac platform specific options
			Mac: &mac.Options{
				TitleBar: &mac.TitleBar{
					TitlebarAppearsTransparent: false,
					HideTitle:                  false,
					HideTitleBar:               false,
					FullSizeContent:            false,
					UseToolbar:                 false,
					HideToolbarSeparator:       true,
				},
				Appearance:           mac.NSAppearanceNameDarkAqua,
				WebviewIsTransparent: true,
				WindowIsTranslucent:  true,
				About: &mac.AboutInfo{
					Title:   "bg-player",
					Message: "",
					Icon:    s.icon,
				},
			},
		})

	s.createSystray()
	// wails3 systray
	// Tray: &tray.Options{
	// 	Icon:    "frontend/public/icon.png", // or use icon bytes
	// 	Menu:    trayMenu,
	// 	Tooltip: "Wails Tray App",
	// },
	// systray := app.NewSystemTray()
	// systray.SetLabel("My App")
	// systray.SetIcon(iconBytes)
	// systray.Run()

	// trayMenu := menu.NewMenu()
	// trayMenu.Append(menu.Text("Show App", nil, func(_ *menu.CallbackData) {
	// 	// You can add logic to show the main window here
	// }))
	// trayMenu.Append(menu.Text("Quit", nil, func(_ *menu.CallbackData) {
	// 	wails.Quit()
	// }))

	return app
}

func (s *guiState) Debug(args ...any) {
	length := len(args) + 1
	buf := make([]string, length)
	buf[0] = time.Now().Format("15:04:05.000")
	for i, val := range args {
		buf[i+1] = fmt.Sprint(val)
	}
	fmt.Println(strings.Join(buf, " "))
}

func (s *guiState) Write(p []byte) (n int, err error) {
	return fmt.Println(string(p))
}

func (s *guiState) Error(args ...any) {
	newArgs := []any{"ERR:"}
	newArgs = append(newArgs, args...)
	s.Debug(newArgs...)
}

func (s *guiState) createSystray() {
	onReady := func() {
		systray.SetTemplateIcon(s.icon, s.icon)
		systray.SetTitle("Bg Music")
		systray.SetTooltip("play")

		go s.handleSystray()
	}
	onExit := func() {
		// now := time.Now()
		// ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	}

	systray.Register(onReady, onExit)
}

func (s *guiState) ToggleVisibility() {
	if !s.visible {
		runtime.Show(s.ctx)
		s.mShow.SetTitle("Hide")
		s.visible = true
	} else {
		runtime.Hide(s.ctx)
		s.mShow.SetTitle("Show")
		s.visible = false
	}
}

func (s *guiState) handleSystray() {
	s.mShow = systray.AddMenuItem("Show", "Show the app")
	systray.AddSeparator()
	mOptions := systray.AddMenuItem("Options", "Show options")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	for {
		select {
		case <-s.mShow.ClickedCh:
			// runtime.Focus(s.ctx)
			s.ToggleVisibility()
		case <-mOptions.ClickedCh:
			runtime.Show(s.ctx)
		case <-mQuit.ClickedCh:
			systray.Quit()
		}
	}
}
