package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type debugMenuDriven struct {
	progress ProgressAnimation
}

type debugSpinner struct {
	chars     []string
	done      chan bool
	active    bool
	startTime time.Time
}

func NewDebugMenuDriven() *debugMenuDriven {
	return &debugMenuDriven{}
}

func (p *debugMenuDriven) GetProgressAnimation() ProgressAnimation {
	if p.progress != nil {
		return p.progress
	}
	p.progress = &debugSpinner{
		chars:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		done:   make(chan bool),
		active: false,
	}

	return p.progress
}

func (s *debugSpinner) Start(msg ...any) {
	if s.active {
		return
	}
	s.done = make(chan bool)
	fmt.Println(msg...)
	s.active = true
	s.startTime = time.Now()

	go func() {
		for i := 0; ; i = (i + 1) % len(s.chars) {
			select {
			case <-s.done:
				fmt.Println() // Clear the spinner
				return
			default:
				elapsed := time.Since(s.startTime).Round(time.Second)
				fmt.Printf("%s %s", s.chars[i], elapsed)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (s *debugSpinner) StopWithSuccess(msg ...any) {
	if !s.active {
		return
	}
	fmt.Println(msg...)
	s.done <- true
	s.active = false
}

func (s *debugSpinner) StopWithFailure(msg ...any) {
	if !s.active {
		return
	}
	fmt.Println(msg...)
	s.done <- true
	s.active = false
}

func (s *debugSpinner) Stop() {
	if !s.active {
		return
	}
	s.done <- true
	s.active = false
}

func (p *debugMenuDriven) Confirmation(prompt string, opts ...func(*option) error) (proceed bool, err error) {
	o, err := processOptions(opts)
	if err != nil {
		return false, err
	}

	fmt.Println(prompt)
	fmt.Printf("Proceed? [y/N]{default: %s}: ", color.HiGreenString(o.defaultValue))
	var response string
	_, err = fmt.Scanln(&response)
	if err != nil {
		return false, err
	}
	if len(response) == 0 {
		if o.defaultValue != "" {
			response = o.defaultValue
		} else {
			return false, nil
		}
	}
	if response != "y" && response != "Y" && response != "yes" {
		return false, nil
	}

	fmt.Println()

	return true, nil
}

func (p *debugMenuDriven) TextInput(prompt string, opts ...func(*option) error) (string, error) {
	o, err := processOptions(opts)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s {default: %s}: ", prompt, color.HiGreenString(o.defaultValue))
	response, err := reader.ReadString('\n')
	if err != nil && io.EOF != err {
		return "", err
	}
	response = strings.TrimSpace(response)

	if len(response) == 0 {
		if o.defaultValue != "" {
			response = o.defaultValue
		}
	}

	fmt.Println("Got response:", response)
	fmt.Println()

	return response, nil
}

func (p *debugMenuDriven) TextInputPassword(prompt string) (string, error) {
	fmt.Println(prompt)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (p *debugMenuDriven) DropDown(prompt string, options map[string]string, opts ...func(*option) error) (string, error) {
	o, err := processOptions(opts)
	if err != nil {
		return "", err
	}

	fmt.Printf("%s {default: %s}\n", prompt, color.HiGreenString(o.defaultValue))
	fmt.Println("Options[make sure to enter value and not the key]:")
	for k, v := range options {
		fmt.Printf("%s: %s\n", k, color.HiCyanString(v))
	}

	var response string
	_, err = fmt.Scanln(&response)
	if err != nil {
		return "", err
	}

	if len(response) == 0 {
		if o.defaultValue != "" {
			response = o.defaultValue
		}
	}

	fmt.Println("Got response:", response)
	fmt.Println()

	return response, nil
}

func (p *debugMenuDriven) DropDownList(prompt string, options []string, opts ...func(*option) error) (string, error) {
	o, err := processOptions(opts)
	if err != nil {
		return "", err
	}

	fmt.Printf("%s {default: %s}\n", prompt, color.HiGreenString(o.defaultValue))
	fmt.Println("Options:")
	for _, v := range options {
		fmt.Println(color.HiCyanString(v))
	}

	var response string
	_, err = fmt.Scanln(&response)
	if err != nil {
		return "", err
	}

	if len(response) == 0 {
		if o.defaultValue != "" {
			response = o.defaultValue
		}
	}

	fmt.Println("Got response:", response)
	fmt.Println()

	return response, nil
}
