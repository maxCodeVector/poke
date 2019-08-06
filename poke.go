package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
)

type Game struct {
	alice  string
	bob    string
	result int
}

type Record struct {
	matches *[]Game
}

const (
	NoneCard       = 0
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

var CARD_TABLE = map[string]int{
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"T": 10,
	"J": 11,
	"Q": 12,
	"K": 13,
	"A": 14,
	"s": 0,
	"h": -1,
	"d": -2,
	"c": -3,
}
var SPECIAL_STAIGHT = []int{2, 3, 4, 5, 14}
var CARD_BIT = 16
var CARD_A_PART = CARD_TABLE["A"] * int(math.Pow(16, 4))

type Cards struct {
	cardMap map[int]int
	max int
	min int
	isFlush bool
	cardType int
	score int
}

func(c Cards) NewCards(cards string)  {
	c.max = 0
	c.min = 100
	c.isFlush = true
	c.score = 0
	preColor := CARD_TABLE[string(cards[1])]
	for i := 0; i < len(cards); i+=2 {
		cardNum := CARD_TABLE[string(cards[i])]
		cardColor := CARD_TABLE[string(cards[i + 1])]

		c.max = max2int(c.max, cardNum)
		c.min =min2int(c.min, cardNum)
		_, ok := c.cardMap[cardNum]
		if !ok {
			c.cardMap[cardNum] = 1
		}else {
			c.cardMap[cardNum] += 1
		}

		if cardColor != preColor{
			c.isFlush = false
		}
	}
}


func (c Cards) iniScoreInEqualCase() {
	score := 0
	type kv struct {
		Key   int
		Value int
	}

	var ss []kv
	for k, v := range c.cardMap {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		if ss[i].Value > ss[j].Value{
			return true
		}else if ss[i].Value == ss[j].Value && ss[i].Key > ss[i].Key{
			return true
		}
		return false
	})
	for _, card := range ss{
		for x:=0;x<card.Value;x++{
			score = score * CARD_BIT + card.Key
		}
	}

	c.score += score
}


func loadFont(fileName string) *[] Game {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("path error")
	}
	var record Record
	err = json.Unmarshal(data, &record)
	if err != nil {
		return nil
	}
	return record.matches

}

type Comparator struct {
}

type BaseComparator interface {
	compare(cards1, cards2 Cards) int
}


func (comp Comparator) compare(cards1, cards2 Cards) int{
	comp.judge_cardType(cards1)
	comp.judge_cardType(cards2)
	if cards1.cardType > cards2.cardType {
		return 1
	}else if cards1.cardType < cards2.cardType {
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

func (comp Comparator)judge_cardType(cards Cards) {
	if len(cards.cardMap) == 5 {
		comp.judgeStraightType(cards)
	}else if len(cards.cardMap) == 4{
		cards.cardType = Double_OneCard
	}else if len(cards.cardMap) == 3{
		if maxValueOfMap(&cards.cardMap) == 3 {
			cards.cardType = ThreeCard
		} else {
			cards.cardType = Double_TwoCard
		}
	}else if len(cards.cardMap) == 2 {
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
	if a > b{
		return a
	}
	return b
}

func min2int(a, b int) int {
	if a < b{
		return a
	}
	return b
}

func maxValueOfMap(m *map[int]int) int {
	maxV := 0
	for _, v := range *m{
		maxV = max2int(maxV, v)
	}
	return maxV
}


func (comp Comparator) judgeStraightType(cards Cards) {
	if comp.baseJudgeStaight(cards) {
		cards.cardType = HighCard
	}else if !cards.isFlush {
		cards.cardType = StraightCard
	}else if int(cards.max) == CARD_TABLE["A"] {
		cards.cardType = RoyalFlush
	}else {
		cards.cardType = StraightFlush
	}
}

func( comp Comparator) baseJudgeStaight(cards Cards)bool{
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
		if !ok{
			return false
		}
	}
	return true

}

func main() {
	t := loadFont("result.json")
	var comparator BaseComparator
	comparator = Comparator{}

	for _, game := range *t {
		var aliceCard Cards
		var bobCard Cards
		aliceCard.NewCards(game.alice)
		bobCard.NewCards(game.bob)
		res := comparator.compare(aliceCard, bobCard)
		if res != game.result{
			os.Exit(-1)
		}

	}

	fmt.Print("Are you happy?\n")

}
