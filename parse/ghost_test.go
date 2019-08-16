package parse


import (
"github.com/bmizerany/assert"
"testing"
)


func TestGhost4(t *testing.T) {
	level := NewCardType("TsXnTc4cTdAsQs").GetCard().Level
	assert.Equal(t, level, FourCard)
}


func TestGhost32(t *testing.T) {
	level := NewCardType("TsXnTc4cQdAsQs").GetCard().Level
	assert.Equal(t, level,GourdCard)
}

func TestGhost2(t *testing.T) {
	level := NewCardType("8sXhAc7sJc6hQd").GetCard().Level
	assert.Equal(t, level, DoubleOneCard)
}

func TestGhost22(t *testing.T) {
	level := NewCardType("TsTc4dXc8dAh3h").GetCard().Level
	assert.Equal(t, level, ThreeCard)
}


func TestGhostFlush(t *testing.T) {
	level := NewCardType("TsXs6d4sTdAs7s").GetCard().Level
	assert.Equal(t, level, FlushCard)
}

func TestGhostStraight(t *testing.T) {
	level := NewCardType("TsJcQdKcTdXh3h").GetCard().Level
	assert.Equal(t, level, StraightCard)
}


func TestGhostSStraight(t *testing.T) {
	level := NewCardType("AdXs3s4sTd5s3h").GetCard().Level
	assert.Equal(t, level, StraightCard)
}


func TestGhostFlushStraight(t *testing.T) {
	level := NewCardType("TsJsQsKsTdXs3h").GetCard().Level
	assert.Equal(t, level, StraightFlush)
}

func TestGhostSFlushStraight(t *testing.T) {
	level := NewCardType("Xs2s3s4sTd5s3h").GetCard().Level
	assert.Equal(t, level, StraightFlush)
}


func TestGhostFinal(t *testing.T) {
	level := NewCardType("7sXn3sJsTs8s3d").GetCard().Level
	assert.Equal(t, level, StraightFlush)
}


func TestGhostFlushCard(t *testing.T) {
	level := NewCardType("Ks2h7sXn3sJsTs").GetCard().Level
	assert.Equal(t, level, FlushCard)
}

//alice :=
//bob := ""
func TestGhostFlushCard2(t *testing.T) {
	level := NewCardType("Ts6s8s2sXn9s9h").GetCard().Level
	assert.Equal(t, level, StraightFlush)
}

//bob := ""
func TestGhost31(t *testing.T) {
	level := NewCardType("TcJd2s4c3hJhXn").GetCard().Level
	assert.Equal(t, level, ThreeCard)
}

func TestGhostScoreTypeFlush_1(t *testing.T) {
	alice := "5dXnKd7d6d6sQc"
	bob := "5c2d5dXnKd7d6d"
	aliceScore := NewCardType(alice).GetCard().CardType.GetScore()
	bobScore := NewCardType(bob).GetCard().CardType.GetScore()
	assert.Equal(t, aliceScore.Score == bobScore.Score, true)
}
