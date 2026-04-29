package ui

import (
	"fmt"

	"github.com/wagoodman/dive/runtime/ui/components"
	"github.com/wagoodman/dive/runtime/ui/key"
	"github.com/jroimartin/gocui"
	"github.com/sirupsen/logrus"
)

// UI is the main user interface manager for the dive application.
// It coordinates the layout, keybindings, and rendering of all views.
type UI struct {
	gui        *gocui.Gui
	controller *components.Controller
}

// New creates and initializes a new UI instance with the given controller.
func New(controller *components.Controller) (*UI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, fmt.Errorf("failed to create gui: %w", err)
	}

	u := &UI{
		gui:        g,
		controller: controller,
	}

	g.Cursor = false
	g.Mouse = true
	g.SetManagerFunc(u.layout)

	if err := u.bindKeys(); err != nil {
		g.Close()
		return nil, fmt.Errorf("failed to bind keys: %w", err)
	}

	return u, nil
}

// Run starts the main UI event loop. It blocks until the user quits.
func (u *UI) Run() error {
	defer u.gui.Close()

	if err := u.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return fmt.Errorf("ui main loop error: %w", err)
	}
	return nil
}

// Close releases all resources held by the UI.
func (u *UI) Close() {
	u.gui.Close()
}

// layout is the gocui layout manager function called on every render cycle.
func (u *UI) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if maxX == 0 || maxY == 0 {
		return nil
	}

	if err := u.controller.Layout(g, maxX, maxY); err != nil {
		logrus.Errorf("layout error: %v", err)
		return err
	}

	return nil
}

// bindKeys registers all global keybindings for the application.
func (u *UI) bindKeys() error {
	// Quit keybindings
	for _, k := range []interface{}{gocui.KeyCtrlC, 'q'} {
		if err := u.gui.SetKeybinding("", k, gocui.ModNone, u.quit); err != nil {
			return fmt.Errorf("failed to bind quit key %v: %w", k, err)
		}
	}

	// Tab to cycle between views
	if err := u.gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, u.nextView); err != nil {
		return fmt.Errorf("failed to bind tab key: %w", err)
	}

	// Register component-specific keybindings
	for _, binding := range u.controller.KeyBindings() {
		b := binding // capture range variable
		handler := func(g *gocui.Gui, v *gocui.View) error {
			return b.Handler(g, v)
		}
		if err := u.gui.SetKeybinding(b.ViewName, b.Key, b.Modifier, handler); err != nil {
			return fmt.Errorf("failed to bind key for view %q: %w", b.ViewName, err)
		}
	}

	return nil
}

// quit is the handler invoked when the user requests to exit the application.
func (u *UI) quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

// nextView cycles focus to the next registered view in the controller.
func (u *UI) nextView(g *gocui.Gui, _ *gocui.View) error {
	nextView := u.controller.NextView()
	if nextView == "" {
		return nil
	}
	if _, err := g.SetCurrentView(nextView); err != nil {
		logrus.Debugf("could not set view %q: %v", nextView, err)
	}
	return nil
}

// Refresh triggers a re-render of the UI.
func (u *UI) Refresh() {
	u.gui.Update(func(g *gocui.Gui) error {
		return u.layout(g)
	})
}

// KeyBindingSummary returns a human-readable summary of active keybindings
// suitable for display in a status bar or help overlay.
func (u *UI) KeyBindingSummary() string {
	return key.Summary(u.controller.KeyBindings())
}
