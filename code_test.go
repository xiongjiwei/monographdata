package monographdata_test

import (
	"math/rand"
	"monographdata"
	"sync"
	"testing"
)

func TestConcurrentUpdate(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				monographdata.Update(rand.Intn(monographdata.N), rand.Intn(monographdata.N))
			}
		}()
	}

	wg.Wait()
}
