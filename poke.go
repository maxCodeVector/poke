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
	Matches *[]Game
}

const (
	HighCard       = 1  // 高牌
	Double_OneCard = 2  // 一对
	Double_TwoCard = 3  // 二对
	ThreeCard      = 4  // 三条
	StraightCard   = 6  // 顺子
	FlushCard      = 7  // 同花
	GourdCard      = 8  // 三条加对子（葫芦）
	FourCard       = 9  // 四条
	StraightFlush  = 10 // 同花顺
	RoyalFlush     = 11 // 皇家同花顺
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
	err = json.Unmarshal(data, &record)
	if err != nil {
		return nil
	}
	return record.Matches

}

type Cards struct {
	cardMap  map[int]int
	colorList [4][]int
	finalCards *[]kv
	max      int
	min      int
	isFlush  bool
	flushIndex int
	cardType int
	score    int
}

func (c *Cards) NewCards(cards string) {
	c.max = 0
	c.min = 100
	c.isFlush = false
	c.score = 0
	c.cardMap = make(map[int]int)
	c.colorList[0] = make([]int, 0, 7)
	c.colorList[1] = make([]int, 0, 7)
	c.colorList[2] = make([]int, 0, 7)
	c.colorList[3] = make([]int, 0, 7)

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
		if len(c.colorList[cardColor]) >= 5{
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

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (ss ByValue) Less(i, j int) bool {
	if ss[i].Value > ss[j].Value {
		return true
	} else if ss[i].Value == ss[j].Value && ss[i].Key > ss[j].Key {
		return true
	}
	return false
}

type ByKey []kv

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (ss ByKey) Less(i, j int) bool { return ss[i].Key > ss[j].Key }

func (c *Cards) iniScoreInEqualCase() {
	score := 0
	// if c.finalCards == nil{
	// 	var ss []kv = make([]kv, 0, len(c.cardMap))
	// 	for k, v := range c.cardMap {
	// 		ss = append(ss, kv{k, v})
	// 	}

	// 	sort.Sort(ByValue(ss))
	// 	c.finalCards = &ss
	// }
	// Slice(ss, func(i, j int) bool {
	// 	if ss[i].Value > ss[j].Value {
	// 		return true
	// 	} else if ss[i].Value == ss[j].Value && ss[i].Key > ss[j].Key {
	// 		return true
	// 	}
	// 	return false
	// })
	for _, card := range (*c.finalCards) {
		for x := 0; x < card.Value; x++ {
			score = score*CARD_BIT + card.Key
		}
	}

	c.score += score
}

type Comparator struct {
}

type SuperComparator struct {
}

type BaseComparator interface {
	compare(cards1, cards2 *Cards) int
}

func (comp *SuperComparator) compare(cards1, cards2 *Cards) int {
	comp.judge_cardType(cards1)
	comp.judge_cardType(cards2)
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


func (comp *SuperComparator) judge_cardType(cards *Cards){
	var flushType int
	var strightType int
	var pairType int
	var ss1 *[]kv
	var ss2 *[]kv
	var ss3 *[]kv
	if cards.isFlush {
		flushType, ss1 = comp.judgeFlushType(cards)
	}else {
		strightType, ss2 = comp.judgeIsPureStaight(&cards.cardMap)
		pairType, ss3 = comp.judgePairCardType(cards)
	}
	if flushType > pairType {
		if flushType > strightType{
			cards.finalCards = ss1
			cards.cardType = flushType
		}else{
			cards.finalCards = ss2
			cards.cardType = strightType
		}
	}else if pairType > strightType{
		cards.finalCards = ss3
		cards.cardType = pairType
	}else{
		cards.finalCards = ss2
		cards.cardType = strightType
	}	
}


func (comp *SuperComparator) judgeFlushType(cards *Cards) (int, *[]kv) {
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
	typeRes, ss := comp.judgeIsPureStaight(&m)
	if typeRes == StraightCard{
		if (*ss)[0].Key == CARD_TABLE['A']{
			return RoyalFlush, ss
		}
		return StraightFlush, ss
	}
	return FlushCard, ss
}


func (comp *SuperComparator) judgeIsPureStaight(m *map[int]int) (int, *[]kv) {
	var ss []kv = make([]kv, 0, len(*m))
	for k, v := range *m {
		for x:=0;x<v;x++{
			ss = append(ss, kv{k, 1})
		}
	}
	sort.Sort(ByKey(ss))
	for x:= 4; x < len(ss); x++{
		if (ss[x].Key - ss[x-4].Key) == 5{
			ss = ss[x-4:x+1]
			return StraightCard, &ss
		}
	}
	if isKeysInKeys(&SPECIAL_STAIGHT, m){
		ss = []kv{
			kv{5, 1},
			kv{4, 1},
			kv{3, 1},
			kv{2, 1},
			kv{1, 1},
		}
		return StraightCard, &ss
	}
	ss = ss[:5]
	return HighCard, &ss
}

func (comp *SuperComparator) judgePairCardType(cards *Cards) (int, *[]kv) {

	var ss []kv = make([]kv, 0, len(cards.cardMap))
	for k, v := range cards.cardMap {
		ss = append(ss, kv{k, v})
	}

	sort.Sort(ByValue(ss))

	// get the five biggest pair cards
	tail := len(ss) -1
	if ss[tail].Value > 2{
		ss[tail].Value -= 2
	}else if ss[tail].Value == 2{
		ss = ss[:tail]
	}else{
		tail --
		if ss[tail].Value > 1{
			ss = ss[:tail+1]
			ss[tail].Value --
		}else{
			ss = ss[:tail]
		}
	}

	switch len(ss) {
	case 4:
		return Double_OneCard, &ss
	case 3:
		if ss[0].Value == 3 {
			return ThreeCard, &ss
		} else {
			return Double_TwoCard, &ss
		}
	case 2:
		if ss[0].Value == 4 {
			return FourCard, &ss
		} else {
			return GourdCard, &ss
		}
	}
	return 0, &ss
}




func (comp *Comparator) compare(cards1, cards2 *Cards) int {
	comp.judge_cardType(cards1)
	comp.judge_cardType(cards2)
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

func (comp *Comparator) judge_cardType(cards *Cards) {
	if len(cards.cardMap) == 5 {
		comp.judgeStraightType(cards)
	} else if len(cards.cardMap) == 4 {
		cards.cardType = Double_OneCard
	} else if len(cards.cardMap) == 3 {
		if maxValueOfMap(&cards.cardMap) == 3 {
			cards.cardType = ThreeCard
		} else {
			cards.cardType = Double_TwoCard
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

func (comp *Comparator) judgeStraightType(cards *Cards) {
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

func (comp *Comparator) baseJudgeStaight(cards *Cards) bool {
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
	comparator = new(SuperComparator)
	const threadNum = 1
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
	i:=0
	for _, game := range (*t)[start:end] {
		aliceCard.NewCards(game.Alice)
		bobCard.NewCards(game.Bob)
		res := (*comparator).compare(&aliceCard, &bobCard)
		if res != game.Result {
			i ++
			// panic("result is not true")
		}
	}
	fmt.Printf("incorrect: %d %d\n", i, len(*t))
	flag <- 1
}
