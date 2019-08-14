package parse

import (
	"github.com/bmizerany/assert"
	"testing"
)
func TestScoreType1(t *testing.T) {
	level := NewCardType("Ts7cKd4cJdAh3h").GetCard().Level
	assert.Equal(t, level, HighCard)
}


func TestScoreType2(t *testing.T) {
	level := NewCardType("8sAhAc7sJc6hQd").GetCard().Level
	assert.Equal(t, level, DoubleOneCard)
}

func TestScoreType32(t *testing.T) {
	level := NewCardType("TsTc4d4cTdAh3h").GetCard().Level
	assert.Equal(t, level, GourdCard)
}

func TestScoreType22(t *testing.T) {
	level := NewCardType("TsTc4d4cAdAh3h").GetCard().Level
	assert.Equal(t, level, DoubleTwoCard)
}

func TestScoreType4(t *testing.T) {
	level := NewCardType("TsTcTd4cTdAh3h").GetCard().Level
	assert.Equal(t, level, FourCard)
}


func TestScoreTypeFlush(t *testing.T) {
	level := NewCardType("Ts3sTd4sTdAs7s").GetCard().Level
	assert.Equal(t, level, FlushCard)
}

func TestScoreTypeStraight(t *testing.T) {
	level := NewCardType("TsJcQdKcTdAh3h").GetCard().Level
	assert.Equal(t, level, StraightCard)
}

func TestScoreTypeFinal(t *testing.T) {
	level := NewCardType("TsJsQsKsTdAs3h").GetCard().Level
	assert.Equal(t, level, RoyalFlush)
}

func TestScoreTypeFlushStraight(t *testing.T) {
	level := NewCardType("TsJsQsKsTd9s3h").GetCard().Level
	assert.Equal(t, level, StraightFlush)
}

func TestScoreTypeSFlushStraight(t *testing.T) {
	level := NewCardType("As2s3s4sTd5s3h").GetCard().Level
	assert.Equal(t, level, StraightFlush)
}

func TestScoreTypeSStraight(t *testing.T) {
	level := NewCardType("Ad2s3s4sTd5s3h").GetCard().Level
	assert.Equal(t, level, StraightCard)
}