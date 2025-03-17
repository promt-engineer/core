package utils

import (
	"reflect"
	"testing"
)

const wild = 7

var scatter = 8

type window [][]int

func (w window) GetWidth() int {
	return len(w)
}

func (w window) GetHeight(col int) int {
	return len(w[col])
}

func (w window) GetSymbol(colIndex int, rowIndex int) int {
	return w[colIndex][rowIndex]
}

func Test_CheckWindow(t *testing.T) {
	type args[Symbol comparable] struct {
		window MegaWaysWindow[int]
	}
	type testCase[Symbol comparable] struct {
		name string
		args args[Symbol]
		want []MegaWayWin[Symbol]
	}
	tests := []testCase[int]{
		{
			name: "1 symbol 1 path",
			args: args[int]{window: window{
				{1},
				{1},
				{1},
			}},
			want: []MegaWayWin[int]{
				{
					Symbol: 1,
					Path:   []int{0, 0, 0},
				},
			},
		},
		{
			name: "1 symbol 2 path",
			args: args[int]{window: window{
				{1},
				{1, 1},
				{1},
			}},
			want: []MegaWayWin[int]{
				{
					Symbol: 1, Path: []int{0, 0, 0},
				},
				{
					Symbol: 1, Path: []int{0, 1, 0},
				},
			},
		},
		{
			name: "2 symbols 4 path",
			args: args[int]{window: window{
				{1, 0},
				{1, 1, 0, 1},
				{1},
			}},
			want: []MegaWayWin[int]{
				{
					Symbol: 1, Path: []int{0, 0, 0},
				},
				{
					Symbol: 1, Path: []int{0, 1, 0},
				},
				{
					Symbol: 1, Path: []int{0, 3, 0},
				},
				{
					Symbol: 0, Path: []int{1, 2},
				},
			},
		},
		{
			name: "3 symbols 6 path",
			args: args[int]{window: window{
				{1, 0, 2},
				{1, 1, 0, 1},
				{1, 0, 0},
			}},
			want: []MegaWayWin[int]{
				{
					Symbol: 1, Path: []int{0, 0, 0},
				},
				{
					Symbol: 1, Path: []int{0, 1, 0},
				},
				{
					Symbol: 1, Path: []int{0, 3, 0},
				},
				{
					Symbol: 0, Path: []int{1, 2, 1},
				},
				{
					Symbol: 0, Path: []int{1, 2, 2},
				},
				{
					Symbol: 2, Path: []int{2},
				},
			},
		},
		{
			name: "2 symbols 3 path with wild",
			args: args[int]{window: window{
				{1, 2},
				{1, wild},
				{1},
			}},
			want: []MegaWayWin[int]{
				{
					Symbol: 1, Path: []int{0, 0, 0},
				},
				{
					Symbol: 1, Path: []int{0, 1, 0},
				},
				{
					Symbol: 2, Path: []int{1, 1},
				},
			},
		},
		{
			name: "4 scatters",
			args: args[int]{window: window{
				{0, 0, 0, scatter},
				{1, 1},
				{2, scatter},
				{3, scatter, scatter},
			}},
			want: []MegaWayWin[int]{
				{
					Symbol: 0, Path: []int{0},
				},
				{
					Symbol: 0, Path: []int{1},
				},
				{
					Symbol: 0, Path: []int{2},
				},
				{
					Symbol: scatter, Path: []int{0, 0, 0, 0}, // path is not correct, just represent a count of scatters
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckWindow(tt.args.window, wild, &scatter)
			for i, w := range got {
				if !reflect.DeepEqual(w.GetSymbol(), tt.want[i].GetSymbol()) || !reflect.DeepEqual(w.GetIndexes(), tt.want[i].GetIndexes()) {
					t.Errorf("checkWindow() = %v, want %v", w.GetIndexes(), tt.want[i].GetIndexes())
				}
			}
		})
	}
}
