package parse

import (
	"github.com/bmizerany/assert"
	"testing"
)
func TestNewCardType1(t *testing.T) {
	level := NewCardType("Ts7cKd4cJdAh3h").GetCard().Level
	assert.Equal(t, level, HighCard)
}


func TestNewCardType2(t *testing.T) {
	level := NewCardType("8sAhAc7sJc6hQd").GetCard().Level
	assert.Equal(t, level, DoubleOneCard)
}

func TestNewCardType22(t *testing.T) {
	level := NewCardType("TsTc4d4cAdAh3h").GetCard().Level
	assert.Equal(t, level, DoubleTwoCard)
}

func TestNewCardType32(t *testing.T) {
	level := NewCardType("TsTc4d4cTdAh3h").GetCard().Level
	assert.Equal(t, level, GourdCard)
}



func TestNewCardType4(t *testing.T) {
	level := NewCardType("TsTcTd4cTdAh3h").GetCard().Level
	assert.Equal(t, level, FourCard)
}


func TestNewCardTypeFlush(t *testing.T) {
	level := NewCardType("Ts3sTd4sTdAs7s").GetCard().Level
	assert.Equal(t, level, FlushCard)
}

func TestNewCardTypeStraight(t *testing.T) {
	level := NewCardType("TsJcQdKcTdAh3h").GetCard().Level
	assert.Equal(t, level, StraightCard)
}

func TestNewCardTypeFinal(t *testing.T) {
	level := NewCardType("TsJsQsKsTdAs3h").GetCard().Level
	assert.Equal(t, level, RoyalFlush)
}

func TestNewCardTypeFlushStraight(t *testing.T) {
	level := NewCardType("TsJsQsKsTd9s3h").GetCard().Level
	assert.Equal(t, level, StraightFlush)
}

func TestNewCardTypeSFlushStraight(t *testing.T) {
	level := NewCardType("As2s3s4sTd5s3h").GetCard().Level
	assert.Equal(t, level, StraightFlush)
}

func TestNewCardTypeSStraight(t *testing.T) {
	level := NewCardType("Ad2s3s4sTd5s3h").GetCard().Level
	assert.Equal(t, level, StraightCard)
}

func TestNewCardType5(t *testing.T) {
	level := NewCardType("6s5h4c3s2c").GetCard().Level
	assert.Equal(t, level, StraightCard)
	level = NewCardType("As2h3s4c5s").GetCard().Level
	assert.Equal(t, level, StraightCard)
}

