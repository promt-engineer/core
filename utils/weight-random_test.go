package utils

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

var testData = map[int]int{
	2:   500,
	3:   400,
	5:   300,
	8:   200,
	10:  160,
	12:  150,
	15:  140,
	18:  30,
	20:  20,
	25:  10,
	30:  5,
	35:  3,
	50:  2,
	100: 1,
}

func TestChooser_Pick(t *testing.T) {
	client, err := rng.NewMockClient(nil)
	require.NoError(t, err)

	choices := make([]Choice[int, int], len(testData))
	var i int
	for value, weight := range testData {
		choices[i] = NewChoice(value, weight)
		i++
	}

	chooser, err := NewChooser(client, choices...)
	require.NoError(t, err)

	var (
		result uint64
		count  = 1_000_000
		gCount = 8
	)

	results := make(chan uint64, gCount)

	wg := sync.WaitGroup{}
	wg.Add(gCount)

	for i := 0; i < gCount; i++ {
		go func() {
			var r uint64
			for j := 0; j < count; j++ {
				v, err := chooser.Pick()
				require.NoError(t, err)
				r += uint64(v)
			}
			results <- r
			wg.Done()
		}()
	}

	wg.Wait()
	close(results)

Loop:
	for {
		select {
		case v, ok := <-results:
			if !ok {
				break Loop
			}
			result += v
		}
	}

	value := float64(result*1000/uint64(count*gCount)) / 1000
	t.Log("Value: ", value)
	require.InDelta(t, 6.478, value, 0.01)
}

func TestSurjectivity(t *testing.T) {
	testData := map[int]int{
		2: 1,
		3: 1,
		4: 1,
	}

	client, err := rng.NewMockClient(nil)

	choices := make([]Choice[int, int], len(testData))
	var i int
	for value, weight := range testData {
		choices[i] = NewChoice(value, weight)
		i++
	}

	chooser, err := NewChooser(client, choices...)
	require.NoError(t, err)

	dataCounter := make(map[int]int)

	for i := 0; i < 30000; i++ {
		value, err := chooser.Pick()
		require.NoError(t, err)
		dataCounter[value]++
	}

	fmt.Println(dataCounter)

	require.InDelta(t, 10000, dataCounter[2], 100)
	require.InDelta(t, 10000, dataCounter[3], 100)
	require.InDelta(t, 10000, dataCounter[4], 100)
}
