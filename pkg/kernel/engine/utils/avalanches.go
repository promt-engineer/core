package utils

type Avalanche[Symbol comparable] struct {
	Window   [][]Symbol        `json:"window"`
	PayItems []PayItem[Symbol] `json:"pay_items"`
}

type PayItem[Symbol comparable] struct {
	Symbol  Symbol  `json:"symbol"`
	Indexes [][]int `json:"indexes"`
	Award   int64   `json:"award"`
}

func Spin[Symbol comparable](
	ag AwardGetter[Symbol], window AvalancheWindow[Symbol], reels []map[int]Symbol, stops []int,
) (avalanches []Avalanche[Symbol], award int64, err error) {
	iReels := NewIndexedReels(reels)

	for {
		var avalanche Avalanche[Symbol]

		if err := window.Compute(stops, iReels.Copy(), iReels.MaxIndexes(), iReels.DeletedIndexes()); err != nil {
			return nil, 0, err
		}

		avalanche.Window = window.Matrix()

		winItems := window.CheckWin()

		win := false
		for _, winItems := range winItems {

			symbol := winItems.GetSymbol()
			qty := winItems.Count()

			if window.GetScatterSymbol() != nil && symbol == *window.GetScatterSymbol() {
				window.SetScatterQty(qty)

				continue
			}

			if window.GetMultiplierSymbol() != nil && symbol == *window.GetMultiplierSymbol() {
				window.SetMultiplierQty(qty)

				continue
			}

			if currentAward := ag.GetAward(symbol, qty); currentAward > 0 {
				award += currentAward

				avalanche.PayItems = append(
					avalanche.PayItems,
					PayItem[Symbol]{Symbol: symbol, Indexes: winItems.GetIndexes(), Award: currentAward},
				)

				iReels.Delete(window.GetAbsIndexesBySymbol(symbol))

				win = true
			}
		}

		avalanches = append(avalanches, avalanche)

		if !win {
			break
		}
	}

	return
}
