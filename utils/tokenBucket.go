// TokenBucket算法实现

package utils

import (
	"sync"
	"time"
)

type TokenBucket struct {
	Capacity   int64      // 令牌桶容量
	Tokens     int64      // 令牌数
	FillRate   int64      // 每分钟填充令牌数
	LastRefill time.Time  // 上次填充时间
	mu         sync.Mutex // 互斥锁
}

func NewTokenBucket(capacity int64, fillRate int64) *TokenBucket {
	return &TokenBucket{
		Capacity:   capacity,
		Tokens:     capacity,
		FillRate:   fillRate,
		LastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elasped := now.Sub(tb.LastRefill).Minutes() // 计算距离上次填充过去了多少分钟
	refill := int64(elasped) * tb.FillRate      // 计算需要填充多少令牌

	tb.Tokens += refill
	if tb.Tokens > tb.Capacity {
		tb.Tokens = tb.Capacity
	}

	if tb.Tokens >= 1 {
		tb.Tokens--
		tb.LastRefill = now
		return true
	}

	return false
}
