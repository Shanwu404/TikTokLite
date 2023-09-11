// ipLimiter 限制ip访问频率

package utils

import "sync"

var (
	buckets    = make(map[string]*TokenBucket)
	bucketLock sync.Mutex
)

const (
	BUCKET_CAPACITY = 20 // 令牌桶容量
	FILL_RATE       = 10 // 每分钟填充令牌数
)

// IsRateLimited 检查给定的IP地址是否受到速率限制
func IsRateLimited(ip string) bool {
	bucketLock.Lock()
	defer bucketLock.Unlock()

	if _, exists := buckets[ip]; !exists {
		buckets[ip] = NewTokenBucket(BUCKET_CAPACITY, FILL_RATE)
	}

	return !buckets[ip].Take()
}
