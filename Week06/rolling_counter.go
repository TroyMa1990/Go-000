package faststats

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// RollingCounter 是滑动窗口计数器，用于统计一段时间窗口内，对每个时间片上发生的事件进行计数
// RollingCounter 使用一片存储桶，通过滑动窗口跟踪一段时间内的事件计数
type RollingCounter struct {

	//len（bucket）是常量，不可变
	//各个bucket的值是原子的，因此它们不接受互斥
	buckets []AtomicInt64

	//两者都不需要锁定（原子操作）
	rollingSum AtomicInt64
	totalSum   AtomicInt64

	rollingBucket RollingBuckets
}

//NewRollingCounter使用桶宽度和桶数初始化滚动计数器
func NewRollingCounter(bucketWidth time.Duration, numBuckets int, now time.Time) RollingCounter {
	ret := RollingCounter{
		buckets: make([]AtomicInt64, numBuckets),
		rollingBucket: RollingBuckets{
			NumBuckets:  numBuckets,
			BucketWidth: bucketWidth,
			StartTime:   now,
		},
	}
	return ret
}

var _ json.Marshaler = &RollingCounter{}
var _ json.Unmarshaler = &RollingCounter{}
var _ fmt.Stringer = &RollingCounter{}

type jsonCounter struct {
	Buckets       []AtomicInt64
	RollingSum    *AtomicInt64
	TotalSum      *AtomicInt64
	RollingBucket *RollingBuckets
}

//MarshalJSON对计数器进行编码。它是线程安全的。
func (r *RollingCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonCounter{
		Buckets:       r.buckets,
		RollingSum:    &r.rollingSum,
		TotalSum:      &r.totalSum,
		RollingBucket: &r.rollingBucket,
	})
}

//UnmarshalJSON存储以前的JSON编码。注意，这不是线程安全的。
func (r *RollingCounter) UnmarshalJSON(b []byte) error {
	var into jsonCounter
	if err := json.Unmarshal(b, &into); err != nil {
		return err
	}
	r.buckets = into.Buckets
	r.rollingSum = *into.RollingSum
	r.totalSum = *into.TotalSum
	r.rollingBucket = *into.RollingBucket
	return nil
}

// String for debugging
func (r *RollingCounter) String() string {
	return r.StringAt(time.Now())
}

//StringAt在给定时间将计数器转换为字符串。
func (r *RollingCounter) StringAt(now time.Time) string {
	b := r.GetBuckets(now)
	parts := make([]string, 0, len(r.buckets))
	for _, v := range b {
		parts = append(parts, strconv.FormatInt(v, 10))
	}
	return fmt.Sprintf("rolling_sum=%d total_sum=%d parts=(%s)", r.RollingSumAt(now), r.TotalSum(), strings.Join(parts, ","))
}

//Inc向当前bucket添加单个事件
func (r *RollingCounter) Inc(now time.Time) {
	r.totalSum.Add(1)
	if len(r.buckets) == 0 {
		return
	}
	idx := r.rollingBucket.Advance(now, r.clearBucket)
	if idx < 0 {
		return
	}
	r.buckets[idx].Add(1)
	r.rollingSum.Add(1)
}

//RollingSumAt返回滚动时间窗口中的事件总数
func (r *RollingCounter) RollingSumAt(now time.Time) int64 {
	r.rollingBucket.Advance(now, r.clearBucket)
	return r.rollingSum.Get()
}

//RollingSum返回滚动时间窗口中的事件总数（随时间变化）时间到了。(With time time.Now())
func (r *RollingCounter) RollingSum() int64 {
	r.rollingBucket.Advance(time.Now(), r.clearBucket)
	return r.rollingSum.Get()
}

//TotalSum返回所有时间的事件总数
func (r *RollingCounter) TotalSum() int64 {
	return r.totalSum.Get()
}

//getbucket按时间倒序返回bucket的副本
func (r *RollingCounter) GetBuckets(now time.Time) []int64 {
	r.rollingBucket.Advance(now, r.clearBucket)
	startIdx := int(r.rollingBucket.LastAbsIndex.Get() % int64(r.rollingBucket.NumBuckets))
	ret := make([]int64, r.rollingBucket.NumBuckets)
	for i := 0; i < r.rollingBucket.NumBuckets; i++ {
		idx := startIdx - i
		if idx < 0 {
			idx += r.rollingBucket.NumBuckets
		}
		ret[i] = r.buckets[idx].Get()
	}
	return ret
}

func (r *RollingCounter) clearBucket(idx int) {
	toDec := r.buckets[idx].Swap(0)
	r.rollingSum.Add(-toDec)
}

//将计数器重置为所有零值。
func (r *RollingCounter) Reset(now time.Time) {
	r.rollingBucket.Advance(now, r.clearBucket)
	for i := 0; i < r.rollingBucket.NumBuckets; i++ {
		r.clearBucket(i)
	}
}
