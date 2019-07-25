package ap

import "math"

// Solve updates A, U, V, and Z to the optimal AP solution.
func (ap *AP) Solve() {
	if !ap.initialized {
		ap.initialize()
	}

	for i := 0; i < ap.Size; i++ {
		if ap.f[i] < 0 {
			j := ap.path(i)
			if j >= 0 {
				ap.increase(i, j)
			}
		}
	}

	ap.Z = 0
	for i := 0; i < ap.Size; i++ {
		ap.Z += ap.U[i] + ap.V[i]
	}
}

func (ap AP) path(i int) int {
	lr := []int{i}       // Vector of labelled rows
	uc := map[int]bool{} // Set of unlabelled columns

	for j := 0; j < ap.Size; j++ {
		uc[j] = true
		ap.pi[j] = math.MaxInt64
	}

	for {
		r := lr[len(lr)-1]
		if r >= ap.Size {
			break
		}

		for j := range uc {
			val := ap.A[r][j] - ap.U[r] - ap.V[j]
			if val < ap.pi[j] {
				ap.pi[j] = val
				ap.c[j] = r
			}
		}

		found := false
		for j := range uc {
			if ap.pi[j] == 0 {
				found = true
				break
			}
		}

		if !found {
			// d = min { pi[j] | j in uc }
			first := true
			var d int64
			for j := range uc {
				if first || ap.pi[j] < d {
					first = false
					d = ap.pi[j]
				}
			}

			for _, h := range lr {
				ap.U[h] += d
			}

			for j := 0; j < ap.Size; j++ {
				if ap.pi[j] == 0 {
					ap.V[j] -= d
				} else {
					ap.pi[j] -= d
				}
			}
		}

		// j = first column in { k in uc | pi[k] = 0 }
		j := -1
		for k := range uc {
			if ap.pi[k] == 0 {
				j = k
				break
			}
		}

		if j >= 0 && ap.fBar[j] >= 0 {
			lr = append(lr, ap.fBar[j])
			delete(uc, j)
		}

		if j >= 0 && ap.fBar[j] < 0 {
			return j
		}
	}

	return -1
}

func (ap AP) increase(i, j int) {
	for {
		l := ap.c[j]
		ap.fBar[j] = l
		k := ap.f[l]
		ap.f[l] = j
		j = k

		if l == i {
			break
		}
	}
}
