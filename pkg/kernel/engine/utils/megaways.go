package utils

func CheckWindow[Symbol comparable](window MegaWaysWindow[Symbol], wildSymbol Symbol, scatterSymbol *Symbol) []Win[Symbol] {
	var wins []Win[Symbol]

	for i := 0; i < window.GetHeight(0); i++ {
		symbol, cors := checkWindowRecursive(window, window.GetSymbol(0, i), wildSymbol, scatterSymbol, 1, []int{i})

		for _, cor := range uniquePayLines(cors) {
			wins = append(wins, &MegaWayWin[Symbol]{
				Symbol: symbol,
				Path:   cor,
			})
		}
	}

	if scatterSymbol != nil {
		scatterCount := 0

		for i := 0; i < window.GetWidth(); i++ {
			for j := 0; j < window.GetHeight(i); j++ {
				if window.GetSymbol(i, j) == *scatterSymbol {
					scatterCount++
				}
			}
		}

		if scatterCount > 0 {
			wins = append(wins, &MegaWayWin[Symbol]{
				Symbol: *scatterSymbol,
				Path:   make([]int, scatterCount),
			})
		}
	}

	return wins
}

func checkWindowRecursive[Symbol comparable](window MegaWaysWindow[Symbol], symbol, wild Symbol, scatter *Symbol, colIndex int, coords []int) (Symbol, [][]int) {
	var res [][]int

	if colIndex == window.GetWidth() {
		res = append(res, coords)

		return symbol, res
	}

	for rowIndex := 0; rowIndex < window.GetHeight(colIndex); rowIndex++ {
		s := window.GetSymbol(colIndex, rowIndex)

		if scatter != nil && symbol == *scatter {
			continue
		}

		if symbol == s || s == wild {
			c := make([]int, len(coords))
			copy(c, coords)
			c = append(c, rowIndex)

			_, coords := checkWindowRecursive(window, symbol, wild, scatter, colIndex+1, c)

			res = append(res, coords...)

			continue
		}

		res = append(res, coords)
	}

	return symbol, res
}

func uniquePayLines(payLines [][]int) [][]int {
	for i := 0; i < len(payLines); i++ {
		for j := 0; j < len(payLines); j++ {
			if i != j && isFirstSubSlice(payLines[j], payLines[i]) {
				payLines = append(payLines[:j], payLines[j+1:]...)
				j--

				if j < i {
					i--
				}
			}
		}
	}

	return payLines
}

func isFirstSubSlice(cors1, cors2 []int) bool {
	if len(cors1) > len(cors2) {
		return false
	}

	for i := 0; i < len(cors1); i++ {
		if cors1[i] != cors2[i] {
			return false
		}
	}

	return true
}
