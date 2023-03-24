package monographdata

import (
	"sync/atomic"
)

const (
	N = 100000
)

var (
	Data       = [N]int{}
	ShadowData = [N]int{}
	lm         *lockManager
)

const (
	w int32 = 2 << 20
)

type Pair struct {
	I int
	J int
}

type lockManager struct {
	locks [N]int32
}

func newLockManager() *lockManager {
	return &lockManager{}
}

func (l *lockManager) needLockObjs(i, j int, sort bool) (objs [4]int, cnt int) {
	objs = [4]int{i, (i + 1) % N, (i + 2) % N, j}
	cnt = 4
	if j == i || j == (i+1)%N || j == (i+2)%N {
		objs[(j+N-i)%N] = objs[3]
		cnt = 3
	}
	if sort {
		for i1 := 0; i1 < cnt; i1++ {
			for i2 := i1 + 1; i2 < cnt; i2++ {
				if objs[i1] > objs[i2] {
					objs[i1] ^= objs[i2]
					objs[i2] ^= objs[i1]
					objs[i1] ^= objs[i2]
				}
			}
		}
	}

	return
}

func (l *lockManager) lock(i, j int) {
	objs, cnt := l.needLockObjs(i, j, true)

	for idx := 0; idx < cnt; idx++ {
		done := false
		for !done {
			tp := l.locks[objs[idx]]
			if tp == w {
				continue
			}
			if objs[idx] != j {
				done = atomic.CompareAndSwapInt32(&l.locks[objs[idx]], tp, tp+1)
			} else {
				done = atomic.CompareAndSwapInt32(&l.locks[objs[idx]], 0, w)
			}
		}
	}
}

func (l *lockManager) unlock(i, j int) {
	objs, cnt := l.needLockObjs(i, j, false)

	for idx := 0; idx < cnt; idx++ {
		if objs[idx] != j {
			atomic.AddInt32(&l.locks[objs[idx]], -1)
		} else {
			atomic.SwapInt32(&l.locks[objs[idx]], 0)
		}
	}
}

func Update(i, j int, ch chan<- Pair) {
	lm.lock(i, j)
	Data[j] = Data[i] + Data[(i+1)%N] + Data[(i+2)%N]
	ch <- Pair{I: i, J: j}
	lm.unlock(i, j)
}

func Reset() {
	lm = newLockManager()
	for i := range Data {
		Data[i] = i
		ShadowData[i] = i
	}
}
