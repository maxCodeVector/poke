package parse

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
)

type Game struct {
	Alice  string
	Bob    string
	Result int
}

type Record struct {
	Matches []Game
}


func LoadJsonFile(fileName string, rep int) *[] Game {

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
var cardMap = make(map[int64]*Cards, 20000)

var sdb *gorm.DB

func InitFromDB(dbPath string) {
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
	//InitFromDB("records.sqlite")
}
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