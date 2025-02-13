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

type ProgressAnimation interface {
	Start(msg ...any)
	StopWithSuccess(msg ...any)
	StopWithFailure(msg ...any)
	Stop()
}

type option struct {
	defaultValue string
}

func WithDefaultValue(defaultValue string) func(*option) error {
	return func(o *option) error {
		o.defaultValue = defaultValue
		return nil
	}
}

type MenuDriven interface {
	GetProgressAnimation() ProgressAnimation
	Confirmation(prompt string, opts ...func(*option) error) (proceed bool, err error)
	TextInput(prompt string, opts ...func(*option) error) (string, error)
	TextInputPassword(prompt string) (string, error)
	DropDown(prompt string, options map[string]string, opts ...func(*option) error) (string, error)
	DropDownList(prompt string, options []string, opts ...func(*option) error) (string, error)
}

func processOptions(opts []func(*option) error) (option, error) {
	var o option
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return o, err
		}
	}
	return o, nil
}
