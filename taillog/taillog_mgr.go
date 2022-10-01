package taillog

import (
	"test/etcd"
	// "time"
	"fmt"
)

var tskMgr *taillogMgr

type taillogMgr struct {
	logEntries []*etcd.LogEntry         // 存儲最原始的配置信息
	tskMap map[string]*TailTask         // 存儲實例出來的監聽任務實例
	newConfChan chan []*etcd.LogEntry   // 配置更新通道
}


func Init(logEntries []*etcd.LogEntry) {
	tskMgr = &taillogMgr {
		logEntries: logEntries,                    // 最原始的配置項（其實沒什麼用）
		tskMap: make(map[string]*TailTask, 16),    // 用來快速比對查找用的 map，key 為 topic_path，值為對應的任務
		newConfChan: make(chan []*etcd.LogEntry),  // 無緩存通道，目的在於沒任務就讓該 gorountine 堵塞
	}

	// 生成每個追蹤任務，並用 topic_path 存起來
	for _, logentry := range logEntries {

		key := fmt.Sprintf("%s_%s", logentry.Topic, logentry.Path)
		tskMgr.tskMap[key] = NewTailTask(logentry.Topic, logentry.Path)
	}

	go tskMgr.run()
}


// 追蹤 etcd 數據更新通道
func (t *taillogMgr)run() {
	for {
		newConf := <- t.newConfChan
		fmt.Println("== 配置更新 ==")

		// 查找原本有而新配置沒有的，要刪除
		for originKey := range tskMgr.tskMap {
			willDelete := true


			for _, conf := range newConf {
				newConfKey := fmt.Sprintf("%s_%s", conf.Topic, conf.Path)
				if originKey == newConfKey {
					willDelete = false
					break
				}
			}

			if willDelete {
				tskMgr.tskMap[originKey].cancelFn()  // 停止追蹤
				delete(tskMgr.tskMap, originKey)     // 從對應表中移除
			}
		}

		// 查找原本沒有而新配置有的，要新增
		for _, conf := range newConf {
			key := fmt.Sprintf("%s_%s", conf.Topic, conf.Path)

			_, ok := tskMgr.tskMap[key]
			if !ok {
				tskMgr.tskMap[key] = NewTailTask(conf.Topic, conf.Path)
			}
		}

	}
}


// 暴露追蹤通道，要給 etcd 模塊塞更新的數據近來
func GetNewConfChan() chan <- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
