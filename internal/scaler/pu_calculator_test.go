package scaler_test

import (
	"fmt"
	"testing"

	"github.com/nktks/dev-span-pu-scaler/internal/scaler"
	"github.com/stretchr/testify/require"
)

func TestPUCalculator_DesiredPU(t *testing.T) {
	tc := []struct {
		buffer  int
		dbCount int
		want    int
	}{
		{
			buffer:  5,
			dbCount: 0,
			want:    100,
		},
		{
			buffer:  5,
			dbCount: 4,
			want:    100,
		},
		{
			buffer:  5,
			dbCount: 14,
			want:    200,
		},
		{
			buffer:  5,
			dbCount: 15,
			want:    300,
		},
		{
			buffer:  5,
			dbCount: 84,
			want:    900,
		},
		{
			buffer:  5,
			dbCount: 85,
			want:    1000,
		},
		{
			buffer:  5,
			dbCount: 100,
			want:    1000,
		},
		{
			buffer:  5,
			dbCount: 200,
			want:    1000,
		},
		{
			buffer:  3,
			dbCount: 6,
			want:    100,
		},
		{
			buffer:  3,
			dbCount: 7,
			want:    200,
		},
	}

	for _, c := range tc {

		calc := scaler.NewPUCalculator(c.dbCount, c.buffer)
		require.Equal(t, c.want, calc.DesiredPU(), fmt.Sprintf("%#v", c))
	}
}
