package etcd

import (
	"time"
	"context"
	"fmt"

	"go.etcd.io/etcd/client/v3"

	"encoding/json"
)

var (
	client *clientv3.Client
)


/**
 * = 預期存在 etcd 的 value 是一個 json
 */
type LogEntry struct {
	Topic string `json="topic"` // 日誌的路徑
	Path string `json="path"`   // 日誌要發往 kafka 的 topic
}

func Init(addr string, timeout time.Duration) (*clientv3.Client, error) {
	var err  error

	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{ addr },
		DialTimeout: timeout,
	})

	return client, err
}



/**
 * @desc 獲取配置信息
 * @param { key } 表示該機器的鍵值，不同的機器就傳入不同的 key
 * @ps
 * = 一台機器可能有多個日誌需要搜集，所以返回的是一個 logEntries slice，裝載著該機器想要搜集的日誌文件
 * = 也就是說機器 key -> value [{"topic":"日誌的主旨","path":"該日誌的路徑"}, ...]
 */
func GetConfig(key string)(logEntries []*LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	gresp, err := client.Get(ctx, key)
	cancel()
	if err != nil {
		return
	}

	for _, evt := range gresp.Kvs {
		err = json.Unmarshal(evt.Value, &logEntries)
		if err != nil {
			fmt.Println("json unmarshal err: ", err)
			return
		}
	}

	return
}


/**
 * 持續監聽存儲在 etcd 的配置改變
 */
func WatchConfig(key string, newConfChan chan <- []*LogEntry) {
	ch := client.Watch(context.Background(), key)
	for wresp := range ch {
		for _, evt := range wresp.Events {
			fmt.Printf("%s, %s, %s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)

			// 如果不是刪除操作的話
			if evt.Type != clientv3.EventTypeDelete {
				var newConf []*LogEntry
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					fmt.Println("json unmarshal err: ", err)
					continue
				}
				newConfChan <- newConf  // 把新配置塞到配置更新通道中
			}
		}
	}
}
