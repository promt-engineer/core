package utils

import "go.uber.org/zap"

type PayLineIterator[PayItem comparable] struct {
	PayLine   []PayItem
	direction string

	i int
}

func NewPayLineIterator[PayItem comparable](payLine []PayItem, direction string) PayLineIterator[PayItem] {
	it := PayLineIterator[PayItem]{PayLine: payLine}

	switch direction {
	case LeftToRightDirection:
		it.direction = LeftToRightDirection
	case RightToLeftDirection:
		it.direction = RightToLeftDirection
	default:
		{
			zap.S().Warn("wrong direction")
			it.direction = LeftToRightDirection
		}
	}

	it.Init()

	return it
}

func (it *PayLineIterator[PayItem]) Init() *PayLineIterator[PayItem] {
	switch it.direction {
	case LeftToRightDirection:
		it.i = 0
	case RightToLeftDirection:
		it.i = len(it.PayLine) - 1
	}

	return it
}

func (it *PayLineIterator[PayItem]) Next() {
	switch it.direction {
	case LeftToRightDirection:
		it.i++
	case RightToLeftDirection:
		it.i--
	}
}

func (it *PayLineIterator[PayItem]) Valid() bool {
	switch it.direction {
	case LeftToRightDirection:
		return it.i < len(it.PayLine)
	case RightToLeftDirection:
		return it.i >= 0
	}

	return false
}

func (it *PayLineIterator[PayItem]) Value() PayItem {
	return it.PayLine[it.i]
}

func (it *PayLineIterator[PayItem]) Index() int {
	return it.i
}
