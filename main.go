package main

import (
	"fmt"
	"time"
	"sync"

	"test/kafkalog"
	"test/taillog"
	"test/config"
	"test/etcd"
	"test/utils"

	"gopkg.in/ini.v1"

)

var (
	cnf = new(config.AppConf)   // * 初始化配置對象
)


func main() {

	// * 加載配置文件
	err := ini.MapTo(cnf, "./config/config.ini")
	if err != nil {
		fmt.Println("加載配置失敗: ", err)
		return
	}


	// * 初始化 kafka
	kafkaClient, err := kafkalog.Init(cnf.KafkaConf.Address, cnf.KafkaConf.MaxBuffSize)
	if err != nil {
		fmt.Println("kafka connect err: ", err)
		return
	}
	defer kafkaClient.Close()
	fmt.Println("kafka 連接成功")


	// * 初始化 etcd
	etcdClient, err := etcd.Init(
		cnf.EtcdConf.Address,
		time.Duration(cnf.EtcdConf.Timeout) * time.Second,
	)
	if err != nil {
		fmt.Println("etcd connect err", err)
		return
	}
	defer etcdClient.Close()
	fmt.Println("etcd 連接成功")



	// * 獲取配置
	// > 先拿到本機 IP，用來做配置 etcd 的 key，這樣不同的機器就有不同的 key
	// > 不過通常不會直接拿 IP 區分，而會拿業務項來區分
	ip, err := utils.GetOutboundIP()
	if err != nil {
		fmt.Println("獲取 IP 是敗", err)
		return
	}
	key := fmt.Sprintf(cnf.EtcdConf.Key, ip)  // 配置有留 %s 用來拼接的，拿 IP 去塞佔位符
	fmt.Println(key)

	// > 根據 key 獲取配置數據
	logEntries, err := etcd.GetConfig(key)
	if err != nil {
		fmt.Println("獲取配置失敗")
		return
	}

	// * 追蹤 etcd 配置裡追蹤的所有文件
	taillog.Init(logEntries)

	// * 獲取追蹤的監聽通道
	newConfChan := taillog.GetNewConfChan()


	// * 持續監聽 etcd 數據
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConfig(key, newConfChan)
	wg.Wait()
}

