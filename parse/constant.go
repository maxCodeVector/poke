package parse

import "math"

const (
	HighCard      = 1  // 高牌
	DoubleOneCard = 2  // 一对
	DoubleTwoCard = 3  // 二对
	ThreeCard     = 4  // 三条
	StraightCard  = 5  // 顺子
	FlushCard     = 6  // 同花
	GourdCard     = 7  // 三条加对子（葫芦）
	FourCard      = 8  // 四条
	StraightFlush = 9 // 同花顺
	RoyalFlush    = 10 // 皇家同花顺
	CARD_BIT      = 16
)

var CARD_A_PART = 14 * int(math.Pow(16, 4))
var COLOR_TABLE = map[byte]int{
	's': 0,
	'h': 1,
	'd': 2,
	'c': 3,
}
var CARD_TABLE = map[byte]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}
var SPECIAL_STAIGHT = []int{2, 3, 4, 5, 14}