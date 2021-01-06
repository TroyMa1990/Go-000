package faststats

import (
	"fmt"
	"time"
)

//RollingBuckets模拟项目桶的时间滚动列表。使用JSON对这个对象进行编码是安全的，以线程安全的方式。
//这是为了不带锁。前提是桶的总尺寸（NumBuckets*BucketWidth）小于Advance执行所需时间。NumBuckets*BucketWidth>=1。
type RollingBuckets struct {
	NumBuckets   int
	StartTime    time.Time
	BucketWidth  time.Duration
	LastAbsIndex AtomicInt64
}

var _ fmt.Stringer = &RollingBuckets{}

func (r *RollingBuckets) String() string {
	return fmt.Sprintf("RollingBucket(num=%d, width=%s)", r.NumBuckets, r.BucketWidth)
}

func (r *RollingBuckets) Advance(now time.Time, clearBucket func(int)) int {
	if r.NumBuckets == 0 {
		return -1
	}
	diff := now.Sub(r.StartTime)
	if diff < 0 {
		return -1
	}
	absIndex := int(diff.Nanoseconds() / r.BucketWidth.Nanoseconds())
	lastAbsVal := int(r.LastAbsIndex.Get())
	indexDiff := absIndex - lastAbsVal
	if indexDiff == 0 {
		return absIndex % r.NumBuckets
	}
	if indexDiff < 0 {
		if indexDiff >= r.NumBuckets {
			return -1
		}
		return absIndex % r.NumBuckets
	}
	for i := 0; i < r.NumBuckets && lastAbsVal < absIndex; i++ {
		if !r.LastAbsIndex.CompareAndSwap(int64(lastAbsVal), int64(lastAbsVal)+1) {
			return r.Advance(now, clearBucket)
		}
		lastAbsVal++
		clearBucket(lastAbsVal % r.NumBuckets)
	}
	r.LastAbsIndex.CompareAndSwap(int64(lastAbsVal), int64(absIndex))
	return r.Advance(now, clearBucket)
}
