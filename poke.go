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
		if c1.Score > c2.Score {
			return 1
		} else if c1.Score < c2.Score {
			return 2
		}
		return 0
	}
}

func main() {
	t := parse.LoadJsonFile("json/result.json", 1)
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
			if game.Result == parse.HighCard{
				//panic(fmt.Sprintf("%d am panic!!!", 2))
			}
			i++
		}
	}
	fmt.Printf("%d fail\n", i)
	flag <- 1
}

func hashCards(s string) int64 {
	return 0
}
