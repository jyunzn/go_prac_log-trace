package taillog

import (
	"fmt"
	"time"
	"context"

	"github.com/hpcloud/tail"

	"test/kafkalog"
)

type TailTask struct {
	topic string
	path string
	instance *tail.Tail          // 追蹤的模塊實例
	ctx context.Context
	cancelFn context.CancelFunc  // 用來停止追蹤模塊的 gorountine
}

func (t *TailTask)init() {
	config := tail.Config {
		Follow: true,
		ReOpen: true,
		Location: &tail.SeekInfo { Offset: 0, Whence: 2 },
		MustExist: false,
		Poll: true,
	}

	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("追蹤檔案有誤: ", err)
		return
	}

	go t.run()
}

func (t *TailTask)run() {
	for {
		select {
			// 如果這個配置被刪掉時，結束 gorounine
			case <- t.ctx.Done():
				fmt.Printf("topic: %s 對應的 path: %s 即將結束監聽\n", t.topic, t.path)
				return

			// 如果追蹤的 log 有新的數據，就發送到 kafka 的緩衝通道中
			case line := <- t.instance.Lines:
				// kafkalog.SendMsgToKafka(t.topic, line.Text)  // => 直接發送到 kafka 的話，這裡就需要同步去等他發送完
				kafkalog.SendToChan(t.topic, line.Text)         // => 所以就改成將數據都先發到一個通道中，然後另開一個進程去那個通道不停的取東西出來發給 kafka
			default:
				time.Sleep(time.Second)
		}
	}
}



func NewTailTask(topic, path string) (*TailTask){
	ctx, cancelFn := context.WithCancel(context.Background())
	tailTaskObj := &TailTask {
		topic: topic,
		path: path,
		ctx: ctx,
		cancelFn: cancelFn,
	}

	tailTaskObj.init()
	return tailTaskObj
}
