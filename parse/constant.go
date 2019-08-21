package parse

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
)


var TAG_TABLE [128]int

func init()  {
	TAG_TABLE['s'] = 0
	TAG_TABLE['h'] = 1
	TAG_TABLE['d'] = 2
	TAG_TABLE['c'] = 3
	TAG_TABLE['n'] = -1

	TAG_TABLE['2'] = 2
	TAG_TABLE['3'] = 3
	TAG_TABLE['4'] = 4
	TAG_TABLE['5'] = 5
	TAG_TABLE['6'] = 6
	TAG_TABLE['7'] = 7
	TAG_TABLE['8'] = 8
	TAG_TABLE['9'] = 9
	TAG_TABLE['T'] = 10
	TAG_TABLE['J'] = 11
	TAG_TABLE['Q'] = 12
	TAG_TABLE['K'] = 13
	TAG_TABLE['A'] = 14
	TAG_TABLE['X'] = 15
}