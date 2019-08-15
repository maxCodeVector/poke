package main

import (
	"fmt"
	"poke/parse"
	"runtime"
	"time"
)

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

type BaseComparator interface {
	compare(c1, c2 *parse.Cards) int
}

type MapComparator struct {
}


func (m *MapComparator) compare(c1, c2 *parse.Cards) int {
	if c1.Level > c2.Level {
		return 1
	} else if c1.Level < c2.Level {
		return 2
	} else {
		return compareInEqualLevel(c1, c2)
	}
}

func compareInEqualLevel(c1 *parse.Cards, c2 *parse.Cards) int {
	c1.CardType.GetScore()
	c2.CardType.GetScore()
	if c1.Score > c2.Score {
		return 1
	} else if c1.Score < c2.Score {
		return 2
	}
	return 0
}

func main() {
	t := parse.LoadJsonFile("json/seven_cards_with_ghost.result.json", 1)
	startTime := time.Now().UnixNano() //纳秒
	var comparator BaseComparator
	comparator = new(MapComparator)
	const threadNum = 1
	runtime.GOMAXPROCS(threadNum)

	var flags [threadNum]chan int
	for x := 0; x < threadNum; x++ {
		flags[x] = make(chan int)
		start := x * len(*t) / threadNum
		end := min2int(start+len(*t)/threadNum, len(*t))
		go thread(t, &comparator, start, end, flags[x])
	}
	for x := 0; x < threadNum; x++ {
		<-flags[x]
	}

	fmt.Printf("cards %d, go thread: %d\n", len(*t), (time.Now().UnixNano()-startTime)/1000000)
	fmt.Printf("Are you happy?\n")

}

func thread(t *[]parse.Game, comparator *BaseComparator, start int, end int, flag chan int) {
	//var aliceCard Cards
	//var bobCard Cards
	i := 0
	for _, game := range (*t)[start:end] {
		aliceCard := parse.NewCardType(game.Alice)
		bobCard := parse.NewCardType(game.Bob)
		res := (*comparator).compare(aliceCard.GetCard(), bobCard.GetCard())
		if res != game.Result{
			i++
			goto test
			if aliceCard.Cards.Level != bobCard.Cards.Level{
				panic("level has wrong")
			}
			if aliceCard.GetCard().Level == parse.HighCard{
				//fmt.Printf("alice := \"%s\"\n", game.Alice)
				//fmt.Printf("bob := \"%s\"\n", game.Bob)
				//panic("high card has wrong")
			}
			if aliceCard.GetCard().Level == parse.DoubleOneCard{
				fmt.Printf("alice := \"%s\"\n", game.Alice)
				fmt.Printf("bob := \"%s\"\n", game.Bob)
				panic("2 card has wrong")
			}
			if aliceCard.GetCard().Level == parse.DoubleTwoCard{
				fmt.Printf("alice := \"%s\"\n", game.Alice)
				fmt.Printf("bob := \"%s\"\n", game.Bob)
				fmt.Printf("result: %d\n", game.Result)
				panic("22 card has wrong")
			}
			if aliceCard.GetCard().Level == parse.ThreeCard{
				fmt.Printf("alice := \"%s\"\n", game.Alice)
				fmt.Printf("bob := \"%s\"\n", game.Bob)
				fmt.Printf("result: %d\n", game.Result)
				panic("3 card has wrong")
			}
			if aliceCard.GetCard().Level == parse.StraightCard{
				fmt.Printf("alice := \"%s\"\n", game.Alice)
				fmt.Printf("bob := \"%s\"\n", game.Bob)
				fmt.Printf("result: %d\n", game.Result)
				//panic("straight card has wrong")
			}
			if aliceCard.GetCard().Level == parse.FlushCard{
				fmt.Printf("alice := \"%s\"\n", game.Alice)
				fmt.Printf("bob := \"%s\"\n", game.Bob)
				fmt.Printf("result: %d\n", game.Result)
				panic("flush card has wrong")
			}
			if aliceCard.GetCard().Level == parse.GourdCard{
				fmt.Printf("alice := \"%s\"\n", game.Alice)
				fmt.Printf("bob := \"%s\"\n", game.Bob)
				fmt.Printf("result: %d\n", game.Result)
				panic("32 card has wrong")
			}
			if aliceCard.GetCard().Level == parse.FourCard{
				fmt.Printf("alice := \"%s\"\n", game.Alice)
				fmt.Printf("bob := \"%s\"\n", game.Bob)
				fmt.Printf("result: %d\n", game.Result)
				panic("four card has wrong")
			}
			if aliceCard.GetCard().Level == parse.StraightFlush{
				fmt.Printf("alice := \"%s\"\n", game.Alice)
				fmt.Printf("bob := \"%s\"\n", game.Bob)
				fmt.Printf("result: %d\n", game.Result)
				panic("straightflush card has wrong")
			}
			if aliceCard.GetCard().Level == parse.RoyalFlush{
				fmt.Printf("alice := \"%s\"\n", game.Alice)
				fmt.Printf("bob := \"%s\"\n", game.Bob)
				fmt.Printf("result: %d\n", game.Result)
				panic("Rflush card has wrong")
			}
			test: i = i
		}
	}
	fmt.Printf("%d fail\n", i)
	flag <- 1
}

