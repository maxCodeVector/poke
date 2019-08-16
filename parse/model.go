package parse

type Cards struct {
	Hash     int64 `gorm:"primary_key;column:hash;index"`
	Level    int   `gorm:"column:level"`
	Score    int   `gorm:"column:score"`
	CardType *CardType
}

type CardType struct {
	pairBitMap     [4]int
	colorBitMap    [4]int
	colorBitMapLen [4]int
	flushColor     int
	Cards          Cards
	origin         string
	cheatNum       int
}

func NewCardType(c string) *CardType {
	var cardType CardType
	cardType.origin = c
	for x := 0; x < len(c); x += 2 {
		if c[x] == 'X' {
			cardType.cheatNum++
			continue
		}
		color := COLOR_TABLE[c[x+1]]
		cardMinusOne := CARD_TABLE[c[x]] - 1
		i := 0
		for {
			if cardType.pairBitMap[i]&(1<<uint(cardMinusOne)) == 0 {
				cardType.pairBitMap[i] |= 1 << uint(cardMinusOne)
				break
			}
			i++
		}
		cardType.colorBitMap[color] |= 1 << uint(cardMinusOne)
		cardType.colorBitMapLen[color] ++
		if cardType.colorBitMapLen[color] >= 5 {
			cardType.flushColor = color + 1
		}
		if cardMinusOne == CARD_TABLE['A']-1 {
			cardType.pairBitMap[0] |= 1
			cardType.colorBitMap[color] |= 1
		}
	}
	cardType.Cards.CardType = &cardType
	return &cardType
}

func getOneBitNumber(x int) int {
	countx := 0
	for {
		if x == 0 {
			break
		}
		countx++
		x = x & (x - 1)
	}
	return countx
}

func getConnBit(cardMap int) int {
	x := cardMap
	k := 0
	maxBitPos := 0
	for {
		maxBitPos = x
		x = x & (x << 1)
		k++
		if x != 0 {
			maxBitPos = x
		} else {
			break
		}
	}
	if k >= 5 {
		return maxBitPos
	}
	return 0
}

func (c *CardType) GetScore() *Cards {
	if len(c.origin) == 10 {
		return c.getScoreOf5Cards()
	}
	return c.getScoreOf7Cards()
}

func (c *CardType) getScoreOf5Cards() *Cards {
	c.pairBitMap[0] &= ^1
	switch c.Cards.Level {
	case HighCard:
		c.Cards.Score = c.pairBitMap[0] & ^1
	case DoubleOneCard, DoubleTwoCard:
		c.Cards.Score = (c.pairBitMap[1] << 16) + c.pairBitMap[0] & ^c.pairBitMap[1]
	case ThreeCard:
		c.Cards.Score = (c.pairBitMap[2] << 16) + c.pairBitMap[0] & ^c.pairBitMap[2]
	case FlushCard:
		c.Cards.Score = c.colorBitMap[c.flushColor-1] & ^1
	case GourdCard:
		threeNum := c.pairBitMap[2]
		// got 3 Cards
		twoNum := c.pairBitMap[1] & ^threeNum
		c.Cards.Score = (threeNum << 16) + twoNum
	case FourCard:
		fourNum := c.pairBitMap[3]
		oneNum := c.pairBitMap[0] & ^fourNum
		c.Cards.Score += (fourNum << 16) + oneNum
	case StraightCard, StraightFlush, RoyalFlush:
	}
	return &c.Cards
}

func (c *CardType) getScoreOf7Cards() *Cards {
	switch c.Cards.Level {
	case HighCard:
		score := c.pairBitMap[0]
		score = score & ^1
		score &= score - 1
		score &= score - 1
		c.Cards.Score = score
	case DoubleOneCard:
		c.Cards.Score = c.pairBitMap[1] << 16
		score := c.pairBitMap[0]
		score = score & ^c.pairBitMap[1]
		score = score & ^1
		score &= score - 1
		score &= score - 1
		c.Cards.Score += score
	case DoubleTwoCard:
		twoNum := c.pairBitMap[1]
		if getOneBitNumber(twoNum) == 3 {
			twoNum &= twoNum - 1
		}
		// got 4 Cards
		oneNum := c.pairBitMap[0] & ^twoNum
		for {
			if getOneBitNumber(oneNum) <= 1 {
				break
			}
			oneNum &= oneNum - 1
		}
		c.Cards.Score = (twoNum << 16) + oneNum
	case ThreeCard:
		threeNum := c.pairBitMap[2]
		oneNum := c.pairBitMap[0] & ^threeNum
		for {
			if getOneBitNumber(oneNum) <= 2 {
				break
			}
			oneNum &= oneNum - 1
		}
		c.Cards.Score = (threeNum << 16) + oneNum
	case FlushCard:
		score := c.colorBitMap[c.flushColor-1]
		score = score & ^1
		numFlushCard := c.colorBitMapLen[c.flushColor-1]
		for {
			if numFlushCard <= 5 {
				break
			}
			score &= score - 1
			numFlushCard--
		}
		c.Cards.Score = score
	case GourdCard:
		threeNum := c.pairBitMap[2]
		if getOneBitNumber(threeNum) == 2 {
			threeNum &= threeNum - 1
		}
		// got 3 Cards
		twoNum := c.pairBitMap[1] & ^threeNum
		if getOneBitNumber(twoNum) == 2 {
			twoNum &= twoNum - 1
		}
		c.Cards.Score = (threeNum << 16) + twoNum
	case FourCard:
		fourNum := c.pairBitMap[3]
		oneNum := c.pairBitMap[0] & ^fourNum
		for {
			if getOneBitNumber(oneNum) <= 1 {
				break
			}
			oneNum &= oneNum - 1
		}
		c.Cards.Score = (fourNum << 16) + oneNum
	case StraightCard, StraightFlush, RoyalFlush:
	}
	return &c.Cards
}

func (c *CardType) GetGhostCard() *Cards {
	ghost := NewGhost()
	var level int
	for x := len(ghost.GhostTable) - 1; x > 0; x-- {
		if ghost.GhostTable[x] == nil {
			continue
		}
		if c.Cards.Level > x{
			break
		}
		if ghost.GhostTable[x](c) <= c.cheatNum { // todo , need to consider equal case
			level = x
			break
		}
	}
	if level != 0 {
		c.Cards.Level = level
	}
	return &c.Cards
}

func (c *CardType) GetCard() *Cards {
	//defer  c.GetScore()
	if c.flushColor != 0 {
		maxBitPos := getConnBit(c.colorBitMap[c.flushColor-1]) // todo change the minBitPos to adapt ghost
		if maxBitPos != 0 {
			if (^c.colorBitMap[c.flushColor-1] & (0x1F << (14 - 5))) == 0 {
				c.Cards.Level = RoyalFlush
			} else {
				c.Cards.Level = StraightFlush
			}
			c.Cards.Score = maxBitPos
			return &c.Cards
		}
	}
	if c.pairBitMap[3] != 0 {
		c.Cards.Level = FourCard
	} else {
		pairNum := getOneBitNumber(c.pairBitMap[1])
		if c.pairBitMap[2] != 0 {
			if pairNum >= 2 {
				c.Cards.Level = GourdCard
			} else if c.flushColor != 0 {
				c.Cards.Level = FlushCard
			} else if maxBitPos := getConnBit(c.pairBitMap[0]); maxBitPos != 0 {
				c.Cards.Level = StraightCard
				c.Cards.Score = maxBitPos
			} else {
				c.Cards.Level = ThreeCard
			}
		} else if c.flushColor != 0 {
			c.Cards.Level = FlushCard
		} else if maxBitPos := getConnBit(c.pairBitMap[0]); maxBitPos != 0 {
			c.Cards.Level = StraightCard
			c.Cards.Score = maxBitPos
		} else if c.pairBitMap[1] != 0 {
			if pairNum >= 2 {
				c.Cards.Level = DoubleTwoCard
			} else {
				c.Cards.Level = DoubleOneCard
			}
		} else {
			c.Cards.Level = HighCard
		}
	}
	if c.cheatNum != 0 {
		return c.GetGhostCard()
	}
	return &c.Cards
}
