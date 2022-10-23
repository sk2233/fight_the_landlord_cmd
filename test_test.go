/*
@author: sk
@date: 2022/10/22
*/
package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/muesli/termenv"
)

func TestColorLen(t *testing.T) {
	color := termenv.EnvColorProfile().Color // 获取颜色函数(下面有用)
	// 设置样式函数(放入字符串,放回带有样式信息的字符串)
	keyword := termenv.Style{}.Foreground(color("#00ff00")).Background(color("#ff0000")).Styled
	help := termenv.Style{}.Foreground(color("#0000ff")).Styled
	str1 := keyword("keyword")
	str2 := help("help")
	fmt.Println(len([]rune(str1)), str1)
	fmt.Println(len([]rune(str2)), str2)
}

func TestBuffer(t *testing.T) {
	buf := &strings.Builder{}
	buf.WriteString("sdfsdfsdf")
	buf.Reset()
	fmt.Println(buf.String())
}

func TestDrawStr(t *testing.T) {
	DrawString(nil, "line1\nline2\nline3", 0, 0)
}
