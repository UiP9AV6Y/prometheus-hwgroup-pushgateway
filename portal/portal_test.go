package portal

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/common"
)

func TestCalculateValueLog(t *testing.T) {
	testCases := map[string]struct {
		haveSequence  common.IntSequence
		haveExponent  int
		haveThreshold float64
		haveCurrent   float64
		wantValue     float64
		wantChanges   uint64
	}{
		"empty": {
			haveSequence: common.IntSequence([]int{}),
		},
		"one/ones": {
			haveSequence: common.IntSequence([]int{1, 1, 1, 1}),
			haveCurrent:  1.0,
			wantValue:    1.0,
		},
		"zero/ones": {
			haveSequence: common.IntSequence([]int{1, 1, 1, 1}),
			wantChanges:  1,
			wantValue:    1.0,
		},
		"flipping": {
			haveSequence: common.IntSequence([]int{1, 0, 1, 0}),
			wantChanges:  4,
		},
		"descend": {
			haveSequence: common.IntSequence([]int{-1, -2, -3, -4}),
			wantChanges:  4,
			wantValue:    -4.0,
		},
		"oscillating": {
			haveSequence: common.IntSequence([]int{-1, 1, -1, 1}),
			wantChanges:  4,
			wantValue:    1.0,
		},
		"threshold": {
			haveSequence:  common.IntSequence([]int{1, 2, 3, 4}),
			haveThreshold: 2.0,
			wantChanges:   1,
			wantValue:     3.0,
		},
		"big step": {
			haveSequence:  common.IntSequence([]int{10, 20, 24, 28}),
			haveCurrent:   8.0,
			haveThreshold: 5.0,
			wantChanges:   2,
			wantValue:     28.0,
		},
		"temperature": {
			haveSequence:  common.IntSequence([]int{100, 106, 113, 119}),
			haveExponent:  -1,
			haveCurrent:   8.0,
			haveThreshold: 2.0,
			wantChanges:   1,
			wantValue:     10.6,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotValue, gotChanges := calculateValueLog(tc.haveSequence, tc.haveExponent, tc.haveThreshold, tc.haveCurrent)

			assert.Assert(t, gotValue == tc.wantValue, "gotValue=%f; wantValue=%f", gotValue, tc.wantValue)
			assert.Assert(t, gotChanges == tc.wantChanges, "gotChanges=%d; wantChanges=%d", gotChanges, tc.wantChanges)
		})
	}
}
