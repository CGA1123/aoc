package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseGame(t *testing.T) {
	t.Parallel()

	game := ParseGame("Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green")

	require.Equal(t, int64(1), game.ID)
	require.Equal(
		t,
		[]Round{
			{Blue: 3, Red: 4},
			{Blue: 6, Red: 1, Green: 2},
			{Green: 2},
		},
		game.Rounds,
	)
}
