package parse

type Cards struct {
	Hash     int64 `gorm:"primary_key;column:hash;index"`
	Level    int   `gorm:"column:level"`
	Score    int   `gorm:"column:score"`
	CardType *CardType
}

type CardType struct {
	pairBitMap [4]int
	colorBitMap [4]int
	colorBitMapLen [4]int
	flushColor int
	cards Cards
	origin string

}

func NewCardType(c string) *CardType {
	var cardType CardType
	cardType.origin = c
	for x := 0; x < len(c); x += 2 {
		color := COLOR_TABLE[c[x+1]]
		card := CARD_TABLE[c[x]]
		i := 0
		for {
			if cardType.pairBitMap[i] & (1 << uint(card)) == 0 {
				cardType.pairBitMap[i] |= 1 << uint(card)
				break
			}
			i++
		}
		cardType.colorBitMap[color] |= 1 << uint(card)
		cardType.colorBitMapLen[color] ++
		if cardType.colorBitMapLen[color] >= 5{
			cardType.flushColor = color + 1
		}
		if card == CARD_TABLE['A']{
			cardType.pairBitMap[0] |= 2
			cardType.colorBitMap[color] |= 2
		}
	}
	cardType.cards.CardType = &cardType
	return &cardType
}

func getOneBitNumber(x int)int{
	countx := 0
	for {
		if x == 0{
			break
		}
		countx ++
		x = x&(x-1)
	}
	return countx
}

func getConnBit(cardMap int)int{
	x := cardMap
	k := 0
	maxBitPos := 0
	for {
		maxBitPos = x
		x = x & (x << 1)
		k++
		if x != 0{
			maxBitPos = x
		}else {
			break
		}
	}
	if k >= 5{
		return maxBitPos
	}
	return 0
}

func (c *CardType) GetScore() *Cards {
	if len(c.origin) == 10{
		c.cards.Score = c.pairBitMap[0] & ^2
		return &c.cards
	}
	switch c.cards.Level {
	case HighCard:
		score := c.pairBitMap[0]
			score = score & ^2
			score &= score - 1
			score &= score - 1
		c.cards.Score = score
	case DoubleOneCard:
		c.cards.Score = c.pairBitMap[1] << 16
		score := c.pairBitMap[0]
		score = score & ^c.pairBitMap[1]
		score = score & ^2
		score &= score - 1
		score &= score - 1
		c.cards.Score += score
	case DoubleTwoCard:
		c.cards.Score = c.pairBitMap[1] << 16
		score := c.pairBitMap[0]
		score = score & ^c.pairBitMap[1]
		score = score & ^2
		score &= score - 1
		score &= score - 1
		c.cards.Score += score
	case ThreeCard:
		c.cards.Score = c.pairBitMap[2] << 16
		score := c.pairBitMap[0] & ^c.pairBitMap[2]
		score = score & ^2
		score &= score - 1
		score &= score - 1
		c.cards.Score += score
	case FlushCard:
		score := c.colorBitMap[c.flushColor-1]
		score = score & ^2
		numFlushCard := c.colorBitMapLen[c.flushColor-1]
		for{
			if numFlushCard <= 5{
				break
			}
			score &= score - 1
			numFlushCard --
		}
		c.cards.Score = score
	case GourdCard:
		c.cards.Score = c.pairBitMap[2] << 16
		score := c.pairBitMap[0] & ^c.pairBitMap[2]
		score = score & ^2
		score &= score - 1
		score &= score - 1
		c.cards.Score += score
	case FourCard:
		c.cards.Score = c.pairBitMap[3] << 16
		score := c.pairBitMap[0] & ^c.pairBitMap[3]
		score = score & ^2
		score &= score - 1
		score &= score - 1
		c.cards.Score += score
	case StraightCard, StraightFlush, RoyalFlush:
	}
	return &c.cards
}

func (c *CardType) GetCard() *Cards {
	//defer  c.GetScore()
	if  c.flushColor != 0{
		maxBitPos := getConnBit(c.colorBitMap[c.flushColor-1])
		if maxBitPos != 0{
			if (^c.colorBitMap[c.flushColor-1] & (0x1F << (15-5))) == 0{
				c.cards.Level = RoyalFlush
			}else {
				c.cards.Level = StraightFlush
			}
			c.cards.Score = maxBitPos
			return &c.cards
		}
	}
	if c.pairBitMap[3] != 0{
		c.cards.Level = FourCard
	}else {
		pairNum := getOneBitNumber(c.pairBitMap[1])
		if c.pairBitMap[2] != 0 {
			if pairNum >= 2 {
				c.cards.Level = GourdCard
			} else if c.flushColor != 0{
				c.cards.Level = FlushCard
			} else if maxBitPos := getConnBit(c.pairBitMap[0]); maxBitPos != 0{
				c.cards.Level = StraightCard
				c.cards.Score = maxBitPos
			}else{
				c.cards.Level = ThreeCard
			}
		}else if c.flushColor != 0{
			c.cards.Level = FlushCard
		} else if maxBitPos := getConnBit(c.pairBitMap[0]); maxBitPos != 0{
			c.cards.Level = StraightCard
			c.cards.Score = maxBitPos
		} else if c.pairBitMap[1] != 0 {
			if pairNum >= 2 {
				c.cards.Level = DoubleTwoCard
			} else {
				c.cards.Level = DoubleOneCard
			}
		} else {
			c.cards.Level = HighCard
		}
	}
	return &c.cards
}
