/*
@author: sk
@date: 2022/10/22
*/
package main

import (
	"fmt"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetPoints(cards []*Card) []int {
	if cards == nil {
		return nil
	}
	res := make([]int, 0)
	for i := 0; i < len(cards); i++ {
		res = append(res, int(cards[i].Point))
	}
	return res
}

func GetCount(nums []int) map[int]int {
	res := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		res[nums[i]]++
	}
	return res
}

//=======================image========================

func CreateImage(w, h int) [][]rune {
	res := make([][]rune, h)
	for i := 0; i < h; i++ {
		res[i] = CreateLine(w)
	}
	return res
}

func CreateLine(l int) []rune {
	res := make([]rune, l)
	for i := 0; i < l; i++ {
		res[i] = ' '
	}
	return res
}

func CreateBorderImage(w, h int) [][]rune {
	res := CreateImage(w, h)
	// 绘制边框
	for i := 1; i < w-1; i++ {
		res[0][i] = '-'
		res[h-1][i] = '-'
	}
	for i := 1; i < h-1; i++ {
		res[i][0] = '|'
		res[i][w-1] = '|'
	}
	res[0][0] = '+'
	res[0][w-1] = '+'
	res[h-1][0] = '+'
	res[h-1][w-1] = '+'
	return res
}

var (
	cardImagePool  = make(map[*Card][][]rune)
	EmptyCardImage = CreateBorderImage(11, 6)
)

func LoadCardImage(card *Card) [][]rune {
	if _, ok := cardImagePool[card]; !ok {
		cardImagePool[card] = createCardImage(card)
	}
	return cardImagePool[card]
}

var (
	points = strings.Split("3,4,5,6,7,8,9,10,J,Q,K,A,2", ",")
	suits  = strings.Split("♣,♦,♥,♠", ",")
)

func createCardImage(card *Card) [][]rune {
	res := CreateBorderImage(11, 6)
	// 绘制点数与花色
	DrawString(res, fmt.Sprintf("%s\n%s", points[card.Point], suits[card.Suit]), 2, 1)
	return res
}

//=====================math========================

func Max(num1, num2 int) int {
	if num1 > num2 {
		return num1
	} else {
		return num2
	}
}

func Min(num1, num2 int) int {
	if num1 > num2 {
		return num2
	} else {
		return num1
	}
}

//=======================draw==========================

// DrawString 绘制文本不忽略换行  x,y处也占用
func DrawString(canvas [][]rune, str string, x, y int) {
	lines := strings.Split(str, "\n")
	image := make([][]rune, len(lines))
	for i := 0; i < len(image); i++ {
		image[i] = []rune(lines[i])
	}
	DrawImage(canvas, image, x, y)
}

// DrawImage 可以不等维度(为了复用)
func DrawImage(canvas, image [][]rune, x, y int) {
	height, width := len(canvas), len(canvas[0])
	if x >= width || y >= height {
		return
	}
	minY := Max(0, y)
	maxY := Min(height, y+len(image))
	if maxY < 0 {
		return //出界
	}
	minX := Max(0, x)
	for i := minY; i < maxY; i++ {
		content := image[i-minY]
		maxX := Min(width, x+len(content))
		if maxX >= 0 { // 不出界才行
			if x < 0 { // 判断是否需要裁剪
				content = content[-x:]
			}
			copy(canvas[i][minX:maxX], content)
		}
	}
}
