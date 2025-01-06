package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/nikodem-wrona/encryptor/internal"
	"golang.org/x/term"
)

type model struct {
	files     []string
	cursor    int              // which to-do list item our cursor is pointing at
	selected  map[int]struct{} // which to-do items are selected
	password  string
	textInput textinput.Model
}

func initialModel(files []string) model {
	ti := textinput.New()

	ti.Placeholder = "Your password"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		files:     files,
		cursor:    0,
		selected:  make(map[int]struct{}),
		password:  "",
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func getPasswordInputStyle() lipgloss.Style {
	width, height := getTerminalSize()

	return lipgloss.
		NewStyle().
		Height(height).
		Width(width).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground()
}

func getFilePickerStyle() lipgloss.Style {
	width, height := getTerminalSize()

	return lipgloss.
		NewStyle().
		Height(height).
		Width(width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground().
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Padding(1, 1)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}

		case "enter", " ":
			if m.password == "" {
				m.password = m.textInput.Value()
				m.textInput.SetValue("")
				return m, nil
			}

			filepath := internal.GetFilePathByFileName(m.files[m.cursor])

			if internal.IsEncryptedFile(filepath) {
				decryptedFilePath := filepath[:len(filepath)-4]
				internal.DecryptFile(filepath, decryptedFilePath, m.password)
			} else {
				encryptedFilePath := internal.GetFilePathForEncryptedFile(filepath)
				internal.EncryptFile(filepath, encryptedFilePath, m.password)
			}

		default:
			if m.password == "" {
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
		}

		m.files = internal.FindAllFiledInCurrentDir()
	}

	return m, nil
}

func getTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		fmt.Println("Error getting terminal size")
		os.Exit(1)
	}

	return width - 2, height - 2
}

func (m model) View() string {
	if m.password == "" {
		style := getPasswordInputStyle()
		return fmt.Sprintf("%s", style.Render(m.textInput.View()))
	}

	s := "What file should be encrypted?\n\n"
	style := getFilePickerStyle()

	for i, choice := range m.files {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		if internal.IsEncryptedFile(internal.GetFilePathByFileName(choice)) {
			s += fmt.Sprintf("%s \033[31m%s\033[0m\n", cursor, choice)
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
	}

	s += "\nPress q to quit.\n"

	return fmt.Sprintf("%s", style.Render(s))
}

func main() {

	files := internal.FindAllFiledInCurrentDir()

	p := tea.NewProgram(initialModel(
		files,
	), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
