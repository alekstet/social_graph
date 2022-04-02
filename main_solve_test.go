package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolve(t *testing.T) {
	journal_soc := [][3]int{{1, 2, 3}, {1, 4, 1}, {2, 3, 1}}
	max_node := 4
	min := 1
	max := 3
	avg := float32(2)
	matrix := [][]int{{0, 3, 0, 1}, {3, 0, 1, 0}, {0, 1, 0, 0}, {1, 0, 0, 0}}
	info := Info{max, min, avg}
	expected := Resp{matrix, info}

	t.Run("solve data", func(t *testing.T) {
		actual := Solve(max_node, min, max, journal_soc)
		assert.Equal(t, expected, actual)
	})
}
