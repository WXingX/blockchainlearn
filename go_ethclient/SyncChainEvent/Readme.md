# 介绍
- 本系统使用go语言进行追踪区块链上合约事件，讲区块链上的事件同步到数据库中，方便后台服务使用。

# 使用第三方库
1. **github.com/spf13/cobra**   用于构建命令行工具
    ```shell
        go get -u github.com/spf13/cobra
    ```
2. **github.com/spf13/viper** 支持读取多种格式的配置文件
    ```shell
    go get -u github.com/spf13/viper
    ```
3. **github.com/mitchellh/go-homedir** 获取当前y用户主目录
    ```shell
    go get -u github.com/mitchellh/go-homedir
    ```
4. **go.uber.org/zap** 日志
    ```shell
    go get -u go.uber.org/zap
    ```
5. **lumberjack** 日志轮换
    ```shell
    go get -u gopkg.in/natefinch/lumberjack.v2
    ```
6. **gorm** 数据库框架
    ```shell
    go get -u gorm.io/gorm
    ```
7. **gorm.io/driver/mysql** 数据库驱动
    ```shell
    go get -u gorm.io/driver/mysql
    ```
8. **github.com/go-stack/stack** 堆栈库
    ```shell
    go get -u github.com/go-stack/stack
    ```
9. **go-zero**
    ```shell
    go get -u github.com/zeromicro/go-zero
    ```
10. **go-ethereum** 以太坊SDK
    ```shell
    go get -u github.com/ethereum/go-ethereum
    ```
11. **cron** 定时器框架
    ```shell
    go get -u github.com/robfig/cron/v3@v3.0.0
    ```
12. 