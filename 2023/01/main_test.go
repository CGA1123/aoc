package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PartTwo(t *testing.T) {
	t.Parallel()

	require.Equal(t, int64(29), fixedCalibrationValue("two1nine"))
}
