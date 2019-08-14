package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"math"
	"runtime"
	"time"
)

type Cards struct {
	Hash  int64 `gorm:"primary_key;column:hash;index"`
	Level int   `gorm:"column:level"`
	Score int   `gorm:"column:score"`
}

var cardMap = make(map[int64]*Cards, 20000)

var sdb *gorm.DB

func initFromDB(dbPath string) {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Cards{})
	sdb = db
	var records []Cards
	// Get all records
	db.Find(&records)
	for _, r := range records {
		cardMap[r.Hash] = &r
	}
}

func init() {
	initFromDB("records.sqlite")
}

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

func loadJsonFile(fileName string, rep int) *[] Game {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("path error")
	}
	var record Record
	record.Matches = make([]Game, 0, 10000*rep)
	err = json.Unmarshal(data, &record)
	if err != nil {
		panic("format error")
	}
	tempS := record.Matches
	for x := 1; x < rep; x++ {
		record.Matches = append(record.Matches, tempS...)
	}
	return &record.Matches
}

var finalRecords = make([]*Cards, 0, 20000)

func save(id int64, c *Cards) {
	finalRecords = append(finalRecords, c)
}

func finalSave() {
	tx := sdb.Begin()
	for _, r := range finalRecords {
		tx.FirstOrCreate(r)
	}
	tx.Commit()
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

type BaseComparator interface {
	compare(c1, c2 *Cards) int
}

type MapComparator struct {
}

type CardType struct {
	bitMap [4]int
	cards Cards
	pairLevel int
	origin string
}

func NewCardType(c string) *CardType {
	var cardType CardType
	cardType.origin = c
	for x := 0; x < len(c); x += 2 {
		_ = COLOR_TABLE[c[x+1]]
		card := CARD_TABLE[c[x]]
		i := 0
		for {
			if cardType.bitMap[i] & ^(1 << uint(card)) != 0 {
				cardType.bitMap[1] |= 1 << uint(card)
				break
			}
			cardType.bitMap[i] |= 1 << uint(card)
			cardType.pairLevel ++
			break
		}
	}
	return &cardType
}

func (c *CardType) getCard() *Cards {
	x := c.bitMap[1]
	countx := 0
	for {
		if x == 0{
			break
		}
		countx ++
		x = x&(x-1)
	}
	switch c.pairLevel{
	case 4: c.cards.Level = DoubleOneCard
	case 3:
		if countx == 2{
			c.cards.Level = DoubleTwoCard
		}else {
			c.cards.Level = ThreeCard
		}
	case 2:
		if countx == 2{
			c.cards.Level = GourdCard
		}else {
			c.cards.Level = FourCard
		}
	case 5:
		c.cards.Level = HighCard
	}
	return &c.cards
}

func (m *MapComparator) compare(c1, c2 *Cards) int {
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
	t := loadJsonFile("test_file/result.json", 1)
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

func thread(t *[]Game, comparator *BaseComparator, start int, end int, flag chan int) {
	//var aliceCard Cards
	//var bobCard Cards
	i := 0
	for _, game := range (*t)[start:end] {
		aliceCard := NewCardType(game.Alice)
		bobCard := NewCardType(game.Bob)
		res := (*comparator).compare(aliceCard.getCard(), bobCard.getCard())
		if res != game.Result {
			i++
			//panic(fmt.Sprintf("%d am panic!!!", b))
		}
	}
	fmt.Printf("%d fail\n", i)
	flag <- 1
}

func hashCards(s string) int64 {
	return 0
}
