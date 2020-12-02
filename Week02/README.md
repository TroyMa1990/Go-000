# go的error处理
 
go的哲学是只处理一次错误，而不是处理成功，也不是简单的抛出异常，
而且可追溯、可回放的错误记录才是有价值的。

我们可以通过案例去更好地理解这一点。

## 案例学习一
Sentinel Error
预定义的特定错误
缺点：
1、成为你 API 公共部分，如果 API 定义了一个返回特定错误的 interface，则该接口的所有实现都将被限制为仅返回该错误，即使它们可以提供更具描述性的错误。
比如 io.Reader。像 io.Copy 这类函数需要 reader 的实现者比如返回 io.EOF 来告诉调用者没有更多数据了，但这又不是错误。
2、在两个包之间创建了依赖，想要判断一个错误类型必须导入包。

## 案例学习二
Error types
Error type 是实现了 error 接口的自定义类型。


通过自定义错误传递指针的方式，来使用断言转换成特定错误类型，获取更多上下文信息
缺点：调用者要使用类型断言和类型 switch，就要让自定义的 error 变为 public。这种模型会导致和调用者产生强耦合，从而导致 API 变得脆弱。
```go
package main

import (  
    "fmt"
    "math"
)

type areaError struct {  
    err    string
    radius float64
}

func (e *areaError) Error() string {  
    return fmt.Sprintf("radius %0.2f: %s", e.radius, e.err)
}

func circleArea(radius float64) (float64, error) {  
    if radius < 0 {
        return 0, &areaError{"radius is negative", radius}
    }
    return math.Pi * radius * radius, nil
}

func main() {  
    radius := -20.0
    area, err := circleArea(radius)
    if err != nil {
        if err, ok := err.(*areaError); ok {
            fmt.Printf("Radius %0.2f is less than zero", err.radius)
            return
        }
        fmt.Println(err)
        return
    }
    fmt.Printf("Area of rectangle1 %0.2f", area)
}
```

## 案例学习三
Opaque errors
不透明错误，耦合性最低，在少数情况下，这种二分错误处理方法是不够的。例如，与进程外的世界进行交互(如网络活动)，需要调用方调查错误的性质，以确定重试该操作是否合理。在这种情况下，我们可以断言错误实现了特定的行为，而不是断言错误是特定的类型或值。
```go
type temporary interface {
    Temporary() bool
}
func IsTemporary(err error) bool {
    te, ok := err.(temporary)
    return ok && te.Temporary()
}

```
这里的关键是，这个逻辑可以在不导入定义错误的包或者实际上不了解 err 的底层类型的情况下实现——我们只对它的行为感兴趣。

## 案例学习四

位于github.com/pkg/errors 的errors包是Go标准库的替代品。
它提供了一些非常有用的操作用于封装和处理错误。值得注意的是，使用标准库的方式对错误进行封装，会改变其类型并使类型断言失败。

建立errwrap.go ：

```go
package errwrap

import (
	"fmt"

	"github.com/pkg/errors"
)

// WrappedError 演示了如何对错误进行封装
func WrappedError(e error) error {
	return errors.Wrap(e, "An error occurred in WrappedError")
}

type ErrorTyped struct {
	error
}

func Wrap() {
	e := errors.New("standard error")

	fmt.Println("Regular Error - ", WrappedError(e))

	fmt.Println("Typed Error - ", WrappedError(ErrorTyped{errors.New("typed error")}))

	fmt.Println("Nil -", WrappedError(nil))

}
```
建立unwrap.go ：

```go
package errwrap

import (
	"fmt"

	"github.com/pkg/errors"
)

// Unwrap 解除封装并进行断言处理
func Unwrap() {

	err := error(ErrorTyped{errors.New("an error occurred")})
	err = errors.Wrap(err, "wrapped")

	fmt.Println("wrapped error: ", err)

	// 处理错误类型
	switch errors.Cause(err).(type) {
	case ErrorTyped:
		fmt.Println("a typed error occurred: ", err)
	default:
		fmt.Println("an unknown error occurred")
	}
}

// StackTrace 打印错误栈
func StackTrace() {
	err := error(ErrorTyped{errors.New("an error occurred")})
	err = errors.Wrap(err, "wrapped")

	fmt.Printf("%+v\n", err)
}
```
建立main.go ：
```go
package main

import (
	"fmt"

	"github.com/xxx/errwrap"
)

func main() {
	errwrap.Wrap()
	fmt.Println()
	errwrap.Unwrap()
	fmt.Println()
	errwrap.StackTrace()
}
```
这会输出 ：

```js
Regular Error - An error occurred in WrappedError: standard
error
Typed Error - An error occurred in WrappedError: typed error
Nil - <nil>
wrapped error: wrapped: an error occurred
a typed error occurred: wrapped: an error occurred
an error occurred
github.com/agtorre/go-cookbook/chapter4/errwrap.StackTrace
/Users/lothamer/go/src/github.com/agtorre/gocookbook/chapter4/errwrap/unwrap.go:30
main.main
/tmp/go/src/github.com/agtorre/gocookbook/chapter4/errwrap/example/main.go:14
```
对错误为nil的情况进行处理：

```go
func RetError() error{
     err := ThisReturnsAnError()
     return errors.Wrap(err, "This only does something if err != nil")
}
```

## 作业

题目：我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

答：dao层中当遇到一个sql.ErrNoRows错误时，需要Wrap这个error并且抛给上层，因为如果dao层吞掉这个sql.ErrNoRows错误，
可能会使上层调用者无法通过错误日志来回溯业务问题，而且有些时候业务逻辑会把某个数据为空的情况也当做一种业务问题来处理，会有正常流、空数据和错误流等多种处理方式。

[作业代码](./demo.go)
 