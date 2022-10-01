1. 首先開啟 zookeer 跟 kafka
    ```
    /usr/local/Cellar/kafka/3.2.1/bin/zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties
    /usr/local/Cellar/kafka/3.2.1/bin/kafka-server-start /usr/local/etc/kafka/server.properties
    ```

2. 接著開啟 etcd
    ```
    % etcd
    ```

3. 接著執行 main.go
    ```
    go run ./
    ```

4. 接著用 etcdctl 來更新配置，用來測試新增追蹤與刪除追蹤
    ```
    4.1 /usr/local/Cellar/etcd/3.5.5/bin/etcdctl put /192.168.43.38/ '[{"topic":"a","path":"/Users/jyunzn/golearn/src/test/test1.log"},{"topic":"b","path":"/Users/jyunzn/golearn/src/test/test2.log"}]'

    4.2 /usr/local/Cellar/etcd/3.5.5/bin/etcdctl put /192.168.43.38/ '[{"topic":"a","path":"/Users/jyunzn/golearn/src/test/test1.log"},{"topic":"b","path":"/Users/jyunzn/golearn/src/test/test2.log"}, {"topic":"c","path":"/Users/jyunzn/golearn/src/test/test3.log"}]'
    ```

5. 最後就可以測試更新追蹤的文件是否有送到 kafka
    ```
    直接更新 test1.log test2.log test3.log

    如果執行了 4.1，那麼 test3.log 應該沒效果，因為取消追蹤了，反之 4.2 就會有追蹤到而發送
    ```
