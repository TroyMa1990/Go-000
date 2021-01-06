package faststats

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
	"time"
)

// AtomicBoolean是一个helper结构，用于模拟布尔函数上的原子操作
type AtomicBoolean struct{ flag uint32 }

// 获取当前原子值
func (a *AtomicBoolean) Get() bool {
	return atomic.LoadUint32(&a.flag) == 1
}

// 设置原子值
func (a *AtomicBoolean) Set(value bool) {
	if value {
		atomic.StoreUint32(&a.flag, 1)
	} else {
		atomic.StoreUint32(&a.flag, 0)
	}
}

func (a *AtomicBoolean) String() string {
	return strconv.FormatBool(a.Get())
}

var _ json.Marshaler = &AtomicBoolean{}
var _ json.Unmarshaler = &AtomicBoolean{}

// MarshalJSON以线程安全的方式将该值编码为json bool
func (a *AtomicBoolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Get())
}

// UnmarshalJSON以线程安全的方式将该值解码为json bool
func (a *AtomicBoolean) UnmarshalJSON(b []byte) error {
	var into bool
	if err := json.Unmarshal(b, &into); err != nil {
		return err
	}
	a.Set(into)
	return nil
}

//AtomicInt64是一个助手结构，用于模拟int64上的原子操作
//在不使用原子函数的情况下轻松地进行加减运算。
type AtomicInt64 struct{ val int64 }

var _ json.Marshaler = &AtomicInt64{}
var _ json.Unmarshaler = &AtomicInt64{}

//MarshalJSON以线程安全的方式将该值编码为int
func (a *AtomicInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Get())
}

//UnmarshalJSON以线程安全的方式将该值解码为int
func (a *AtomicInt64) UnmarshalJSON(b []byte) error {
	var into int64
	if err := json.Unmarshal(b, &into); err != nil {
		return err
	}
	a.Set(into)
	return nil
}

func (a *AtomicInt64) Get() int64 {
	return atomic.LoadInt64(&a.val)
}

//String以线程安全的方式将整数作为字符串返回
func (a *AtomicInt64) String() string {
	return strconv.FormatInt(a.Get(), 10)
}

//用新值交换当前值
func (a *AtomicInt64) Swap(newValue int64) int64 {
	return atomic.SwapInt64(&a.val, newValue)
}

func (a *AtomicInt64) Add(value int64) int64 {
	return atomic.AddInt64(&a.val, value)
}

func (a *AtomicInt64) Set(value int64) {
	atomic.StoreInt64(&a.val, value)
}

func (a *AtomicInt64) CompareAndSwap(expected int64, newVal int64) bool {
	return atomic.CompareAndSwapInt64(&a.val, expected, newVal)
}

func (a *AtomicInt64) Duration() time.Duration {
	return time.Duration(a.Get())
}
