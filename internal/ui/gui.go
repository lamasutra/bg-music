package ui

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/getlantern/systray"
	"github.com/lamasutra/bg-music/internal/app"
	"github.com/lamasutra/bg-music/internal/ui/gui"
	"github.com/lamasutra/bg-music/pkg/events"
	"github.com/lamasutra/bg-music/pkg/input"

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
	app            *app.AppState
	ctx            context.Context // app
	onStartup      func(a *app.AppState)
	assets         *embed.FS
	icon           []byte
	visible        bool
	playerControls *gui.PlayerControls
	inputManager   input.InputManager
	mShow          *systray.MenuItem
}

type WindowSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (s *guiState) startup(ctx context.Context) {
	s.ctx = ctx
	runtime.Hide(s.ctx)
	events.ListenAll("gui", s.eventDispatcher)
	events.Listen(bgplogger.EV_LOG, "gui", s.renderMessage)
	s.onStartup(s.app)
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
	if event == input.EV_INPUT_PRESSED {
		bgplogger.Info("eventDispatcher", event, values)
	}

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

func NewGui(a *app.AppState) *guiState {
	// Create an instance of the app structure
	app := &guiState{
		app:    a,
		assets: a.Assets,
		icon:   a.Icon,
	}

	return app
}

func (s *guiState) Run(onStartup func(a *app.AppState)) {
	s.onStartup = onStartup

	guiApp := s.createApplication()

	err := guiApp.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *guiState) MinimizeWindow() {
	runtime.WindowMinimise(s.ctx)
}

func (s *guiState) SetWindowSize(width int, height int) {
	runtime.WindowSetSize(s.ctx, width, height)
}

func (s *guiState) GetWindowSize() WindowSize {
	width, height := runtime.WindowGetSize(s.ctx)

	return WindowSize{Width: width, Height: height}
}

func (s *guiState) CloseApp() {
	runtime.Quit(s.ctx)
}

func (s *guiState) createApplication() *application.Application {
	s.playerControls = gui.NewPlayerControls()
	s.inputManager = input.NewManager(200 * time.Millisecond)

	// Create application with options
	app := application.NewWithOptions(
		&options.App{
			Title:             "Background Music Player",
			Width:             500,
			Height:            304,
			MinWidth:          500,
			MinHeight:         304,
			MaxWidth:          500,
			MaxHeight:         768,
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
				s.inputManager,
				s.app,
				s.app.Config,
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

	return app
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
			// runtime.EventsEmit(s.ctx, "navigate", "/options")
			runtime.Show(s.ctx)
			runtime.EventsEmit(s.ctx, "show-options", true)
			// runtime.Show(s.ctx)
		case <-mQuit.ClickedCh:
			systray.Quit()
		}
	}
}
