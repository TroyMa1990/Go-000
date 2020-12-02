package main

import (
	"database/sql"
	"log"
	"github.com/pkg/errors"
	
	codes "xxx.xxx.com/biz/codes"
	control "xxx.xxx.com/biz/control"
	pb "xxx.xxxx.com/biz/protos"

)

var (
	chans        sync.Map
	ctrl         *control.Ctrl
	BizId        =1
)

type GrpcResponse struct {
	bizId 	string
	ord   	*pb.Order
	err   	error
}

type Goods struct {
	Id        	uint64      `json:"id"`
	Extension 	struct{}    `json:"ext"`
}

func Dao() (*Goods, error) {
	return nil, errors.Wrap(sql.ErrNoRows, "can not find data")
}

func Service() (*Goods, error) {
	return Dao()
}

func Business(ch chan GrpcResponse) error {
	ctrl, err = control.NewControler()
	if err != nil {
		log.Printf("Failed To Init Status Redis: %v", err)
	}
	defer ctrl.Close()
	//此处检测redis业务暂停信号
	for {
		if owner, ss := ctrl.GetStop(context.TODO()); ss == codes.BizStatusExists {
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	//此处消费数据库
	resp, err := service()
	//直接交给实时的业务处理队列
	ch <- GrpcResponse{bizId: bizId, ord: resp, err: err} 
	//记录日志
	if  err != nil {
		log.Printf("Some Error To Query Database With Stack:%v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return codes.BizMissionQueueDBErrNoRows
		} else {
			return codes.BizMissionQueueDBError
		}
	}
	
	return codes.OK
}

func BizMissionQueue(ch chan GrpcResponse) {
	go func() {
		for {
			var res GrpcResponse
			select {
			case res = <-ch:
				if res.err ==nil {
					//这里处理正常数据业务逻辑
					...
					continue
				}
				if errors.Is(res.err, sql.ErrNoRows) {
					//这里处理空数据业务逻辑
					...
				} else {
					//这里处理数据库错误的降级业务逻辑
					...
				}	
			case <-time.After(30 * time.Second):
				//检测业务BizMissionQueue停止信号
				err := GetQueueStop(BizId)
				if err != nil {
					log.Printf("BizId[%v] Stop: %v", res.BizId, err)
					return 
				}
				continue
			}
		}
	}
}

func main() {
	ch := make(chan GrpcResponse, 0)
	// 开启业务处理队列
	BizMissionQueue(ch)
	// 保持业务消费持续进行
	for {
		if err := Business(ch); err != codes.OK {
			time.Sleep(3 * time.Second)
			//实际上是记录错误码
			log.Printf("Some Error When Biz Is Going:%v", err)
		}
	}
}
