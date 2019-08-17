package parse

type GhostFunction func(c *CardType) int

const UNREACHABLE = 10

type Ghost struct {
	GhostTable [11]GhostFunction
}

func GetLastBitPos(num int) int {
	x := num
	x = x & (x - 1)
	return ^x & num
}

func GetHighestOneBit(num int) int {
	i := num
	i |= i >> 1 //右移过程使用零扩展
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	//i |= i >> 16
	return i ^ (i >> 1)
}

func NewGhost() *Ghost {
	ghost := Ghost{}
	ghost.GhostTable[RoyalFlush] = DistToRoyalFlush
	ghost.GhostTable[StraightFlush] = DistToStraightFlush
	ghost.GhostTable[FourCard] = DistToFour
	ghost.GhostTable[GourdCard] = DistToThreeTwo
	ghost.GhostTable[FlushCard] = DistToFlush
	ghost.GhostTable[StraightCard] = DistToStraight
	ghost.GhostTable[ThreeCard] = DistToThreeOne
	ghost.GhostTable[DoubleTwoCard] = DistToDoubleTwoCard
	ghost.GhostTable[DoubleOneCard] = DistToDoubleOneCard
	ghost.GhostTable[HighCard] = DistToHightCard
	return &ghost
}

func DistToRoyalFlush(c *CardType) int {
	return UNREACHABLE
}

func DistToStraightFlush(c *CardType) int {
	res := UNREACHABLE
	var oneNumBitMap int
	var numOfBits int
	for i, colorNum := range c.colorBitMapLen {
		if colorNum >= 4 {
			oneNumBitMap = c.colorBitMap[i]
			numOfBits = oneNumBitMap&1 + c.colorBitMapLen[i]
			goto next
		}
	}
	return res
next:
	for ; numOfBits >= 4; numOfBits-- {
		lastBitPos := GetLastBitPos(oneNumBitMap)
		fiveLastBitPos := (lastBitPos << 1) + lastBitPos
		fiveLastBitPos = (fiveLastBitPos << 3) + fiveLastBitPos + (lastBitPos << 2)
		if getOneBitNumber(fiveLastBitPos&oneNumBitMap) >= 4 {
			c.Cards.Score = lastBitPos << 4
			res = 1
		}
		oneNumBitMap &= oneNumBitMap - 1
	}
	return res
}

func DistToFour(c *CardType) int {
	if c.Cards.Level == ThreeCard || c.Cards.Level == GourdCard {
		highPos := GetHighestOneBit(c.pairBitMap[2])
		c.pairBitMap[3] |= highPos
		return 1
	}
	return UNREACHABLE
}

func DistToThreeTwo(c *CardType) int {
	//ThreeCard is impossible
	if c.Cards.Level == DoubleTwoCard {
		highPos := GetHighestOneBit(c.pairBitMap[1])
		c.pairBitMap[2] |= highPos
		return 1
	}
	return UNREACHABLE
}

func DistToFlush(c *CardType) int {
	for i, colorNum := range c.colorBitMapLen {
		if colorNum >= 4 {
			c.flushColor = i + 1
			c.colorBitMap[i] |= 1 << 14
			c.colorBitMapLen[i] ++
			return 1
		}
	}
	return UNREACHABLE
}

func DistToStraight(c *CardType) int {
	res := UNREACHABLE
	oneNumBitMap := c.pairBitMap[0]
	var numOfBits = oneNumBitMap & 1
	for numOfBits += c.pairBitMapLen[0]; numOfBits >= 4; numOfBits-- {
		lastBitPos := GetLastBitPos(oneNumBitMap)
		fiveLastBitPos := (lastBitPos << 1) + lastBitPos
		fiveLastBitPos = (fiveLastBitPos << 3) + fiveLastBitPos + (lastBitPos << 2)
		if getOneBitNumber(fiveLastBitPos&oneNumBitMap) >= 4 {
			c.Cards.Score = lastBitPos << 4
			res = 1
		}
		oneNumBitMap &= oneNumBitMap - 1
	}
	return res
}

func DistToThreeOne(c *CardType) int {
	if c.Cards.Level == DoubleOneCard {
		//highPos := GetHighestOneBit(c.pairBitMap[1])
		c.pairBitMap[2] |= c.pairBitMap[1]
		return 1
	}
	return UNREACHABLE
}

func DistToDoubleTwoCard(c *CardType) int {
	if c.Cards.Level == DoubleOneCard {
		// impossible
		return 1
	}
	return UNREACHABLE
}

func DistToDoubleOneCard(c *CardType) int {
	if c.Cards.Level == HighCard {
		highPos := GetHighestOneBit(c.pairBitMap[0])
		c.pairBitMap[1] |= highPos
		return 1
	}
	return UNREACHABLE
}

func DistToHightCard(c *CardType) int {
	return UNREACHABLE
}

//HighCard      = 1  // 高牌
//DoubleOneCard = 2  // 一对
//DoubleTwoCard = 3  // 二对
//ThreeCard     = 4  // 三条
//StraightCard  = 5  // 顺子
//FlushCard     = 6  // 同花
//GourdCard     = 7  // 三条加对子（葫芦）
//FourCard      = 8  // 四条
//StraightFlush = 9 // 同花顺
//RoyalFlush    = 10 // 皇家同花顺
