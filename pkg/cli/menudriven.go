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
	"github.com/pterm/pterm"
)

type Spinner struct {
	c pterm.SpinnerPrinter
	s *pterm.SpinnerPrinter
}

func GetSpinner() *Spinner {
	spinner := pterm.DefaultSpinner
	spinner.Sequence = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	return &Spinner{
		c: spinner,
	}
}

func (s *Spinner) Start(msg ...any) {
	s.s, _ = s.c.Start(msg...)
}

func (s *Spinner) StopWithSuccess(msg ...any) {
	s.s.Success(msg...)
}

func (s *Spinner) StopWithFailure(msg ...any) {
	s.s.Fail(msg...)
}

func Confirmation(prompt, defaultOption string) (proceed bool, err error) {
	x := pterm.DefaultInteractiveConfirm
	if len(defaultOption) != 0 {
		x = *x.WithDefaultText(defaultOption)
	}
	return x.Show(prompt)
}

func TextInput(prompt string) (string, error) {
	return pterm.DefaultInteractiveTextInput.Show(prompt)
}

func TextInputPassword(prompt string) (string, error) {
	x := pterm.DefaultInteractiveTextInput.WithMask("*")
	return x.Show(prompt)
}

func DropDown(prompt string, options map[string]string, defaultOption string) (string, error) {
	var _options []string
	for k := range options {
		_options = append(_options, k)
	}

	x := pterm.DefaultInteractiveSelect.WithOptions(_options)
	if len(defaultOption) != 0 {
		for k, v := range options {
			if v == defaultOption {
				defaultOption = k
				break
			}
		}
		x = x.WithDefaultOption(defaultOption)
	}

	if v, err := x.Show(prompt); err != nil {
		return "", err
	} else {
		// pterm.DefaultArea.Clear()
		return options[v], nil
	}

}
