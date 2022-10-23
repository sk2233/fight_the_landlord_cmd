/*
@author: sk
@date: 2022/10/22
*/
package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
)

type MainApp struct {
	buffer    *strings.Builder
	canvas    [][]rune
	emptyLine []rune
}

func (m *MainApp) Init() tea.Cmd {
	return frame()
}

func (m *MainApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.checkQuit(msg) {
		return m, tea.Quit
	}
	m.updateMouse(msg)
	GameManager.Action()
	if _, ok := msg.(time.Time); ok {
		return m, frame()
	}
	return m, nil
}

func (m *MainApp) View() string {
	m.clearCanvas()
	CardManager.Draw(m.canvas)
	GameManager.Draw(m.canvas)
	m.fillBuffer()
	return m.buffer.String()
}

func (m *MainApp) checkQuit(msg tea.Msg) bool {
	key, ok := msg.(tea.KeyMsg)
	return ok && key.Runes[0] == 'q'
}

func (m *MainApp) updateMouse(msg tea.Msg) {
	press := MouseManager.currentPress
	if mouse, ok := msg.(tea.MouseMsg); ok {
		MouseManager.X, MouseManager.Y = mouse.X, mouse.Y
		press = mouse.Type == tea.MouseLeft
	}
	// 无论如何也要更新
	MouseManager.UpdatePress(press)
}

func (m *MainApp) clearCanvas() {
	for i := 0; i < WINDOW_HEIGHT; i++ {
		copy(m.canvas[i], m.emptyLine)
	}
}

func (m *MainApp) fillBuffer() {
	m.buffer.Reset()
	for i := 0; i < len(m.canvas); i++ {
		if i != 0 {
			m.buffer.WriteRune('\n')
		}
		m.buffer.WriteString(string(m.canvas[i]))
	}
}

func NewMainApp() *MainApp {
	res := &MainApp{}
	res.buffer = &strings.Builder{}
	res.canvas = CreateImage(WINDOW_WIDTH, WINDOW_HEIGHT)
	res.emptyLine = CreateLine(WINDOW_WIDTH)
	return res
}

//=====================self-utils========================

func frame() tea.Cmd {
	return tea.Every(time.Millisecond*FRAME_TIME, func(time0 time.Time) tea.Msg {
		return time0
	})
}
