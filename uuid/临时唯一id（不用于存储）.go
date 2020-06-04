package common

import (
	"fmt"
	"math"
	"sync/atomic"
)

type uuid struct {
	i uint64
}

func (this *uuid)GetUint64() uint64 {
	atomic.CompareAndSwapUint64(&this.i, math.MaxUint64, 0)
	return atomic.AddUint64(&this.i, 1)
}

func (this *uuid)GetString() string {
	atomic.CompareAndSwapUint64(&this.i, math.MaxUint64, 0)
	return fmt.Sprintf("%d", atomic.AddUint64(&this.i, 1))
}

func NewUUID() *uuid {
	return &uuid{}
}
