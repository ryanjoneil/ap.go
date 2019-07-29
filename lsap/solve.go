package lsap

import "math"

// solve updates A, U, V, and Z to the optimal LSAP solution.
func (a *LSAP) solve() {
	for i := 0; i < a.Size; i++ {
		if a.f[i] < 0 {
			j := a.path(i)
			if j >= 0 {
				a.increase(i, j)
			}
		}
	}

	a.Z = 0
	for i := 0; i < a.Size; i++ {
		a.Z += a.U[i] + a.V[i]
	}
}

func (a *LSAP) path(i int) int {
	lr := []int{i}       // Vector of labelled rows
	uc := map[int]bool{} // Set of unlabelled columns

	for j := 0; j < a.Size; j++ {
		uc[j] = true
		a.pi[j] = math.MaxInt64
	}

	for {
		r := lr[len(lr)-1]
		if r >= a.Size {
			break
		}

		for j := range uc {
			val := a.A[r][j] - a.U[r] - a.V[j]
			if val < a.pi[j] {
				a.pi[j] = val
				a.c[j] = r
			}
		}

		found := false
		for j := range uc {
			if a.pi[j] == 0 {
				found = true
				break
			}
		}

		if !found {
			// d = min { pi[j] | j in uc }
			first := true
			var d int64
			for j := range uc {
				if first || a.pi[j] < d {
					first = false
					d = a.pi[j]
				}
			}

			for _, h := range lr {
				a.U[h] += d
			}

			for j := 0; j < a.Size; j++ {
				if a.pi[j] == 0 {
					a.V[j] -= d
				} else {
					a.pi[j] -= d
				}
			}
		}

		// j = first column in { k in uc | pi[k] = 0 }
		j := -1
		for k := range uc {
			if a.pi[k] == 0 {
				j = k
				break
			}
		}

		if j >= 0 && a.fBar[j] >= 0 {
			lr = append(lr, a.fBar[j])
			delete(uc, j)
		}

		if j >= 0 && a.fBar[j] < 0 {
			return j
		}
	}

	return -1
}

func (a *LSAP) increase(i, j int) {
	for {
		l := a.c[j]
		a.fBar[j] = l
		k := a.f[l]
		a.f[l] = j
		j = k

		if l == i {
			break
		}
	}
}