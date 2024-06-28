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

package tcell

import (
	"testing"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/kylelemons/godebug/pretty"
	"github.com/woodliu/termdash/cell"
	"github.com/woodliu/termdash/terminal/terminalapi"
)

func TestNewTerminalColorMode(t *testing.T) {
	tests := []struct {
		desc string
		opts []Option
		want *Terminal
	}{
		{
			desc: "default options",
			want: &Terminal{
				colorMode: terminalapi.ColorMode256,
			},
		},
		{
			desc: "sets color mode",
			opts: []Option{
				ColorMode(terminalapi.ColorModeNormal),
			},
			want: &Terminal{
				colorMode: terminalapi.ColorModeNormal,
			},
		},
	}

	tcellNewScreen = func() (tcell.Screen, error) { return nil, nil }
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := newTerminal(tc.opts...)
			if err != nil {
				t.Errorf("newTerminal => unexpected error:\n%v", err)
				return
			}

			// Ignore these fields.
			got.screen = nil
			got.events = nil
			got.done = nil
			got.clearStyle = nil

			if diff := pretty.Compare(tc.want, got); diff != "" {
				t.Errorf("newTerminal => unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestNewTerminalClearStyle(t *testing.T) {
	tests := []struct {
		desc string
		opts []Option
		want *Terminal
	}{
		{
			desc: "default options",
			want: &Terminal{
				colorMode: terminalapi.ColorMode256,
				clearStyle: &cell.Options{
					FgColor: cell.ColorDefault,
					BgColor: cell.ColorDefault,
				},
			},
		},
		{
			desc: "sets clear style",
			opts: []Option{
				ClearStyle(cell.ColorRed, cell.ColorBlue),
			},
			want: &Terminal{
				colorMode: terminalapi.ColorMode256,
				clearStyle: &cell.Options{
					FgColor: cell.ColorRed,
					BgColor: cell.ColorBlue,
				},
			},
		},
	}

	tcellNewScreen = func() (tcell.Screen, error) { return nil, nil }
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := newTerminal(tc.opts...)
			if err != nil {
				t.Errorf("newTerminal => unexpected error:\n%v", err)
				return
			}

			// Ignore these fields.
			got.screen = nil
			got.events = nil
			got.done = nil

			if diff := pretty.Compare(tc.want, got); diff != "" {
				t.Errorf("newTerminal => unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}
