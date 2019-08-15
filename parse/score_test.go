package parse

import (
	"github.com/bmizerany/assert"
	"testing"
)
func TestScoreType1_1(t *testing.T) {
	alice :="TsKdJh6d4sQh2s"
	bob := "7c2hTsKdJh6d4s"

	aliceScore := NewCardType(alice).GetCard().Score
	bobScore := NewCardType(bob).GetCard().Score
	assert.Equal(t, aliceScore > bobScore, true)
}

func TestScoreType21_1(t *testing.T) {
	alice := "5sJc9d3d2s3cAc"
	bob := "5dQs5sJc9d3d2s"
	aliceScore := NewCardType(alice).GetCard()
	bobScore := NewCardType(bob).GetCard()
	assert.Equal(t, aliceScore.Score < bobScore.Score, true)
}

func TestScoreType22_1(t *testing.T) {
	alice := "Qs6s5hQhAdAc6c"
	bob := "AsKsQs6s5hQhAd"
	aliceScore := NewCardType(alice).GetCard()
	bobScore := NewCardType(bob).GetCard()
	assert.Equal(t, aliceScore.Score < bobScore.Score, true)
}

func TestScoreTypeStraight_1(t *testing.T) {
	alice := "QsJsTdKdKh5c9s"
	bob := "9hKcQsJsTdKdKh"
	aliceScore := NewCardType(alice).GetCard()
	bobScore := NewCardType(bob).GetCard()
	assert.Equal(t, aliceScore.Score ==  bobScore.Score, true)
}

func TestScoreTypeFlush_1(t *testing.T) {
	alice := "4h7dKh7h8h6h6s"
	bob := "Qd3h4h7dKh7h8h"
	aliceScore := NewCardType(alice).GetCard()
	bobScore := NewCardType(bob).GetCard()
	assert.Equal(t, aliceScore.Score >  bobScore.Score, true)
}