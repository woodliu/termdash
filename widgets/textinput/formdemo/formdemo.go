// Copyright 2020 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Binary formdemo creates a form that accepts text inputs and supports
// keyboard navigation.
package main

import (
	"context"
	"fmt"
	"os/user"
	"time"

	"github.com/woodliu/termdash"
	"github.com/woodliu/termdash/align"
	"github.com/woodliu/termdash/cell"
	"github.com/woodliu/termdash/container"
	"github.com/woodliu/termdash/keyboard"
	"github.com/woodliu/termdash/linestyle"
	"github.com/woodliu/termdash/terminal/tcell"
	"github.com/woodliu/termdash/widgets/button"
	"github.com/woodliu/termdash/widgets/text"
	"github.com/woodliu/termdash/widgets/textinput"
)

// buttonChunks creates the text chunks for a button from the provided text.
func buttonChunks(text string) []*button.TextChunk {
	if len(text) == 0 {
		return nil
	}
	first := string(text[0])
	rest := string(text[1:])

	return []*button.TextChunk{
		button.NewChunk(
			"<",
			button.TextCellOpts(cell.FgColor(cell.ColorWhite)),
			button.FocusedTextCellOpts(cell.FgColor(cell.ColorBlack)),
			button.PressedTextCellOpts(cell.FgColor(cell.ColorBlack)),
		),
		button.NewChunk(
			first,
			button.TextCellOpts(cell.FgColor(cell.ColorRed)),
		),
		button.NewChunk(
			rest,
			button.TextCellOpts(cell.FgColor(cell.ColorWhite)),
			button.FocusedTextCellOpts(cell.FgColor(cell.ColorBlack)),
			button.PressedTextCellOpts(cell.FgColor(cell.ColorBlack)),
		),
		button.NewChunk(
			">",
			button.TextCellOpts(cell.FgColor(cell.ColorWhite)),
			button.FocusedTextCellOpts(cell.FgColor(cell.ColorBlack)),
			button.PressedTextCellOpts(cell.FgColor(cell.ColorBlack)),
		),
	}
}

// form contains the elements of a text input form.
type form struct {
	// userInput is a text input that accepts user name.
	userInput *textinput.TextInput
	// uidInput is a text input that accepts UID.
	uidInput *textinput.TextInput
	// gidInput is a text input that accepts GID.
	gidInput *textinput.TextInput
	// homeInput is a text input that accepts path to the home folder.
	homeInput *textinput.TextInput

	// submitB is a button that submits the form.
	submitB *button.Button
	// cancelB is a button that exist the application.
	cancelB *button.Button
}

// newForm returns a new form instance.
// The cancel argument is a function that terminates the application when called.
func newForm(cancel context.CancelFunc) (*form, error) {
	var username string
	u, err := user.Current()
	if err != nil {
		username = "mum4k"
	} else {
		username = u.Username
	}

	userInput, err := textinput.New(
		textinput.Label("Username: ", cell.FgColor(cell.ColorNumber(33))),
		textinput.DefaultText(username),
		textinput.MaxWidthCells(20),
		textinput.ExclusiveKeyboardOnFocus(),
	)
	uidInput, err := textinput.New(
		textinput.Label("UID:      ", cell.FgColor(cell.ColorNumber(33))),
		textinput.DefaultText("1000"),
		textinput.MaxWidthCells(20),
		textinput.ExclusiveKeyboardOnFocus(),
	)
	gidInput, err := textinput.New(
		textinput.Label("GID:      ", cell.FgColor(cell.ColorNumber(33))),
		textinput.DefaultText("1000"),
		textinput.MaxWidthCells(20),
		textinput.ExclusiveKeyboardOnFocus(),
	)
	homeInput, err := textinput.New(
		textinput.Label("Home:     ", cell.FgColor(cell.ColorNumber(33))),
		textinput.DefaultText(fmt.Sprintf("/home/%s", username)),
		textinput.MaxWidthCells(20),
		textinput.ExclusiveKeyboardOnFocus(),
	)

	submitB, err := button.NewFromChunks(buttonChunks("Submit"), nil,
		button.Key(keyboard.KeyEnter),
		button.GlobalKeys('s', 'S'),
		button.DisableShadow(),
		button.Height(1),
		button.TextHorizontalPadding(0),
		button.FillColor(cell.ColorBlack),
		button.FocusedFillColor(cell.ColorNumber(117)),
		button.PressedFillColor(cell.ColorNumber(220)),
	)
	if err != nil {
		panic(err)
	}

	cancelB, err := button.NewFromChunks(buttonChunks("Cancel"), func() error {
		cancel()
		return nil
	},
		button.FillColor(cell.ColorNumber(220)),
		button.Key(keyboard.KeyEnter),
		button.GlobalKeys('c', 'C'),
		button.DisableShadow(),
		button.Height(1),
		button.TextHorizontalPadding(0),
		button.FillColor(cell.ColorBlack),
		button.FocusedFillColor(cell.ColorNumber(117)),
		button.PressedFillColor(cell.ColorNumber(220)),
	)
	if err != nil {
		panic(err)
	}

	return &form{
		userInput: userInput,
		uidInput:  uidInput,
		gidInput:  gidInput,
		homeInput: homeInput,
		submitB:   submitB,
		cancelB:   cancelB,
	}, nil
}

// formLayout updates the container into a layout with text inputs and buttons.
func formLayout(c *container.Container, f *form) error {
	return c.Update("root",
		container.KeyFocusNext(keyboard.KeyTab),
		container.KeyFocusPrevious(keyboard.KeyBacktab),
		container.KeyFocusGroupsNext(keyboard.KeyArrowDown, 1),
		container.KeyFocusGroupsPrevious(keyboard.KeyArrowUp, 1),
		container.KeyFocusGroupsNext(keyboard.KeyArrowRight, 2),
		container.KeyFocusGroupsPrevious(keyboard.KeyArrowLeft, 2),
		container.SplitHorizontal(
			container.Top(
				container.Border(linestyle.Light),
				container.SplitHorizontal(
					container.Top(
						container.SplitHorizontal(
							container.Top(
								container.Focused(),
								container.KeyFocusGroups(1),
								container.PlaceWidget(f.userInput),
							),
							container.Bottom(
								container.KeyFocusGroups(1),
								container.KeyFocusSkip(),
								container.PlaceWidget(f.uidInput),
							),
						),
					),
					container.Bottom(
						container.SplitHorizontal(
							container.Top(
								container.KeyFocusGroups(1),
								container.KeyFocusSkip(),
								container.PlaceWidget(f.gidInput),
							),
							container.Bottom(
								container.KeyFocusGroups(1),
								container.KeyFocusSkip(),
								container.PlaceWidget(f.homeInput),
							),
						),
					),
				),
			),
			container.Bottom(
				container.SplitHorizontal(
					container.Top(
						container.SplitVertical(
							container.Left(
								container.KeyFocusGroups(1, 2),
								container.PlaceWidget(f.submitB),
								container.AlignHorizontal(align.HorizontalRight),
								container.PaddingRight(5),
							),
							container.Right(
								container.KeyFocusGroups(1, 2),
								container.PlaceWidget(f.cancelB),
								container.AlignHorizontal(align.HorizontalLeft),
								container.PaddingLeft(5),
							),
						),
					),
					container.Bottom(
						container.KeyFocusSkip(),
					),
					container.SplitFixed(3),
				),
			),
			container.SplitFixed(6),
		),
	)
}

// submitLayout updates the container into a layout that displays the submitted data.
// The cancel argument is a function that terminates Termdash when called.
func submitLayout(c *container.Container, f *form, cancel context.CancelFunc) error {
	t, err := text.New()
	if err != nil {
		return err
	}

	if err := t.Write("Submitted data:\n\n"); err != nil {
		return err
	}
	if err := t.Write(fmt.Sprintf("Username: %s\n", f.userInput.Read())); err != nil {
		return err
	}
	if err := t.Write(fmt.Sprintf("UID: %s\n", f.uidInput.Read())); err != nil {
		return err
	}
	if err := t.Write(fmt.Sprintf("GID: %s\n", f.gidInput.Read())); err != nil {
		return err
	}
	if err := t.Write(fmt.Sprintf("Home: %s\n", f.homeInput.Read())); err != nil {
		return err
	}

	okB, err := button.NewFromChunks(buttonChunks("OK"), func() error {
		cancel()
		return nil
	},
		button.FillColor(cell.ColorNumber(220)),
		button.Key(keyboard.KeyEnter),
		button.GlobalKeys('o', 'O'),
		button.DisableShadow(),
		button.Height(1),
		button.TextHorizontalPadding(0),
		button.FillColor(cell.ColorBlack),
		button.FocusedFillColor(cell.ColorNumber(117)),
		button.PressedFillColor(cell.ColorNumber(220)),
	)
	if err != nil {
		return err
	}

	return c.Update("root",
		container.SplitHorizontal(
			container.Top(
				container.SplitVertical(
					container.Left(),
					container.Right(
						container.PlaceWidget(t),
					),
					container.SplitPercent(33),
				),
			),
			container.Bottom(
				container.Focused(),
				container.PlaceWidget(okB),
			),
			container.SplitFixed(7),
		),
	)
}

func main() {
	t, err := tcell.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	ctx, cancel := context.WithCancel(context.Background())
	c, err := container.New(t, container.ID("root"))
	if err != nil {
		panic(err)
	}

	f, err := newForm(cancel)
	if err != nil {
		panic(err)
	}
	f.submitB.SetCallback(func() error {
		return submitLayout(c, f, cancel)
	})
	if err := formLayout(c, f); err != nil {
		panic(err)
	}

	if err := termdash.Run(ctx, t, c, termdash.RedrawInterval(100*time.Millisecond)); err != nil {
		panic(err)
	}
}
