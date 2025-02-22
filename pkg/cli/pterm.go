// Copyright 2025 Ksctl Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

type ptermMenuDriven struct {
	progress ProgressAnimation
}

type spinner struct {
	c      pterm.SpinnerPrinter
	s      *pterm.SpinnerPrinter
	active bool
}

func NewPtermMenuDriven() *ptermMenuDriven {
	return &ptermMenuDriven{}
}

func (p *ptermMenuDriven) GetProgressAnimation() ProgressAnimation {
	if p.progress != nil {
		return p.progress
	}

	spin := pterm.DefaultSpinner
	spin.Sequence = []string{"⡏", "⡟", "⡿", "⢿", "⣻", "⣽", "⣾", "⣷", "⣯", "⣏"}

	p.progress = &spinner{
		c: spin,
	}
	return p.progress
}

func (s *spinner) Start(msg ...any) {
	if s.active {
		return
	}
	s.s, _ = s.c.Start(msg...)
}

func (s *spinner) StopWithSuccess(msg ...any) {
	s.s.Success(msg...)
	s.s = nil
	s.active = false
}

func (s *spinner) Stop() {
	_, _ = fmt.Fprint(os.Stderr, "\r"+strings.Repeat(" ", pterm.GetTerminalWidth())) // Clear the spinner
	_ = s.s.Stop()
	s.s = nil
	s.active = false
}

func (s *spinner) StopWithFailure(msg ...any) {
	s.s.Fail(msg...)
	s.s = nil
	s.active = false
}

func (p *ptermMenuDriven) Confirmation(prompt string, opts ...func(*option) error) (proceed bool, err error) {
	o, err := processOptions(opts)
	if err != nil {
		return false, err
	}

	x := pterm.DefaultInteractiveConfirm
	if len(o.defaultValue) != 0 {
		x = *x.WithDefaultValue(o.defaultValue == "yes")
	}
	return x.Show(prompt)
}

func (p *ptermMenuDriven) TextInput(prompt string, opts ...func(*option) error) (string, error) {
	o, err := processOptions(opts)
	if err != nil {
		return "", err
	}

	if len(o.defaultValue) == 0 {
		return pterm.DefaultInteractiveTextInput.Show(prompt)
	}
	x := pterm.DefaultInteractiveTextInput.WithDefaultValue(o.defaultValue)
	return x.Show(prompt)
}

func (p *ptermMenuDriven) TextInputPassword(prompt string) (string, error) {
	x := pterm.DefaultInteractiveTextInput.WithMask("*")
	return x.Show(prompt)
}

func (p *ptermMenuDriven) DropDown(prompt string, options map[string]string, opts ...func(*option) error) (string, error) {
	o, err := processOptions(opts)
	if err != nil {
		return "", err
	}

	var _options []string
	for k := range options {
		_options = append(_options, k)
	}

	x := pterm.DefaultInteractiveSelect.WithOptions(_options)
	if len(o.defaultValue) != 0 {
		for k, v := range options {
			if v == o.defaultValue {
				o.defaultValue = k
				break
			}
		}
		x = x.WithDefaultOption(o.defaultValue)
	}

	if v, err := x.Show(prompt); err != nil {
		return "", err
	} else {
		return options[v], nil
	}

}

func (p *ptermMenuDriven) DropDownList(prompt string, options []string, opts ...func(*option) error) (string, error) {
	o, err := processOptions(opts)
	if err != nil {
		return "", err
	}

	x := pterm.DefaultInteractiveSelect.WithOptions(options)
	if len(o.defaultValue) != 0 {
		x = x.WithDefaultOption(o.defaultValue)
	}

	if v, err := x.Show(prompt); err != nil {
		return "", err
	} else {
		return v, nil
	}
}
