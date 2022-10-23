/*
@author: sk
@date: 2022/10/22
*/
package main

const (
	WINDOW_WIDTH  = 116
	WINDOW_HEIGHT = 26
)

const (
	FRAME_TIME = 50
)

type Point int

const (
	POINT_3 Point = iota
	POINT_4
	POINT_5
	POINT_6
	POINT_7
	POINT_8
	POINT_9
	POINT_10
	POINT_J
	POINT_Q
	POINT_K
	POINT_A
	POINT_2
)

type Suit int

const (
	SUIT_CLUB Suit = iota
	SUIT_DIAMOND
	SUIT_HEART
	SUIT_SPADE
)

type Decide int

const (
	DECIDE_IDLE Decide = iota
	DECIDE_YES
	DECIDE_NO
)
