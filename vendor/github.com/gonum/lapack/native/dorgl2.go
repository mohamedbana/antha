// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

// Dorgl2 generates an m×n matrix Q with orthonormal rows defined by the
// first m rows product of elementary reflectors as computed by Dgelqf.
//  Q = H(0) * H(2) * ... * H(k-1)
// len(tau) >= k, 0 <= k <= m, 0 <= m <= n, len(work) >= m.
// Dorgl2 will panic if these conditions are not met.
func (impl Implementation) Dorgl2(m, n, k int, a []float64, lda int, tau, work []float64) {
	checkMatrix(m, n, a, lda)
	if len(tau) < k {
		panic(badTau)
	}
	if k > m {
		panic(kGTM)
	}
	if k > m {
		panic(kGTM)
	}
	if m > n {
		panic(nLTM)
	}
	if len(work) < m {
		panic(badWork)
	}
	if m == 0 {
		return
	}
	bi := blas64.Implementation()
	if k < m-1 {
		for i := k; i < m; i++ {
			for j := 0; j < n; j++ {
				a[i*lda+j] = 0
			}
		}
		for j := k; j < m; j++ {
			a[j*lda+j] = 1
		}
	}
	for i := k - 1; i >= 0; i-- {
		if i < n-1 {
			if i < m-1 {
				a[i*lda+i] = 1
				impl.Dlarf(blas.Right, m-i-1, n-i, a[i*lda+i:], 1, tau[i], a[(i+1)*lda+i:], lda, work)
			}
			bi.Dscal(n-i-1, -tau[i], a[i*lda+i+1:], 1)
		}
		a[i*lda+i] = 1 - tau[i]
		for l := 0; l < i; l++ {
			a[i*lda+l] = 0
		}
	}
}
