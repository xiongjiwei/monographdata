package monographdata_test

import (
	"math/rand"
	"monographdata"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcurrentUpdate(t *testing.T) {
	test(t)
}

func test(t *testing.T) {
	monographdata.Reset()
	ch := make(chan monographdata.Pair, 1000)
	var chwg, rvwg sync.WaitGroup
	for i := 0; i < 12; i++ {
		chwg.Add(1)
		go func() {
			defer chwg.Done()
			for j := 0; j < 10000; j++ {
				monographdata.Update(rand.Intn(monographdata.N), rand.Intn(monographdata.N), ch)
			}
		}()
	}

	rvwg.Add(1)
	go func() {
		defer rvwg.Done()
		for pair := range ch {
			monographdata.ShadowData[pair.J] = monographdata.ShadowData[pair.I] +
				monographdata.ShadowData[(pair.I+1)%monographdata.N] +
				monographdata.ShadowData[(pair.I+2)%monographdata.N]
		}
		require.Equal(t, monographdata.ShadowData, monographdata.Data)
	}()

	chwg.Wait()
	close(ch)
	rvwg.Wait()
}
