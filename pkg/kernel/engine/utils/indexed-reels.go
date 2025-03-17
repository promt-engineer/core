package utils

type IndexedReels[Symbol comparable] interface {
	Delete(indexes [][]int)
	Copy() []map[int]Symbol
	MaxIndexes() []int
	DeletedIndexes() []map[int]struct{}
}

type indexedReels[Symbol comparable] struct {
	data           []map[int]Symbol
	maxIndexes     []int              // max symbol index for each reel
	deletedIndexes []map[int]struct{} // deleted indexes for each reel
}

func NewIndexedReels[Symbol comparable](reels []map[int]Symbol) IndexedReels[Symbol] {

	// find max indexes for each reel
	maxIndexes := make([]int, len(reels))
	for reelIndex, reel := range reels {
		maxIndexes[reelIndex] = len(reel)
	}

	// create deleted indexes
	deletedIndexes := make([]map[int]struct{}, len(reels))
	for i := 0; i < len(reels); i++ {
		deletedIndexes[i] = make(map[int]struct{})
	}

	return &indexedReels[Symbol]{
		data:           reels,
		maxIndexes:     maxIndexes,
		deletedIndexes: deletedIndexes,
	}
}

func (r *indexedReels[Symbol]) Delete(indexes [][]int) {
	for reelIndex, symbolIndexes := range indexes {
		for _, symbolIndex := range symbolIndexes {
			r.deletedIndexes[reelIndex][symbolIndex] = struct{}{}

			// if deleted symbol index is max index, find new max index
			if symbolIndex == r.maxIndexes[reelIndex] {
				for i := symbolIndex - 1; i > 0; i-- {
					if _, ok := r.data[reelIndex][i]; ok {
						r.maxIndexes[reelIndex] = i
						break
					}
				}
			}
		}
	}
}

func (r *indexedReels[Symbol]) DeletedIndexes() []map[int]struct{} {
	return r.deletedIndexes
}

func (r *indexedReels[Symbol]) Copy() []map[int]Symbol {
	return r.data
}

func (r *indexedReels[Symbol]) MaxIndexes() []int {
	return r.maxIndexes
}
