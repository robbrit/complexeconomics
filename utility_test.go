package econerra

import (
	"math"

	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
)

import "testing"

func ones() []float64 {
	x := make([]float64, NumGoods)
	for i := range x {
		// Don't use the degenerate case of zero.
		x[i] = 1.0
	}
	return x
}

func TestUtilityGradient(t *testing.T) {
	g := fd.Gradient(nil, utility, ones(), nil)

	// Confirm that dU/dx > 0 for all x.
	for _, i := range UtilityGoodsList() {
		if g[int(i)] <= 0.0 {
			t.Errorf("got dx %s = %.2f, want > 0", i, g[int(i)])
		}
	}
}

func TestUtilityHessian(t *testing.T) {
	h := fd.Hessian(nil, utility, ones(), nil)

	// Hessians for complements should be positive, negative for substitutes, and
	// zero for goods that are independent.
	signMap := map[Good]map[Good]int{
		Meat: {
			Vegetables: -1,
			Beer:       1,
		},
		Vegetables: {
			Meat: -1,
			Beer: 1,
		},
		Beer: {
			Meat:       1,
			Vegetables: 1,
		},
		// No need to put clothing, Hessian is always zero.
		Clothing: {},
	}

	for _, i := range UtilityGoodsList() {
		for _, j := range UtilityGoodsList() {
			got := h.At(int(i), int(j))
			var want int
			if i == j {
				want = -1
			} else {
				want = signMap[i][j]
			}

			if want == 0 {
				if d := math.Abs(got - 0.0); d > 0.001 {
					t.Errorf("got d2x %s/%s = %.2f, want 0", i, j, got)
				}
			} else {
				gotsb := math.Signbit(got)
				wantsb := math.Signbit(float64(want))
				if gotsb != wantsb {
					t.Errorf("got d2x %s/%s sign = %t, want %t", i, j, gotsb, wantsb)
				}
			}
		}
	}

	// Lastly the Hessian should be negative semi-definite. We can check that by
	// verifying that all the eigenvalues are <= 0.
	l := mat.EigenSym{}
	if !l.Factorize(h, false) {
		t.Errorf("could not calculate eigenvalues for Hessian")
	} else {
		for _, li := range l.Values(nil) {
			if li > 0 {
				t.Errorf("got Hessian eigenvalue %2f, want non-positive", li)
			}
		}
	}
}
