// /equipment/manual/cli/cui.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK

package cli

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/antha-lang/antha/microArch/logger"
	"github.com/antha-lang/antha/internal/github.com/jroimartin/gocui"
)

//MultiLevelMessage represents an aggregating struct to hold indentable strings
type MultiLevelMessage struct {
	Message  string
	Children []MultiLevelMessage
}

//LeveledString prints in the given buffer an indented version of the MultiLevelMessage using the level string as padding
func (m *MultiLevelMessage) LeveledString(level string, out *bytes.Buffer) {
	out.WriteString(fmt.Sprintf("%s%s\n", level, m.Message))
	for _, v := range m.Children {
		v.LeveledString(level+level, out)
	}
}

//String will return the whole MultiLevelMessage indented with spaces in a multiline string
func (m *MultiLevelMessage) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("%s\n", m.Message))
	for _, v := range m.Children {
		v.LeveledString("  ", &out)
	}
	return out.String()
}
func NewMultiLevelMessage(message string, children []MultiLevelMessage) *MultiLevelMessage {
	mlm := new(MultiLevelMessage)
	mlm.Message = message
	mlm.Children = children
	return mlm
}
func (m *MultiLevelMessage) ChildrenText() string {
	var out bytes.Buffer
	for _, v := range m.Children {
		v.LeveledString("  ", &out)
	}
	return out.String()
}

//CommandRequest represents a request to command line in which a question is made, a set of options are presented and
// one of those options must be selected. Expected is the option which should be selected for the command to execute
// successfully
type CUICommandRequest struct {
	Id      string
	Message MultiLevelMessage
}

func NewCUICommandRequest(id string, message MultiLevelMessage) *CUICommandRequest {
	cr := new(CUICommandRequest)
	cr.Id = id
	cr.Message = message
	return cr
}

//CommandResult is the result of a CommandRequest. It has a bool representing whether the action was performed successfully
// and an Answer holding a string with a textual representation of the result should it be needed
type CUICommandResult struct {
	Id    string
	Error error
}

func NewCUICommandResult(id string, err error) *CUICommandResult {
	cr := new(CUICommandResult)
	cr.Id = id
	cr.Error = err
	return cr
}

//CUI represents a gocui interface for the manual driver commands with log input
type CUI struct {
	//CmdIn channel to receive input commands
	CmdIn chan CUICommandRequest
	//cmdList internal slice with the list of command requests as where received
	cmdList []CUICommandRequest
	//selectedCommand holds the value for the currently selected command in the action list view
	selectedCommand int
	//CmdOut channel to output the command results
	CmdOut chan CUICommandResult
	//LogIn channel to receive Log messages. This messages must be interpretable by the CUI. strings must be supported
	LogIn chan interface{}
	//Exit is a channel to be closed on exit to allow
	Exit chan interface{}
	//G reference to the underlying gocui GUI interface
	G              *gocui.Gui
	capturedstdout *os.File
}

func NewCUI() *CUI {
	c := new(CUI)
	c.LogIn = make(chan interface{})
	c.CmdIn = make(chan CUICommandRequest)
	c.CmdOut = make(chan CUICommandResult)
	c.Exit = make(chan interface{})
	c.cmdList = make([]CUICommandRequest, 0)
	c.selectedCommand = -1
	c.G = gocui.NewGui()
	return c
}

//Init instantiates the Gui, sets the layout, keybindings and general colours
func (c *CUI) Init() error {
	if err := c.G.Init(); err != nil {
		return err
	}
	c.G.SetLayout(c.layout)
	if err := c.keybindings(c.G); err != nil {
		return err
	}

	//Good looks
	c.G.SelBgColor = gocui.ColorGreen
	c.G.SelFgColor = gocui.ColorBlack
	c.G.ShowCursor = true

	return nil
}

//Close waits for the user to exit the interface, then closes the gocui
func (c *CUI) Close() {
	<-c.Exit
	c.G.Close()
	os.Stdout = c.capturedstdout
}
func (c *CUI) RunCLI() error {
	//Capture stdout and save it for shutdown
	c.capturedstdout = os.Stdout
	os.Stdout = nil

	go func() {
		defer c.Close()
		go func() {
			for v := range c.CmdIn {
				c.cmdList = append(c.cmdList, v)
				err := c.newAction(v)
				if err != nil {
					log.Panicln(err)
				}
			}
		}()

		go func() {
			for v := range c.LogIn {
				err := c.newLog(v)
				if err != nil {
					log.Panicln(err)
				}
			}
		}()
		ticker := time.NewTicker(time.Millisecond * 200) //Refresh rate
		go func() {
			//TODO we should implement this in a fork of gocui in the event processing loop
			for _ = range ticker.C {
				c.G.Flush()
			}
		}()
		err := c.G.MainLoop()
		if err != nil && err != gocui.Quit {
			panic(err)
		}
	}()
	return nil
}

//getCursorListPos returns the position of selection in the given view
func (c *CUI) getCursorListPos(g *gocui.Gui, v *gocui.View) (y int) {
	_, oy := v.Origin()
	_, cy := v.Cursor()
	y = oy + cy
	return
}

//cursorDown moves the cursor down one position. When reaching the bottom, windows the text
func (c *CUI) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

//cursorUp moves the cursor up one position. When reaching the top, windows the text
func (c *CUI) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

//cmdDone deletes the selected command from the list and redraws the list
func (c *CUI) cmdDone(i int) error {
	if len(c.cmdList) == 1 {
		if i != 0 {
			return errors.New(fmt.Sprintf("Impossible to delete command %d from CmdList", i))
		}
		c.cmdList = make([]CUICommandRequest, 0)
		c.selectedCommand = -1
	} else if i == 0 {
		c.cmdList = c.cmdList[1:]
	} else {
		c.cmdList = append(c.cmdList[0:i], c.cmdList[i+1:]...)
	}
	c.selectedCommand = -1
	//get cmd list view and redraw
	cmdlistview, err := c.G.View("ActionListView")
	if err != nil {
		return err
	}
	cmdlistview.Clear()
	cmdlistview.SetCursor(0, 0)
	for _, a := range c.cmdList {
		c.newAction(a)
	}
	return nil
}

//ActionError displays a window to give a description of the error occurred on a specific action
func (c *CUI) ActionError(g *gocui.Gui, v *gocui.View) error {
	if c.selectedCommand < 0 {
		panic(errors.New("Invalid command selection value"))
	}
	//send an error on

	maxX, maxY := g.Size()
	if v, err := g.SetView("ErrorMessageTittleView", maxX/2-30, maxY/4-2, maxX/2+30, maxY/4); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, " Write a description for the error (<Enter> when done)")
		v.FgColor = gocui.ColorRed
		if err := g.SetCurrentView("ActionView"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("ErrorMessageView", maxX/2-30, maxY/4, maxX/2+30, 3*maxY/4); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		if err := g.SetCurrentView("ActionView"); err != nil {
			return err
		}
		v.Editable = true
	}
	if err := g.SetCurrentView("ErrorMessageView"); err != nil {
		return err
	}

	return nil
}

//ReportError instantiates a command result with the contents of the ErrorMessageView as the message for the error,
// sends the error on the channel and deletes the views on error reporting. Deletes the action from the window too
func (c *CUI) ReportError(g *gocui.Gui, v *gocui.View) error {
	res := *NewCUICommandResult(c.cmdList[c.selectedCommand].Id, errors.New(v.Buffer())) //TODO strip \n ??
	c.CmdOut <- res
	err := c.cmdDone(c.selectedCommand)
	if err != nil {
		return err
	}
	//find all the views and exit
	if err := c.G.DeleteView("ErrorMessageTittleView"); err != nil {
		return err
	}
	if err := c.G.DeleteView("ErrorMessageView"); err != nil {
		return err
	}
	return c.deleteView(g, v)
}

//ActionBack goes back from the commadn dialogue without performing any operation
func (c *CUI) ActionBack(g *gocui.Gui, v *gocui.View) error {
	return c.deleteView(g, v)
}

//ActionDone reports a commandresult for the action with a positive feedback
func (c *CUI) ActionDone(g *gocui.Gui, v *gocui.View) error {
	if c.selectedCommand < 0 {
		panic(errors.New("Invalid command selection value"))
	}
	res := *NewCUICommandResult(c.cmdList[c.selectedCommand].Id, nil)
	c.CmdOut <- res
	c.cmdDone(c.selectedCommand)
	return c.deleteView(g, v)
}

//deleteView deletes ActionView + ActionViewTittle and gives focus to ActionListView
func (c *CUI) deleteView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("ActionView"); err != nil {
		return err
	}
	if err := g.DeleteView("ActionViewTitle"); err != nil {
		return err
	}
	if err := g.SetCurrentView("ActionListView"); err != nil {
		return err
	}
	if err := c.printGeneralHelp(g); err != nil {
		return err
	}
	return nil
}

//selectAction shows information about the action that is selected in the ActionListView
func (c *CUI) selectAction(g *gocui.Gui, v *gocui.View) error {
	c.selectedCommand = c.getCursorListPos(g, v)

	if c.selectedCommand > len(c.cmdList)-1 {
		return nil
	}

	var l string
	var t string
	//Load the message from the commandList
	l = c.cmdList[c.selectedCommand].Message.ChildrenText()
	t = c.cmdList[c.selectedCommand].Message.Message

	maxX, maxY := g.Size()
	if v, err := g.SetView("ActionViewTitle", maxX/2-30, maxY/3-2, maxX/2+30, maxY/3); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, t)
	}
	if v, err := g.SetView("ActionView", maxX/2-30, maxY/3, maxX/2+30, 2*maxY/3); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Wrap = true
		fmt.Fprintln(v, l)
		if err := g.SetCurrentView("ActionView"); err != nil {
			return err
		}
	}
	err := c.printActionHelp(g)
	return err
}

//layout sets the initial layou of the gui putting the windows in place
func (c *CUI) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("ActionListTitle", -1, -1, int(0.7*float32(maxX)), 2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprint(v, " Action List")
		v.FgColor = gocui.ColorGreen
	}
	if v, err := g.SetView("ActionListView", -1, 2, int(0.7*float32(maxX)), maxY-5); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Highlight = true
		if err := g.SetCurrentView("ActionListView"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("LogViewTitle", int(0.7*float32(maxX)), -1, maxX, 2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintf(v, " Log View")
		v.FgColor = gocui.ColorGreen
	}
	if v, err := g.SetView("LogView", int(0.7*float32(maxX)), 2, maxX, maxY-5); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Highlight = true
	}
	if v, err := g.SetView("HelpView", 0, maxY-5, maxX, maxY-1); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Frame = false
		v.FgColor = gocui.ColorGreen
		if err := c.printGeneralHelp(g); err != nil {
			return err
		}
	}
	return nil
}

func (c *CUI) printGeneralHelp(g *gocui.Gui) error {
	v, err := g.View("HelpView")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprint(v, "Use up/down keys to highlight Actions. <Enter> to select. <Tab> to change to Log Messages List.\n<Ctrl-X> to exit at any time")
	return nil
}
func (c *CUI) printActionHelp(g *gocui.Gui) error {
	v, err := g.View("HelpView")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprint(v, "Use <Enter> to acknowledge action. <Backspace> to go Back. <Space> to report Error.\n<Ctrl-X> to exit at any time")
	return nil
}

//quit is called when the user wants to finish the execution of the gui. It should display a warning screen and give an
// error result to all the pending actions. If no pending actions exists, it just exists
func (c *CUI) quit(g *gocui.Gui, v *gocui.View) error {
	if len(c.cmdList) > 0 {
		maxX, maxY := g.Size()
		if v, err := g.SetView("QuitWarnView", maxX/2-30, maxY/4, maxX/2+30, 3*maxY/4); err != nil {
			if err != gocui.ErrorUnkView {
				return err
			}
			fmt.Fprintf(v, "By pressing <Ctrl-X> all pending actions will be cancelled.\nPress <Enter> to go back")
			v.BgColor = gocui.ColorYellow
			v.FgColor = gocui.ColorBlack
			if err := g.SetCurrentView("QuitWarnView"); err != nil {
				return err
			}
		}
		return nil
	}
	return c.quitAcknowledged(g, v)
}

//quitAcknowledge exits the gui
func (c *CUI) quitAcknowledged(g *gocui.Gui, v *gocui.View) error {
	if err := c.cancelPendingActions(); err != nil {
		return err
	}
	close(c.Exit)
	//	return gocui.Quit
	return nil
}

//cancelPendingActions will output a result for every pending action with a cancelled error
func (c *CUI) cancelPendingActions() error {
	for _, action := range c.cmdList {
		res := *NewCUICommandResult(action.Id, errors.New("User aborted protocol."))
		c.CmdOut <- res
	}
	return nil
}

//abortQuit goes back to the ActionListView if the quit message is rejected
func (c *CUI) abortQuit(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("QuitWarnView"); err != nil {
		return err
	}
	if err := g.SetCurrentView("ActionListView"); err != nil {
		return err
	}
	return nil
}

//nextView jumps between the ActionListView and the LogView
func (c *CUI) nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "ActionListView" {
		return g.SetCurrentView("LogView")
	}
	return g.SetCurrentView("ActionListView")
}

//keybindings configures the keyboard actions for each window
func (c *CUI) keybindings(g *gocui.Gui) error {
	//Everybody
	if err := g.SetKeybinding("", gocui.KeyCtrlX, gocui.ModNone, c.quit); err != nil {
		return err
	}

	//ActionListView
	if err := g.SetKeybinding("ActionListView", gocui.KeyTab, gocui.ModNone, c.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("ActionListView", gocui.KeyArrowDown, gocui.ModNone, c.cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("ActionListView", gocui.KeyArrowUp, gocui.ModNone, c.cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("ActionListView", gocui.KeyEnter, gocui.ModNone, c.selectAction); err != nil {
		return err
	}

	//LogView
	if err := g.SetKeybinding("LogView", gocui.KeyArrowDown, gocui.ModNone, c.cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("LogView", gocui.KeyArrowUp, gocui.ModNone, c.cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("LogView", gocui.KeyTab, gocui.ModNone, c.nextView); err != nil {
		return err
	}

	//Action View
	if err := g.SetKeybinding("ActionView", gocui.KeyEnter, gocui.ModNone, c.ActionDone); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("ActionView", gocui.KeyBackspace, gocui.ModNone, c.ActionBack); err != nil {
		return err
	}
	if err := g.SetKeybinding("ActionView", gocui.KeyBackspace2, gocui.ModNone, c.ActionBack); err != nil {
		return err
	}
	if err := g.SetKeybinding("ActionView", gocui.KeySpace, gocui.ModNone, c.ActionError); err != nil {
		return err
	}

	if err := g.SetKeybinding("ErrorMessageView", gocui.KeyEnter, gocui.ModNone, c.ReportError); err != nil {
		return err
	}

	if err := g.SetKeybinding("QuitWarnView", gocui.KeyCtrlX, gocui.ModNone, c.quitAcknowledged); err != nil {
		return err
	}
	if err := g.SetKeybinding("QuitWarnView", gocui.KeyEnter, gocui.ModNone, c.abortQuit); err != nil {
		return err
	}

	return nil
}

//newAction inserts a new action from a command request in the command list
func (c *CUI) newAction(action CUICommandRequest) error {
	v, err := c.G.View("ActionListView")
	if err != nil {
		return err
	}
	v.Clear()
	for _, val := range c.cmdList {
		fmt.Fprint(v, val.Message.Message+"\n") //TODO strip new lines if they exist
	}
	return nil
}

//newLog inserts a new log in the logView
func (c *CUI) newLog(log interface{}) error {
	var shortDesc string
	switch l := log.(type) {
	case string:
		shortDesc = l
	case logger.LogEntry:
		logview, err := c.G.View("LogView")
		if err != nil {
			return err
		}
		lx, _ := logview.Size()
		var mess string
		if len(l.Message) > lx-9 {
			mess = l.Message[:lx-9]
		} else {
			mess = l.Message
		}
		shortDesc = fmt.Sprintf("[%s] %s...", l.Level.String()[:1], mess)
	default:
		//ignore by default
		return nil
	}
	v, err := c.G.View("LogView")
	if err != nil {
		return err
	}
	fmt.Fprint(v, shortDesc+"\n")
	return nil
}
