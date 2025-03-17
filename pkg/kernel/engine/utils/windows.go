package utils

type Window[PayItem, Symbol comparable] interface {
	GetSymbol(reelIndex int, payItem PayItem) Symbol
	GetByIndexes(reelIndex, symbolIndex int) (PayItem, Symbol)
	GetHeight() int
	GetWidth() int
}

type MegaWaysWindow[Symbol comparable] interface {
	GetWidth() int
	GetHeight(col int) int
	GetSymbol(colIndex int, rowIndex int) Symbol
}

type AvalancheWindow[Symbol comparable] interface {
	// Compute computes window with check if symbol is not deleted.
	Compute(stops []int, reels []map[int]Symbol, maxIndexes []int, deletedIndexes []map[int]struct{}) error

	// CheckWin returns slice of win combinations including unpaid.
	CheckWin() []Win[Symbol]

	// GetIndexesBySymbol returns indexes of symbol in window.
	GetIndexesBySymbol(symbol Symbol) [][]int

	// GetAbsIndexesBySymbol returns absolute indexes of symbol in window. It needs to delete symbols from reels.
	GetAbsIndexesBySymbol(symbol Symbol) [][]int

	// Matrix returns 2D matrix of window. It needs to return window to the frontend.
	Matrix() [][]Symbol

	GetScatterSymbol() *Symbol
	SetScatterQty(qty int)
	GetScatterQty() int

	GetMultiplierSymbol() *Symbol
	SetMultiplierQty(qty int)
	GetMultiplierQty() int
}

type Win[Symbol comparable] interface {
	GetSymbol() Symbol
	GetIndexes() [][]int
	Count() int
}
