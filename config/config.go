package config


type KafkaConf struct {
	Address string `ini:"address"`
	MaxBuffSize int `init:"max_buff_size"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Timeout uint `ini:"timeout"`
	Key string `ini:"collect_log_key"`
}

type AppConf struct {
	KafkaConf `ini:"kafkalog"`
	EtcdConf `ini:"etcd"`
}
