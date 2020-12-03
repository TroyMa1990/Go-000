package main

import (
	"database/sql"
	"log"
	"github.com/pkg/errors"
	"time"

	codes "xxx.xxx.com/biz/codes"
	control "xxx.xxx.com/biz/control"
	pb "xxx.xxxx.com/biz/protos"

)

var (
	ctrl  	 *control.Ctrl
	BizId 	 string = "1"
	identKey string = "CycleIdent"
	cnt      int64  = 0
)

type GrpcResponse struct {
	bizId string
	ord   *pb.Order
	err   error
}

type Goods struct {
	Id        int64   `json:"id"`
	Extension struct{} `json:"ext"`
}

func Dao() (*Goods, error) {
	return nil, errors.Wrap(sql.ErrNoRows, "Can Not Find Data")
}

func Service() (*Goods, error) {
	return Dao()
}

func Business(ch chan GrpcResponse,ident int) error {
	//此处消费数据库
	resp, err := Service()
	//直接交给实时的业务处理队列
	ch <- GrpcResponse{bizId: BizId, ord: resp, err: err} 
	if  err != nil {
		//记录错误堆栈日志
		log.Printf("Ident[%d] Some Error To Query Database With Stack:%v\n", cnt, err)
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
					if res.err == nil {
						//这里处理正常数据业务逻辑
						...
						log.Printf("Ident[%d] BizId[%v] BizMissionQueue Current Handle: %v", cnt, BizId, err)
						continue
					}
					if errors.Is(res.err, sql.ErrNoRows) {
						//这里处理空数据业务逻辑
						...
						log.Printf("Ident[%d] BizId[%v] BizMissionQueue Empty Handle: %v", cnt, BizId, err)
					} else {
						//这里处理数据库错误的降级业务逻辑
						...
						log.Printf("Ident[%d] BizId[%v] BizMissionQueue Degradation Handle: %v", cnt, BizId, err)
					}	
				case <-time.After(30 * time.Second):
					//检测业务BizMissionQueue停止信号
					err := GetQueueStop(BizId)
					if err != nil {
						//记录业务错误码
						log.Printf("Ident[%d] BizId[%v] BizMissionQueue Stop: %v", cnt, BizId, err)
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
	ctrl, err = control.NewControler()
	if err != nil {
		log.Printf("Failed To Init Status Redis: %v", err)
	}
	defer ctrl.Close()
	for {
		//获取redis持久化执行轮次，每次加一
		cnt, err = ctrl.GetIdent(context.TODO(), identKey)
		log.Printf("Ident[%d] Main Run Stage", cnt)
		//此处检测redis业务暂停信号
		for {
			if owner, err := ctrl.GetStop(context.TODO()); err == codes.BizStatusExists {
				time.Sleep(3 * time.Second)
				continue
			}
			break
		}
		// 保持业务消费持续进行
		if err := Business(ch,ident); err != codes.OK {
			time.Sleep(3 * time.Second)
			//记录业务错误码
			log.Printf("Ident[%d] Some Error When BizId[%v] Is Going:%v", cnt, BizId, err)
		}
	}
}
