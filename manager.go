/*
@author: sk
@date: 2022/10/22
*/
package main

import (
	"math/rand"
	"sort"
	"time"
)

var (
	MouseManager = newMouseManager()
	GameManager  = newGameManager()
	CardManager  = newCardManager()
)

//====================mouseManager======================

type mouseManager struct {
	X, Y                    int
	lastPress, currentPress bool
}

//func (m *mouseManager) Draw(canvas [][]rune) {
//	DrawString(canvas, fmt.Sprintf("pos:%d,%d", m.X, m.Y), m.X, m.Y)
//}

func newMouseManager() *mouseManager {
	return &mouseManager{}
}

func (m mouseManager) JustPress() bool {
	return m.currentPress && !m.lastPress
}

func (m *mouseManager) UpdatePress(press bool) {
	m.lastPress = m.currentPress
	m.currentPress = press
}

//=======================gameManager=====================

type gameManager struct {
	player    *Player  //用户玩家
	maxPlayer PlayAble // 最大出牌玩家
	winImage  [][]rune
	players   []PlayAble
	state     int // 0 发牌  1 抢地主    2  出牌     3 结束
	actions   []func()
	index     int
}

func (g *gameManager) Draw(canvas [][]rune) {
	for i := 0; i < len(g.players); i++ {
		g.players[i].Draw(canvas)
	}
	if g.state == 3 && g.winImage != nil {
		DrawImage(canvas, g.winImage, 50, 11)
	}
}

func (g *gameManager) Action() {
	g.actions[g.state]()
}

func (g *gameManager) sendCard() {
	// 移动牌 失败  重新发牌   发牌 失败 进入下一个阶段
	if !CardManager.MoveCard(g.players[g.index]) {
		g.index = (g.index + 1) % len(g.players)
		if !CardManager.SendCard(g.players[g.index]) {
			g.state++
			g.player.PrepareDecideLandlord()
		}
	}
}

func (g *gameManager) playCard() {
	cards, ok := g.players[g.index].PlayCard()
	if ok {
		if cards != nil {
			CardManager.MaxCards = cards
			g.maxPlayer = g.players[g.index]
		} else {
			if g.maxPlayer == g.players[(g.index+1)%len(g.players)] {
				CardManager.MaxCards = nil
			}
		}
		if g.players[g.index].IsWin() {
			g.state++
			text := "FARMERS WIN"
			if g.players[g.index].IsLandlord() {
				text = "LANDLORD WIN" // 都按 12 来
			}
			g.winImage = CreateBorderImage(16, 3)
			DrawString(g.winImage, text, 2, 1)
			return
		}
		g.index = (g.index + 1) % len(g.players)
		g.players[g.index].PreparePlay()
	}
}

func (g *gameManager) gameEnd() {
	// DO NOTHING
}

func (g *gameManager) decideLandlord() {
	decide := g.player.DecideLandlord()
	if decide == DECIDE_YES || decide == DECIDE_NO {
		g.index = 0
		if decide == DECIDE_NO {
			g.index = rand.Intn(2) + 1
		}
		g.players[g.index].MarkLandlord()
		cards := CardManager.GetLandlordCards()
		for i := 0; i < len(cards); i++ {
			g.players[g.index].AddCard(cards[i])
		}
		g.players[g.index].PreparePlay()
		g.state++
	}
}

func newGameManager() *gameManager {
	res := &gameManager{state: 0, index: 2} // 为了先给玩家发牌
	res.player = NewPlayer()
	res.players = []PlayAble{res.player, NewComputer(11, 23), NewComputer(94, 86)}
	res.actions = []func(){res.sendCard, res.decideLandlord, res.playCard, res.gameEnd}
	return res
}

//========================cardManager==========================

type cardManager struct {
	cards    []*Card
	sendCard bool
	MaxCards []*Card
	// 绘制直接使用 emptyCardImage 即可
	moveCardSpeed complex128 // 必须用浮点数  否则移动异常
	moveCardPos   complex128
	moveTimer     int
	// x   58    y  0
}

func (c *cardManager) Draw(canvas [][]rune) {
	if c.sendCard { // 发牌阶段
		DrawImage(canvas, EmptyCardImage, 52, 0)
		x := int(real(c.moveCardPos))
		y := int(imag(c.moveCardPos))
		DrawImage(canvas, EmptyCardImage, x, y)
	} else { //发完牌了显示地主牌
		c.drawCards(canvas, c.cards, 58, 0)
		if c.MaxCards != nil { // 绘制玩家出的牌
			c.drawCards(canvas, c.MaxCards, 58, 9)
		}
	}
}

func (c *cardManager) drawCards(canvas [][]rune, cards []*Card, x int, y int) {
	x -= ((len(cards)-1)*4 + 11) / 2
	for i := 0; i < len(cards); i++ {
		img := LoadCardImage(cards[i])
		DrawImage(canvas, img, x, y)
		x += 4
	}
}

func (c *cardManager) MoveCard(player PlayAble) bool {
	if c.moveTimer <= 0 {
		return false
	}
	c.moveTimer--
	c.moveCardPos += c.moveCardSpeed
	if c.moveTimer == 0 {
		player.AddCard(c.cards[0])
		c.cards = c.cards[1:]
	}
	return true
}

func (c *cardManager) SendCard(player PlayAble) bool {
	if len(c.cards) <= 1 { // 暂时1张地主牌
		c.sendCard = false
		return false
	}
	c.moveTimer = 10
	c.moveCardPos = 52
	c.moveCardSpeed = (complex(float64(player.GetX()), float64(player.GetY())) - c.moveCardPos) / 10
	return true
}

func (c *cardManager) GetLandlordCards() []*Card {
	return c.cards
}

func newCardManager() *cardManager {
	cards := make([]*Card, 0)
	for point := 0; point < 13; point++ {
		for suit := 0; suit < 4; suit++ {
			cards = append(cards, &Card{Point: Point(point), Suit: Suit(suit)})
		}
	}
	rand.Seed(time.Now().Unix())
	sort.Slice(cards, func(i, j int) bool {
		return rand.Float64() < 0.5
	})
	return &cardManager{cards: cards, sendCard: true}
}
