package chatui

import (
	"fmt"
	"strings"

	"github.com/augustkang/gochat/client/pkg/chatapp"
	"github.com/marcusolsson/tui-go"
)

func GetUI(userName string, app *chatapp.App) (ui tui.UI, cbox *tui.Box) {

	history := tui.NewVBox()
	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	username := strings.TrimSuffix(userName, "\n")
	input.OnSubmit(func(e *tui.Entry) {
		history.Append(tui.NewHBox(
			tui.NewLabel(username+" : "+e.Text()),
			tui.NewSpacer(),
		))
		err := app.WriteToConn(e.Text() + "\n")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		input.SetText("")
	})

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Maximum)

	ui, err := tui.New(chat)
	if err != nil {
		fmt.Println("failed to set ui", err)
		panic(err)
	}
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	return ui, history
}
