/*
@author: sk
@date: 2022/10/22
*/
package main

import tea "github.com/charmbracelet/bubbletea"

func main() {
	//file, err := tea.LogToFile("out.log", "DEBUG")
	//HandleErr(err)
	//defer file.Close()
	err := tea.NewProgram(NewMainApp(), tea.WithAltScreen(), tea.WithMouseAllMotion()).Start()
	HandleErr(err)
}
