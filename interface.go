/*
@author: sk
@date: 2022/10/22
*/
package main

type DrawAble interface {
	Draw(canvas [][]rune)
}

type ClickAble interface {
	Click(x, y int) bool
}

type PlayAble interface {
	DrawAble

	GetX() int
	GetY() int
	AddCard(card *Card)
	MarkLandlord()
	PreparePlay()
	PlayCard() ([]*Card, bool)
	IsWin() bool
	IsLandlord() bool
}
