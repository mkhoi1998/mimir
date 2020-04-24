package chat

import (
	"fmt"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/mkhoi1998/mimir/cmd/handler"
	"github.com/mkhoi1998/mimir/errorer"
)

// Handler create a chat windows on terminal for chatting with bot
func Handler(args []string) error {
	input := make(chan string)
	output := make(chan []string)
	go processAnswer(input, output)
	return drawchat(input, output)
}

func drawchat(input chan string, output chan []string) error {
	// Create a new GUI.
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()
	g.Cursor = true

	// Update the views when terminal changes size.
	g.SetManagerFunc(func(g *gocui.Gui) error {
		termwidth, termheight := g.Size()
		_, err = g.SetView("output", 0, 0, termwidth-1, termheight-5)
		if err != nil {
			return err
		}
		_, err = g.SetView("input", 0, termheight-4, termwidth-1, termheight-1)
		if err != nil {
			return err
		}
		return nil
	})

	// Terminal width and height.
	termwidth, termheight := g.Size()

	// Output.
	ov, err := g.SetView("output", 0, 0, termwidth-1, termheight-4)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	// ov.FgColor = gocui.ColorRed
	ov.Autoscroll = true
	ov.Wrap = true
	ov.Frame = false

	// // Send a welcome message.
	// err = renderChat(ov, "Press Ctrl-C to quit.", 0, termwidth/3)
	// if err != nil {
	// 	return err
	// }

	// Input.
	iv, err := g.SetView("input", 0, termheight-3, termwidth-1, termheight-1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	iv.Wrap = true
	iv.Editable = true
	err = iv.SetCursor(0, 0)
	if err != nil {
		return err
	}

	// Bind Ctrl-C so the user can quit.
	err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})
	if err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlK, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(ov, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlJ, gocui.ModNone,
		func(g *gocui.Gui, res *gocui.View) error {
			scrollView(ov, 1)
			return nil
		}); err != nil {
		return err
	}

	// Bind enter key to input to send new messages.
	err = g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, iv *gocui.View) error {
		// Read buffer from the beginning.
		iv.Rewind()

		switch iv.Buffer() {
		case ":q\n":
			return gocui.ErrQuit
		case "":
			return nil
		}

		// Get output view and print.
		ov, err := g.View("output")
		if err != nil {
			return err
		}

		err = renderChat(ov, iv.Buffer(), 2*termwidth/3, termwidth)
		if err != nil {
			return err
		}

		input <- iv.Buffer()

		// Reset input.
		iv.Clear()

		// Reset cursor.
		err = iv.SetCursor(0, 0)
		if err != nil {
		}
		return err
	})
	if err != nil {
	}

	// Set the focus to input.
	_, err = g.SetCurrentView("input")
	if err != nil {
	}

	go func() {
		for {
			select {
			case mimirOutput := <-output:
				for i := range mimirOutput {
					if mimirOutput[i] == "" {
						break
					}
					if len(mimirOutput[i]) > 300 || i%5 == 0 {
						time.Sleep(5 * time.Second)
					}
					renderChat(ov, mimirOutput[i], 0, termwidth/2)
					// Refresh view
					g.Update(func(g *gocui.Gui) error {
						return nil
					})
				}
			}
		}
	}()

	// Start the main loop.
	err = g.MainLoop()

	return nil
}

func scrollView(v *gocui.View, dy int) {
	_, y := v.Size()
	ox, oy := v.Origin()

	if oy+dy > strings.Count(v.ViewBuffer(), "\n")-y-1 {
		v.Autoscroll = true
	} else {
		v.Autoscroll = false
		v.SetOrigin(ox, oy+dy)
	}
}

func renderChat(v *gocui.View, chat string, x0, x1 int) error {
	isUser := x0 != 0
	termwidth, _ := v.Size()
	width := x1 - x0
	chat = strings.ReplaceAll(chat, "( ", "(")
	chat = strings.ReplaceAll(chat, " )", ")")
	parts := strings.Split(strings.TrimSuffix(chat, "\n"), " ")
	var r string
	var res []string
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
		if len(r)+len(parts[i]) > width {
			res = append(res, r)
			r = ""
		}
		r += parts[i] + " "
	}
	res = append(res, r)
	if isUser {
		for i := range res {
			if i == 0 {
				res[i] = fmt.Sprintf("%v%v\033[1;36m%v\033[0m\n", strings.Repeat(" ", termwidth-len(res[i])-3), res[i], "Me")
				continue
			}
			res[i] = fmt.Sprintf("%v%v\n", strings.Repeat(" ", termwidth-len(res[i])-3), res[i])
		}
		_, err := fmt.Fprint(v, strings.Join(res, ""), "\n")
		return err
	}
	for i := range res {
		if i == 0 {
			res[i] = fmt.Sprintf("\033[1;36m%v\033[0m %v\n", "Mimir", res[i])
			continue
		}
		res[i] = fmt.Sprintf("%v%v\n", strings.Repeat(" ", 6), res[i])
	}
	_, err := fmt.Fprint(v, strings.Join(res, ""), "\n")
	return err

}

func processAnswer(input chan string, output chan []string) {
	for {
		select {
		case userInput := <-input:
			keywords := handler.ExtractKeywords(userInput)

			var res []string

			switch len(keywords) {
			case 0:
				res = append(res, errorer.ErrEmptyQuestion.Error())
			case 1:
				res = handler.GetStackWiki(keywords[0])
			case 2:
				res = strings.Split(handler.ExtractStackWiki(keywords), "\n")[1:]
			default:
				res = handler.GetStackoverflow(keywords)
			}
			output <- res
		}
	}

}
