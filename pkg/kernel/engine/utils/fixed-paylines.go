package utils

const (
	LeftToRightDirection = "left-to-right"
	RightToLeftDirection = "right-to-left"
)

type AwardGetter[Symbol comparable] interface {
	GetAward(symbol Symbol, size int) int64
}

type ScatterTrigger[PayItem comparable] struct {
	PayLineItems [][]PayItem // row -> symbol
}

func (st ScatterTrigger[PayItem]) Count() int {
	res := 0
	for _, row := range st.PayLineItems {
		res += len(row)
	}

	return res
}

// CalcBasePayLines scatter and wild are nullable
func CalcBasePayLines[PayItem, Symbol comparable](payLines [][]PayItem, window Window[PayItem, Symbol], ag AwardGetter[Symbol], scatter, wild *Symbol, direction string) []PayLine[PayItem, Symbol] {
	var foundPayLines []PayLine[PayItem, Symbol]

	for payLineIndex, payLine := range payLines {
		symbol, payLineItems, resDirection := SwitchPayLine(payLine, window, wild, direction)

		if scatter != nil && symbol == *scatter {
			continue
		}

		if award := ag.GetAward(symbol, len(payLineItems)); award > 0 {
			foundPayLines = append(foundPayLines, PayLine[PayItem, Symbol]{
				PayLineIndex: payLineIndex,
				PayLineItems: payLineItems,
				PaySymbol:    symbol,
				Award:        award,

				Direction: resDirection,
			})
		}
	}

	return foundPayLines
}

func CalcScatter[PayItem, Symbol comparable](window Window[PayItem, Symbol], scatter Symbol) ScatterTrigger[PayItem] {
	var scatterPayLine ScatterTrigger[PayItem]
	w, h := window.GetWidth(), window.GetHeight()

	for i := 0; i < w; i++ {
		reelIndexes := []PayItem{}

		for j := 0; j < h; j++ {
			idx, symbol := window.GetByIndexes(i, j)
			if symbol == scatter {
				reelIndexes = append(reelIndexes, idx)
			}
		}

		scatterPayLine.PayLineItems = append(scatterPayLine.PayLineItems, reelIndexes)
	}

	return scatterPayLine
}

func EmptySymbolOrWild[Symbol comparable](symbol Symbol, wild *Symbol) bool {
	var defaultSymbol Symbol

	if wild == nil {
		if symbol == defaultSymbol {
			return true
		}

		return false
	}

	return symbol == defaultSymbol || symbol == *wild
}

func SwitchPayLine[PayItem, Symbol comparable](payLine []PayItem, window Window[PayItem, Symbol], wild *Symbol, direction string) (symbol Symbol, payLineItems []PayItem, resDirection string) {
	switch direction {
	case LeftToRightDirection:
		symbol, payLineItems = CheckPayLine(payLine, window, wild, LeftToRightDirection)
		return symbol, payLineItems, LeftToRightDirection
	case RightToLeftDirection:
		symbol, payLineItems = CheckPayLine(payLine, window, wild, RightToLeftDirection)
		return symbol, payLineItems, RightToLeftDirection
	default:
		symbol, payLineItems = CheckPayLine(payLine, window, wild, LeftToRightDirection)
		return symbol, payLineItems, LeftToRightDirection
	}
}

func CheckPayLine[PayItem, Symbol comparable](payLine []PayItem, window Window[PayItem, Symbol], wild *Symbol, direction string) (symbol Symbol, payLineItems []PayItem) {
	for it := NewPayLineIterator(payLine, direction); it.Valid(); it.Next() {
		reelIndex, payItem := it.Index(), it.Value()

		s := window.GetSymbol(reelIndex, payItem)

		if EmptySymbolOrWild(symbol, wild) {
			payLineItems = append(payLineItems, payItem)
			symbol = s
		} else {
			if wild != nil && s == *wild {
				payLineItems = append(payLineItems, payItem)
			} else if symbol == s {
				payLineItems = append(payLineItems, payItem)
			} else {
				break
			}
		}
	}

	return
}
