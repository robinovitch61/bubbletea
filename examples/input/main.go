package main

// A simple program that counts down from 5 and then exits.

import (
	"errors"
	"fmt"
	"log"

	"github.com/charmbracelet/tea"
	"github.com/charmbracelet/teaparty/input"
)

type Model struct {
	Input input.Model
	Error error
}

type tickMsg struct{}

func main() {
	tea.UseSysLog("tea")

	p := tea.NewProgram(
		initialize,
		update,
		view,
		subscriptions,
	)

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

func initialize() (tea.Model, tea.Cmd) {
	return Model{
		Input: input.DefaultModel(),
		Error: nil,
	}, nil
}

func update(msg tea.Msg, model tea.Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m, ok := model.(Model)
	if !ok {
		// When we encounter errors in Update we simply add the error to the
		// model so we can handle it in the view. We could also return a command
		// that does something else with the error, like logs it via IO.
		m.Error = errors.New("could not perform assertion on model")
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case tea.ErrMsg:
		m.Error = msg
		return m, nil
	}

	m.Input, cmd = input.Update(msg, m.Input)
	return m, cmd
}

func subscriptions(model tea.Model) tea.Subs {
	return tea.Subs{
		// We just hand off the subscription to the input component, giving
		// it the model it expects.
		"input": func(model tea.Model) tea.Msg {
			m, _ := model.(Model)
			return input.Blink(m.Input)
		},
	}
}

func view(model tea.Model) string {
	m, ok := model.(Model)
	if !ok {
		return "Oh no: could not perform assertion on model."
	} else if m.Error != nil {
		return fmt.Sprintf("Uh oh: %s", m.Error)
	}
	return fmt.Sprintf(
		"What’s your favorite Pokémon?\n\n%s\n\n%s",
		input.View(m.Input),
		"(esc to quit)",
	)
}