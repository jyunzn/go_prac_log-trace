package kafkalog

import (
	"time"
	"fmt"

	"github.com/Shopify/sarama"
)


// 要發給 kafka 的元數據
type logData struct {
	topic string
	data string
}


var (
	client sarama.SyncProducer
	logDataChan chan *logData   // 緩衝通道，用來存儲準備發給 kafka 的數據們
)


// 初始化 kafka api
func Init(addr string, maxsize int) (sarama.SyncProducer, error) {

	// 初始化 kafka
	var err error
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer([]string{ addr }, config)

	// 初始化 log 緩衝通道
	logDataChan = make(chan *logData, maxsize)


	// 起一個進程專門去取通道中的數據，發送到 kafka 中
	go sendMsgToKafka()

	return client, err
}


// 把數據塞到通道中
func SendToChan(topic, msg string) {
	ld := &logData {
		topic: topic,
		data: msg,
	}

	logDataChan <- ld
}


// 從通道中取數據出來發給 kafka
func sendMsgToKafka() {
	for {
		select {
			case ld := <- logDataChan:
				fmt.Println(ld.topic, ld.data)
				msgObj := &sarama.ProducerMessage {
					Topic: ld.topic,
					Value: sarama.StringEncoder(ld.data),
				}
				pid, offset, err := client.SendMessage(msgObj)
				if err != nil {
					fmt.Println("發送信息至 kafka 有誤", err)
					time.Sleep(time.Second)
					continue
				}
				fmt.Println("kafka 發送成功: ", pid, offset)

			default:
				time.Sleep(time.Second)
		}
	}
}
