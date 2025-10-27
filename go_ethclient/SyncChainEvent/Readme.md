# 介绍
- 本系统使用go语言进行追踪区块链上合约事件，讲区块链上的事件同步到数据库中，方便后台服务使用。

# 使用第三方库
1. **github.com/spf13/viper** 支持读取多种格式的配置文件
    ```shell
    go get -u github.com/spf13/viper
    ```
2. **github.com/mitchellh/go-homedir** 获取当前y用户主目录
    ```shell
    go get -u github.com/mitchellh/go-homedir
    ```
3. **go.uber.org/zap** 日志
    ```shell
    go get -u go.uber.org/zap
    ```
4. **lumberjack** 日志轮换
    ```shell
    go get -u gopkg.in/natefinch/lumberjack.v2
    ```
5. **gorm** 数据库框架
    ```shell
    go get -u gorm.io/gorm
    ```
6. **gorm.io/driver/mysql** 数据库驱动
    ```shell
    go get -u gorm.io/driver/mysql
    ```
7. **github.com/go-stack/stack** 堆栈库
    ```shell
    go get -u github.com/go-stack/stack
    ```
8. **go-zero**
    ```shell
    go get -u github.com/zeromicro/go-zero
    ```
9. **go-ethereum** 以太坊SDK
    ```shell
    go get -u github.com/ethereum/go-ethereum
    ```
10. **cron** 定时器框架
    ```shell
    go get -u github.com/robfig/cron/v3@v3.0.0
    ```

# 使用说明
程序运行前需修改配置config_template下的文件，并将配置文件放到程序目录下。需要在数据库中的tbl_sync_record表中添加要第一次同步的记录，以确定程序开始同步时间。在tbl_chain_event表中添加同步事件的hash值。在tbl_point_count_record_sepolia表中添加积分计算开始时间。