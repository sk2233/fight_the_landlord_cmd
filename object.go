/*
@author: sk
@date: 2022/10/22
*/
package main

import (
	"sort"
	"strconv"
)

//=================================base=================================

type BaseObject struct {
	Image [][]rune
	X, Y  int
}

func NewBaseObject() *BaseObject {
	return &BaseObject{}
}

func (b *BaseObject) Click(x, y int) bool {
	left, top := b.X, b.Y
	right, bottom := left+len(b.Image[0])-1, top+len(b.Image)-1
	if x < left || x > right {
		return false
	}
	if y < top || y > bottom {
		return false
	}
	return true
}

func (b *BaseObject) Draw(canvas [][]rune) {
	DrawImage(canvas, b.Image, b.X, b.Y)
}

//===============================Card=====================================

type HandCard struct {
	*BaseObject
	Card   *Card
	Select bool
}

func (c *HandCard) UnSelect() {
	if c.Select {
		c.Select = false
		c.Y += 2
	}
}

func (c *HandCard) ToggleSelect() {
	c.Select = !c.Select
	if c.Select {
		c.Y -= 2
	} else {
		c.Y += 2
	}
}

func NewHandCard(card *Card) *HandCard {
	res := &HandCard{Card: card}
	res.BaseObject = NewBaseObject()
	res.Image = LoadCardImage(card)
	return res
}

//==========================按钮==============================

type Button struct {
	*BaseObject
}

func NewButton(text string) *Button {
	res := &Button{}
	res.BaseObject = NewBaseObject()
	img := CreateBorderImage(len([]rune(text))+2, 3)
	DrawString(img, text, 1, 1)
	res.Image = img
	return res
}

//===========================BasePlayer==================================

type BasePlayer struct {
	X, Y     int
	Landlord bool
}

func (b *BasePlayer) IsLandlord() bool {
	return b.Landlord
}

func (b *BasePlayer) MarkLandlord() {
	b.Landlord = true
}

func (b *BasePlayer) GetX() int {
	return b.X
}

func (b *BasePlayer) GetY() int {
	return b.Y
}

func NewBasePlayer() *BasePlayer {
	return &BasePlayer{Landlord: false}
}

//===========================Computer==============================

type Computer struct {
	*BasePlayer
	cards   []*Card
	timer   int
	giveUp  bool
	actionX int // 显示放弃的偏移
}

func (b *Computer) IsWin() bool {
	return len(b.cards) == 0
}

func (b *Computer) PreparePlay() {
	b.timer = 20
}

func (b *Computer) PlayCard() ([]*Card, bool) {
	if b.timer > 0 {
		b.timer--
		return nil, false
	}
	points := GetBestCards(GetPoints(CardManager.MaxCards), GetPoints(b.cards))
	if points == nil {
		b.giveUp = true
		return nil, true
	}
	counts := GetCount(points)
	res := make([]*Card, 0)
	last := make([]*Card, 0)
	for i := 0; i < len(b.cards); i++ {
		point := int(b.cards[i].Point)
		if counts[point] > 0 {
			counts[point]--
			res = append(res, b.cards[i])
		} else {
			last = append(last, b.cards[i])
		}
	}
	b.cards = last
	b.giveUp = false
	return res, true
}

func (b *Computer) Draw(canvas [][]rune) {
	// 最多显示 7张 多的 显示 数目
	num := Min(7, len(b.cards))
	x := b.X
	y := b.Y
	for i := 0; i < num; i++ {
		DrawImage(canvas, EmptyCardImage, x, y)
		y += 2
	}
	if len(b.cards) > 7 {
		DrawString(canvas, strconv.Itoa(len(b.cards)), b.X+2, b.Y+1)
	}
	if b.Landlord {
		DrawString(canvas, "\uF316", b.X+8, b.Y+1)
	}
	if b.giveUp {
		DrawString(canvas, "GIVE UP", b.actionX, 9)
	}
}

func (b *Computer) AddCard(card *Card) {
	b.cards = append(b.cards, card)
}

func NewComputer(x, actionX int) *Computer {
	res := &Computer{cards: make([]*Card, 0), giveUp: false}
	res.BasePlayer = NewBasePlayer()
	res.X = x
	res.Y = 0
	res.actionX = actionX
	return res
}

//===========================Player=========================

type Player struct {
	*BasePlayer
	cards []*HandCard
	btns  []*Button
}

func (b *Player) IsWin() bool {
	return len(b.cards) == 0
}

func (b *Player) PreparePlay() {
	// 没有大的牌只能 放弃
	res := GetBestCards(GetPoints(CardManager.MaxCards), b.getPoints())
	noBtn := NewButton("Give Up")
	noBtn.X, noBtn.Y = 54, 15
	b.btns = append(b.btns, noBtn)
	if res != nil {
		noBtn.X = 44
		hitBtn := NewButton("Hint")
		hitBtn.X, hitBtn.Y = 55, 15
		b.btns = append(b.btns, hitBtn)
		yesBtn := NewButton("Play")
		yesBtn.X, yesBtn.Y = 65, 15
		b.btns = append(b.btns, yesBtn)
	}
}

func (b *Player) PlayCard() ([]*Card, bool) {
	if !MouseManager.JustPress() {
		return nil, false
	}
	x, y := MouseManager.X, MouseManager.Y
	if b.btns[0].Click(x, y) { // 放弃
		if CardManager.MaxCards == nil { // 自己出牌不能放弃
			return nil, false
		}
		b.btns = make([]*Button, 0)
		b.playEnd()
		return nil, true
	}
	if len(b.btns) > 1 {
		// 处理卡片点击
		for i := len(b.cards) - 1; i >= 0; i-- {
			if b.cards[i].Click(x, y) {
				b.cards[i].ToggleSelect()
				return nil, false
			}
		} // 提示
		if b.btns[1].Click(x, y) {
			res := GetBestCards(GetPoints(CardManager.MaxCards), b.getPoints())
			if res != nil {
				b.tidyCards()
				counts := GetCount(res)
				for i := 0; i < len(b.cards); i++ {
					point := int(b.cards[i].Card.Point)
					if counts[point] > 0 {
						counts[point]--
						b.cards[i].ToggleSelect()
					}
				}
			}
			return nil, false
		} // 出牌
		if b.btns[2].Click(x, y) {
			src := make([]int, 0)
			last := make([]*HandCard, 0)
			res := make([]*Card, 0)
			for i := 0; i < len(b.cards); i++ {
				if b.cards[i].Select {
					src = append(src, int(b.cards[i].Card.Point))
					res = append(res, b.cards[i].Card)
				} else {
					last = append(last, b.cards[i])
				}
			}
			if IsBigger(GetPoints(CardManager.MaxCards), src) {
				b.cards = last
				b.playEnd()
				return res, true
			}
		}
	}
	return nil, false
}

func (b *Player) Draw(canvas [][]rune) {
	// 绘制手卡
	for i := 0; i < len(b.cards); i++ {
		b.cards[i].Draw(canvas)
	}
	// 绘制按钮
	for i := 0; i < len(b.btns); i++ {
		b.btns[i].Draw(canvas)
	}
	// 绘制地主标记
	if b.Landlord {
		DrawString(canvas, "\uF316", 11+6, b.Y+1)
	}
}

func (b *Player) AddCard(card *Card) {
	b.cards = append(b.cards, NewHandCard(card))
	sort.Slice(b.cards, func(i, j int) bool {
		cardi := b.cards[i].Card
		cardj := b.cards[j].Card
		if cardi.Point != cardj.Point {
			return cardi.Point < cardj.Point
		} else {
			return cardi.Suit < cardj.Suit
		}
	})
	b.tidyCards()
}

func (b *Player) PrepareDecideLandlord() {
	yesBtn := NewButton("Get Landlord")
	yesBtn.X, yesBtn.Y = 49, 15
	noBtn := NewButton("Give Up")
	noBtn.X, noBtn.Y = 65, 15
	b.btns = append(b.btns, yesBtn, noBtn)
}

func (b *Player) DecideLandlord() Decide {
	if !MouseManager.JustPress() {
		return DECIDE_IDLE
	}
	x, y := MouseManager.X, MouseManager.Y
	if b.btns[0].Click(x, y) {
		b.btns = make([]*Button, 0)
		return DECIDE_YES
	}
	if b.btns[1].Click(x, y) {
		b.btns = make([]*Button, 0)
		return DECIDE_NO
	}
	return DECIDE_IDLE
}

func (b *Player) tidyCards() {
	x := 58 - ((len(b.cards)-1)*4+11)/2
	y := b.Y
	for i := 0; i < len(b.cards); i++ {
		b.cards[i].UnSelect()
		b.cards[i].X, b.cards[i].Y = x, y
		x += 4
	}
}

func (b *Player) playEnd() {
	b.tidyCards()
	b.btns = make([]*Button, 0)
}

func (b *Player) getPoints() []int {
	res := make([]int, 0)
	for i := 0; i < len(b.cards); i++ {
		res = append(res, int(b.cards[i].Card.Point))
	}
	return res
}

func NewPlayer() *Player {
	res := &Player{cards: make([]*HandCard, 0)}
	res.BasePlayer = NewBasePlayer()
	res.X = 52
	res.Y = 20
	return res
}
