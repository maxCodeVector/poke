package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"runtime"
	"sort"
	"time"
)

type Game struct {
	Alice  string
	Bob    string
	Result int
}

type Record struct {
	Matches []Game
}

const (
	HighCard      = 1  // 高牌
	DoubleOneCard = 2  // 一对
	DoubleTwoCard = 3  // 二对
	ThreeCard     = 4  // 三条
	StraightCard  = 6  // 顺子
	FlushCard     = 7  // 同花
	GourdCard     = 8  // 三条加对子（葫芦）
	FourCard      = 9  // 四条
	StraightFlush = 10 // 同花顺
	RoyalFlush    = 11 // 皇家同花顺
)

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
var CARD_BIT = 16
var CARD_A_PART = CARD_TABLE['A'] * int(math.Pow(16, 4))

func loadJsonFile(fileName string) *[] Game {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("path error")
	}
	var record Record
	record.Matches = make([]Game, 0, 10000)
	err = json.Unmarshal(data, &record)
	if err != nil {
		panic("format error")
	}
	return &record.Matches
}

type Cards struct {
	cardMap    map[int]int
	colorList  [4][]int
	finalCards *[]kv
	max        int
	min        int
	isFlush    bool
	isStraight bool
	flushIndex int
	cardType   int
	score      int
}

func (c *Cards) NewCards(cards string) {
	c.max = 0
	c.min = 100
	c.isFlush = false
	c.isStraight = false
	c.score = 0
	c.cardMap = make(map[int]int)
	c.colorList[0] = make([]int, 0, 7)
	c.colorList[1] = make([]int, 0, 7)
	c.colorList[2] = make([]int, 0, 7)
	c.colorList[3] = make([]int, 0, 7)
	c.finalCards = nil

	for i := 0; i < len(cards); i += 2 {
		cardNum := CARD_TABLE[cards[i]]
		cardColor := COLOR_TABLE[cards[i+1]]

		c.max = max2int(c.max, cardNum) // will not usage
		c.min = min2int(c.min, cardNum)
		_, ok := c.cardMap[cardNum]
		if !ok {
			c.cardMap[cardNum] = 1
		} else {
			c.cardMap[cardNum] += 1
		}
		c.colorList[cardColor] = append(c.colorList[cardColor], cardNum)
		if !c.isFlush && len(c.colorList[cardColor]) >= 5 {
			c.isFlush = true
			c.flushIndex = cardColor
		}
	}
}

type kv struct {
	Key   int
	Value int
}

type ByValue []kv

func (a ByValue) Len() int      { return len(a) }
func (a ByValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (ss ByValue) Less(i, j int) bool {
	return ss[i].Value > ss[j].Value || ss[i].Value == ss[j].Value && ss[i].Key > ss[j].Key
}

type ByKey []kv

func (a ByKey) Len() int            { return len(a) }
func (a ByKey) Swap(i, j int)       { a[i], a[j] = a[j], a[i] }
func (ss ByKey) Less(i, j int) bool { return ss[i].Key > ss[j].Key }

func (c *Cards) iniScoreInEqualCase() {
	score := 0
	if c.finalCards == nil {
		var ss = make([]kv, 0, len(c.cardMap))
		for k, v := range c.cardMap {
			ss = append(ss, kv{k, v})
		}

		sort.Sort(ByValue(ss))
		c.finalCards = &ss
	}
	for _, card := range *c.finalCards {
		for x := 0; x < card.Value; x++ {
			score = score*CARD_BIT + card.Key
		}
	}

	c.score += score
}

func compareScoreInEqualCase(c1, c2 *Cards) int {

	final1 := *c1.finalCards
	final2 := *c2.finalCards
	for x := 0; x < len(final1); x++ {
		if final1[x].Key > final2[x].Key {
			return 1
		} else if final1[x].Key < final2[x].Key {
			return 2
		} else if c1.isStraight {
			return 0
		}
	}
	return 0
}

type Comparator5Cards struct {
}

type Comparator7Cards struct {
}

type BaseComparator interface {
	compare(cards1, cards2 *Cards) int
}

func (comp *Comparator7Cards) compare(cards1, cards2 *Cards) int {
	comp.judgeCardType(cards1)
	comp.judgeCardType(cards2)
	if cards1.cardType > cards2.cardType {
		return 1
	} else if cards1.cardType < cards2.cardType {
		return 2
	} else {
		return compareScoreInEqualCase(cards1, cards2)
	}
}

func (comp *Comparator7Cards) judgeCardType(cards *Cards) {
	var flushType int
	var pairType int
	var ss1 *[]kv
	var ss2 *[]kv

	var ss = make([]kv, 0, len(cards.cardMap))
	for k, v := range cards.cardMap {
		ss = append(ss, kv{k, v})
	}
	if cards.isFlush {
		flushType, ss1 = comp.judgeFlushType(cards)
	} else {
		flushType, ss1 = comp.judgeIsPureStraight(&cards.cardMap, &ss, 0)
	}
	pairType, ss2 = comp.judgePairCardType(cards, &ss)
	if flushType > pairType {
		cards.finalCards = ss1
		cards.cardType = flushType
		if flushType != FlushCard {
			cards.isStraight = true
		}
	} else {
		cards.finalCards = ss2
		cards.cardType = pairType
	}
}

func (comp *Comparator7Cards) judgeFlushType(cards *Cards) (int, *[]kv) {
	cardsList := cards.colorList[cards.flushIndex]
	m := make(map[int]int)
	for _, card := range cardsList {
		_, ok := m[card]
		if !ok {
			m[card] = 1
		} else {
			m[card] += 1
		}
	}
	var ss = make([]kv, 0, len(m))
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}
	typeRes, ssp := comp.judgeIsPureStraight(&m, &ss, len(cardsList))
	if typeRes == StraightCard {
		if (*ssp)[0].Key == CARD_TABLE['A'] {
			return RoyalFlush, ssp
		}
		return StraightFlush, ssp
	}
	return FlushCard, ssp
}

func (comp *Comparator7Cards) judgeIsPureStraight(m *map[int]int, ssp *[]kv, cardNum int) (int, *[]kv) {
	ss := *ssp
	sort.Sort(ByKey(ss))
	for x := 4; x < len(ss); x++ {
		if ss[x-4].Key-ss[x].Key == 4 {
			ss = ss[x-4 : x+1]
			return StraightCard, &ss
		}
	}
	if len(ss) >= 5 && isKeysInKeys(&SPECIAL_STAIGHT, m) {
		ss = []kv{
			{5, 1},
			{4, 1},
			{3, 1},
			{2, 1},
			{1, 1},
		}
		return StraightCard, &ss
	}
	// get the five biggest pair cards
	tail := len(ss)
	for x := cardNum; x > 5; x-- {
		ss[tail-1].Value--
		if ss[tail-1].Value == 0 {
			tail--
		}
	}
	ss = ss[0:tail]
	return 0, &ss
}

func (comp *Comparator7Cards) judgePairCardType(cards *Cards, ssp *[]kv) (int, *[]kv) {

	var ss = make([]kv, 0, len(cards.cardMap))
	for k, v := range cards.cardMap {
		ss = append(ss, kv{k, v})
	}
	sort.Sort(ByValue(ss))

	var tempS []kv
	typeRes := HighCard
	switch ss[0].Value {
	case 1:
		ss = ss[0:5]
		return HighCard, &ss
	case 2:
		if ss[1].Value == 2 {
			tempS = ss[2:]
			sort.Sort(ByKey(tempS))
			tempS[0].Value = 1
			ss = append(ss[0:2], tempS[0])
			return DoubleTwoCard, &ss
		} else {
			tempS = ss[1:]
			typeRes = DoubleOneCard
		}
	case 3:
		if ss[1].Value >= 2 {
			ss = ss[0:2]
			ss[1].Value = 2
			return GourdCard, &ss
		} else {
			tempS = ss[1:]
			typeRes = ThreeCard
		}
	case 4:
		tempS = ss[1:]
		sort.Sort(ByKey(tempS))
		tempS[0].Value = 1
		ss = append(ss[0:2], tempS[0])
		return FourCard, &ss
	}
	sort.Sort(ByKey(tempS))
	ss = append(ss[0:2], tempS...)

	// get the five biggest pair cards
	tail := len(ss)
	for x := 7; x > 5; x-- {
		ss[tail-1].Value--
		if ss[tail-1].Value == 0 {
			tail--
		}
	}
	ss = ss[0:tail]
	return typeRes, &ss
}

func (comp *Comparator5Cards) compare(cards1, cards2 *Cards) int {
	comp.judgeCardType(cards1)
	comp.judgeCardType(cards2)
	if cards1.cardType > cards2.cardType {
		return 1
	} else if cards1.cardType < cards2.cardType {
		return 2
	} else {
		cards1.iniScoreInEqualCase()
		cards2.iniScoreInEqualCase()
		if cards1.score > cards2.score {
			return 1
		} else if cards1.score < cards2.score {
			return 2
		} else {
			return 0
		}
	}
}

func (comp *Comparator5Cards) judgeCardType(cards *Cards) {
	if len(cards.cardMap) == 5 {
		comp.judgeStraightType(cards)
	} else if len(cards.cardMap) == 4 {
		cards.cardType = DoubleOneCard
	} else if len(cards.cardMap) == 3 {
		if maxValueOfMap(&cards.cardMap) == 3 {
			cards.cardType = ThreeCard
		} else {
			cards.cardType = DoubleTwoCard
		}
	} else if len(cards.cardMap) == 2 {
		if maxValueOfMap(&cards.cardMap) == 4 {
			cards.cardType = FourCard
		} else {
			cards.cardType = GourdCard
		}
	}
	if cards.isFlush {
		cards.cardType = max2int(cards.cardType, FlushCard)
	}
}

func max2int(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min2int(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxValueOfMap(m *map[int]int) int {
	maxV := 0
	for _, v := range *m {
		maxV = max2int(maxV, v)
	}
	return maxV
}

func (comp *Comparator5Cards) judgeStraightType(cards *Cards) {
	if !comp.baseJudgeStaight(cards) {
		cards.cardType = HighCard
	} else if !cards.isFlush {
		cards.cardType = StraightCard
	} else if cards.max == CARD_TABLE['A'] {
		cards.cardType = RoyalFlush
	} else {
		cards.cardType = StraightFlush
	}
}

func (comp *Comparator5Cards) baseJudgeStaight(cards *Cards) bool {
	if cards.max-cards.min == 4 {
		return true
	}
	if isKeysInKeys(&SPECIAL_STAIGHT, &cards.cardMap) {
		cards.score = (cards.score-CARD_A_PART)*CARD_BIT + 1
		cards.max = 0
		return true
	}
	return false
}

func isKeysInKeys(l *[]int, m *map[int]int) bool {
	for _, e := range *l {
		_, ok := (*m)[e]
		if !ok {
			return false
		}
	}
	return true

}

func main() {
	startTime := time.Now().UnixNano() //纳秒
	t := loadJsonFile("test_file/seven_cards.result.json")
	var comparator BaseComparator
	comparator = new(Comparator7Cards)
	const threadNum = 3
	runtime.GOMAXPROCS(threadNum)

	var flags [threadNum]chan int
	for x := 0; x < threadNum; x++ {
		flags[x] = make(chan int)
		start := x * len(*t) / threadNum
		end := min2int(start+len(*t)/threadNum, len(*t))
		go thread(t, &comparator, start, end, flags[x])
	}
	fmt.Printf("time: %d, cpu: %d, go thread: %d\n", time.Now().UnixNano()-startTime, runtime.NumCPU(), runtime.NumGoroutine())
	for x := 0; x < threadNum; x++ {
		<-flags[x]
	}
	fmt.Print("Are you happy?\n")

}

func thread(t *[]Game, comparator *BaseComparator, start int, end int, flag chan int) {
	var aliceCard Cards
	var bobCard Cards
	for k, game := range (*t)[start:end] {
		aliceCard.NewCards(game.Alice)
		bobCard.NewCards(game.Bob)
		res := (*comparator).compare(&aliceCard, &bobCard)
		if res != game.Result {
			panic(fmt.Sprintf("result %d is not true", k))
		}
	}
	flag <- 1
}
