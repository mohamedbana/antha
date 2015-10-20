---
layout: default
type: api
navgroup: api
shortname: manual/cli
title: manual/cli
microarch-api:
  published: 2015-06-25
  antha_version: 0.0.2
  package: manual/cli
---
# cli
--
    import "."


## Usage

#### type CUI

```go
type CUI struct {
	//CmdIn channel to receive input commands
	CmdIn chan CUICommandRequest

	//CmdOut channel to output the command results
	CmdOut chan CUICommandResult
	//LogIn channel to receive Log messages. This messages must be interpretable by the CUI. strings must be supported
	LogIn chan interface{}
	//Exit is a channel to be closed on exit to allow
	Exit chan interface{}
	//G reference to the underlying gocui GUI interface
	G *gocui.Gui
}
```

CUI represents a gocui interface for the manual driver commands with log input

#### func  NewCUI

```go
func NewCUI() *CUI
```

#### func (*CUI) ActionBack

```go
func (c *CUI) ActionBack(g *gocui.Gui, v *gocui.View) error
```
ActionBack goes back from the commadn dialogue without performing any operation

#### func (*CUI) ActionDone

```go
func (c *CUI) ActionDone(g *gocui.Gui, v *gocui.View) error
```
ActionDone reports a commandresult for the action with a positive feedback

#### func (*CUI) ActionError

```go
func (c *CUI) ActionError(g *gocui.Gui, v *gocui.View) error
```
ActionError displays a window to give a description of the error occurred on a
specific action

#### func (*CUI) Close

```go
func (c *CUI) Close()
```
Close waits for the user to exit the interface, then closes the gocui

#### func (*CUI) Init

```go
func (c *CUI) Init() error
```
Init instantiates the Gui, sets the layout, keybindings and general colours

#### func (*CUI) ReportError

```go
func (c *CUI) ReportError(g *gocui.Gui, v *gocui.View) error
```
ReportError instantiates a command result with the contents of the
ErrorMessageView as the message for the error, sends the error on the channel
and deletes the views on error reporting. Deletes the action from the window too

#### func (*CUI) RunCLI

```go
func (c *CUI) RunCLI() error
```

#### type CUICommandRequest

```go
type CUICommandRequest struct {
	Id      string
	Message MultiLevelMessage
}
```

CommandRequest represents a request to command line in which a question is made,
a set of options are presented and one of those options must be selected.
Expected is the option which should be selected for the command to execute
successfully

#### func  NewCUICommandRequest

```go
func NewCUICommandRequest(id string, message MultiLevelMessage) *CUICommandRequest
```

#### type CUICommandResult

```go
type CUICommandResult struct {
	Id    string
	Error error
}
```

CommandResult is the result of a CommandRequest. It has a bool representing
whether the action was performed successfully and an Answer holding a string
with a textual representation of the result should it be needed

#### func  NewCUICommandResult

```go
func NewCUICommandResult(id string, err error) *CUICommandResult
```

#### type MultiLevelMessage

```go
type MultiLevelMessage struct {
	Message  string
	Children []MultiLevelMessage
}
```

MultiLevelMessage represents an aggregating struct to hold indentable strings

#### func  NewMultiLevelMessage

```go
func NewMultiLevelMessage(message string, children []MultiLevelMessage) *MultiLevelMessage
```

#### func (*MultiLevelMessage) ChildrenText

```go
func (m *MultiLevelMessage) ChildrenText() string
```

#### func (*MultiLevelMessage) LeveledString

```go
func (m *MultiLevelMessage) LeveledString(level string, out *bytes.Buffer)
```
LeveledString prints in the given buffer an indented version of the
MultiLevelMessage using the level string as padding

#### func (*MultiLevelMessage) String

```go
func (m *MultiLevelMessage) String() string
```
String will return the whole MultiLevelMessage indented with spaces in a
multiline string
