package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"

	"github.com/mattn/go-runewidth"
)

type Ticker struct {
	Symbol string
	Name   string
}

type TickerSelected struct {
	Ticker string
}
type getTickerSuggestions struct{ Fragment string }

// TickerPicker is a form input field.
type TickerPicker struct {
	value *string
	key   string

	// customization
	title       string
	description string
	inline      bool
	tickers     []Ticker

	// error handling
	validate func(string) error
	err      error

	// model
	textinput textinput.Model
	keymap    PickerKeyMap

	// state
	focused bool

	// options
	width      int
	height     int
	accessible bool
	theme      *huh.Theme
	lookup     func(string) ([]Ticker, error)
}

// NewTickerPicker returns a new input field.
func NewTickerPicker(tickerLookup func(string) ([]Ticker, error)) *TickerPicker {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Prompt = "Enter a ticker: "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Focus()
	ti.CharLimit = 8
	ti.Width = 80
	ti.ShowSuggestions = true

	i := &TickerPicker{
		value:     new(string),
		validate:  func(string) error { return nil },
		theme:     huh.ThemeCharm(),
		textinput: ti,
		keymap:    DefaultPickerKeyMap,
		tickers:   []Ticker{},
		lookup:    tickerLookup,
	}

	return i
}

// Value sets the value of the input field.
func (i *TickerPicker) Value(value *string) *TickerPicker {
	i.value = value
	i.textinput.SetValue(*value)
	return i
}

// Key sets the key of the input field.
func (i *TickerPicker) Key(key string) *TickerPicker {
	i.key = key
	return i
}

// Title sets the title of the input field.
func (i *TickerPicker) Title(title string) *TickerPicker {
	i.title = title
	return i
}

// Description sets the description of the input field.
func (i *TickerPicker) Description(description string) *TickerPicker {
	i.description = description
	return i
}

// Inline sets whether the title and input should be on the same line.
func (i *TickerPicker) Inline(inline bool) *TickerPicker {
	i.inline = inline
	return i
}

// Validate sets the validation function of the input field.
func (i *TickerPicker) Validate(validate func(string) error) *TickerPicker {
	i.validate = validate
	return i
}

// Error returns the error of the input field.
func (i *TickerPicker) Error() error {
	return i.err
}

// Skip returns whether the input should be skipped or should be blocking.
func (*TickerPicker) Skip() bool {
	return false
}

// Zoom returns whether the input should be zoomed.
func (*TickerPicker) Zoom() bool {
	return true
}

// Focus focuses the input field.
func (i *TickerPicker) Focus() tea.Cmd {
	i.focused = true
	return tea.Batch(i.textinput.Focus(), i.textinput.Focus())
}

// Blur blurs the input field.
func (i *TickerPicker) Blur() tea.Cmd {
	i.focused = false
	*i.value = i.textinput.Value()
	i.textinput.Blur()
	return textinput.Blink
}

// KeyBinds returns the help message for the input field.
func (i *TickerPicker) KeyBinds() []key.Binding {
	return i.keymap.ShortHelp()
}

// Init initializes the input field.
func (i *TickerPicker) Init() tea.Cmd {
	i.textinput.Blur()
	return nil
}

// Update updates the input field.
func (m *TickerPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickerSelected:
		*m.value = msg.Ticker
		m.textinput.SetValue(msg.Ticker)
		return m, tea.Sequence(m.Blur(), huh.NextField)
	case getTickerSuggestions:
		ts, err := m.lookup(msg.Fragment)
		if err != nil {
			return m, nil
		}

		suggestions := make([]string, len(ts))
		for i, t := range ts {
			suggestions[i] = t.Symbol
		}

		m.textinput.SetSuggestions(suggestions)
		m.tickers = ts

		var tiUpdateCmd, taUpdateCmd tea.Cmd
		m.textinput, tiUpdateCmd = m.textinput.Update(msg)
		return m, tea.Batch(tiUpdateCmd, taUpdateCmd)
	case tea.KeyMsg:
		switch {
		case msg.String() == "enter":
			value := m.textinput.Value()
			m.err = m.validate(value)
			if m.err != nil {
				return m, nil
			}

			return m, cmdize(TickerSelected{Ticker: value})
		}

		m.err = nil
		var tiUpdateCmd tea.Cmd
		m.textinput, tiUpdateCmd = m.textinput.Update(msg)
		lookup := m.textinput.Value()
		cmd := cmdize(getTickerSuggestions{Fragment: lookup})
		return m, tea.Batch(tiUpdateCmd, cmd)
	}

	ti, ticmd := m.textinput.Update(msg)
	m.textinput = ti

	return m, tea.Batch(ticmd)
}

// View renders the input field.
func (m *TickerPicker) View() string {
	styles := m.theme.Blurred
	if m.focused {
		styles = m.theme.Focused
	}

	var sb strings.Builder
	if m.title != "" {
		sb.WriteString(styles.Title.Render(m.title))
		if !m.inline {
			sb.WriteString("\n")
		}
	}
	if m.description != "" {
		sb.WriteString(styles.Description.Render(m.description))
		if !m.inline {
			sb.WriteString("\n")
		}
	}

	sb.WriteString(m.textinput.View() + "\n")

	if m.focused && len(m.textinput.Value()) > 0 && len(m.tickers) > 0 {
		allowedSuggestions := clamp(m.height-3, m.height-3, len(m.tickers))
		allowedWidth := clamp(m.width-2, 1, 80)
		suggested := m.GetSuggestedTicker()

		for _, t := range m.tickers[:allowedSuggestions] {
			style := m.theme.Focused.UnselectedOption
			if suggested != nil && suggested.Symbol == t.Symbol {
				style = m.theme.Focused.SelectedOption
			}
			t := strings.TrimSpace(style.Render(fmt.Sprintf("%s %s", t.Symbol, t.Name)))
			truncated := truncateText(t, allowedWidth, true)
			sb.WriteString(truncated + "\n")
		}
	}

	return styles.Base.Render(sb.String())
}

func (m *TickerPicker) GetSuggestedTicker() *Ticker {
	var result *Ticker
	defer func() {
		if r := recover(); r != nil {
			result = nil // workaround for panic bug in textinput bubble
		}
	}()
	suggested := m.textinput.CurrentSuggestion()

	if len(m.tickers) == 0 {
		return nil
	}

	for _, t := range m.tickers {
		if t.Symbol == suggested {
			result = &t
			break
		}
	}

	return result
}

// Run runs the input field in accessible mode.
func (i *TickerPicker) Run() error {
	if i.accessible {
		return i.runAccessible()
	}
	return huh.Run(i)
}

// runAccessible runs the input field in accessible mode.
func (i *TickerPicker) runAccessible() error {
	fmt.Println(i.theme.Blurred.Base.Render(i.theme.Focused.Title.Render(i.title)))
	fmt.Println()
	*i.value = accessibility.PromptString("Input: ", i.validate)
	fmt.Println(i.theme.Focused.SelectedOption.Render("Input: " + *i.value + "\n"))
	return nil
}

// TODO this func exposed on the huh.Field interface breaks extensibility because it is caked onto types only implemented by huh
func (i *TickerPicker) WithKeyMap(k *huh.KeyMap) huh.Field {
	return i
}

// WithAccessible sets the accessible mode of the input field.
func (i *TickerPicker) WithAccessible(accessible bool) huh.Field {
	i.accessible = accessible
	return i
}

// WithTheme sets the theme of the input field.
func (i *TickerPicker) WithTheme(theme *huh.Theme) huh.Field {
	i.theme = theme
	return i
}

// WithWidth sets the width of the input field.
func (i *TickerPicker) WithWidth(width int) huh.Field {
	i.width = width
	frameSize := i.theme.Blurred.Base.GetHorizontalFrameSize()
	promptWidth := lipgloss.Width(i.textinput.PromptStyle.Render(i.textinput.Prompt))
	titleWidth := lipgloss.Width(i.theme.Focused.Title.Render(i.title))
	descriptionWidth := lipgloss.Width(i.theme.Focused.Description.Render(i.description))
	i.textinput.Width = width - frameSize - promptWidth - 1
	if i.inline {
		i.textinput.Width -= titleWidth
		i.textinput.Width -= descriptionWidth
	}
	return i
}

// WithHeight sets the height of the input field.
func (i *TickerPicker) WithHeight(height int) huh.Field {
	i.height = height
	// i.textarea.SetHeight(height - 2)
	return i
}

// WithPosition sets the position of the input field.
func (i *TickerPicker) WithPosition(p huh.FieldPosition) huh.Field {
	i.keymap.PrevSuggestion.SetEnabled(!p.IsFirst())
	i.keymap.NextSuggestion.SetEnabled(!p.IsLast())
	i.keymap.Submit.SetEnabled(p.IsLast())
	return i
}

// GetKey returns the key of the field.
func (i *TickerPicker) GetKey() string {
	return i.key
}

// GetValue returns the value of the field.
func (i *TickerPicker) GetValue() any {
	return *i.value
}

func truncateText(text string, maxWidth int, addEllipsis bool) string {
	t := strings.TrimSpace(text)
	ellipsis := ""
	if addEllipsis && runewidth.StringWidth(t) > maxWidth {
		ellipsis = "…"
		maxWidth -= runewidth.StringWidth(ellipsis) // Reserve space for the ellipsis
		if maxWidth < 0 {
			maxWidth = 0
		}
	}
	truncate := runewidth.Truncate(t, maxWidth, ellipsis)
	return truncate
}

func clamp(v, low, high int) int {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

func cmdize[T any](t T) tea.Cmd {
	return func() tea.Msg {
		return t
	}
}
