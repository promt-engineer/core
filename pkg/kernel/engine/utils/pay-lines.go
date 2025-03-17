package utils

type PayLine[PayItem, Symbol comparable] struct {
	PayLineIndex int
	PayLineItems []PayItem
	PaySymbol    Symbol

	Direction string

	Award int64
}

type MegaWayWin[Symbol comparable] struct {
	Symbol Symbol
	Path   []int
}

func (p *MegaWayWin[Symbol]) GetSymbol() Symbol {
	return p.Symbol
}

func (p *MegaWayWin[Symbol]) GetIndexes() [][]int {
	result := make([][]int, len(p.Path))
	for i, index := range p.Path {
		result[i] = []int{index}
	}
	return result
}

func (p *MegaWayWin[Symbol]) Count() int {
	return len(p.Path)
}
