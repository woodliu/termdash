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

// Binary heatmapdemo displays a heatmap widget.
// Exist when 'q' is pressed.
package main

import (
	"context"
	"github.com/woodliu/termdash"
	"github.com/woodliu/termdash/container"
	"github.com/woodliu/termdash/linestyle"
	"github.com/woodliu/termdash/terminal/tcell"
	"github.com/woodliu/termdash/terminal/terminalapi"
	"github.com/woodliu/termdash/widgets/heatmap"
)

func main() {
	t, err := tcell.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	hp, err := heatmap.New()
	if err != nil {
		panic(err)
	}

	// TODO: set heatmap's data

	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		container.PlaceWidget(hp),
	)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter)); err != nil {
		panic(err)
	}
}
