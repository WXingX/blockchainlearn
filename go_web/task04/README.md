# 博客管理系统

### 一、介绍
本系统是一个基于gin框架，使用gorm库进行数据操作的博客管理系统后端。

主要实现博客文章的基本管理功能，包括文章的创建、读取、更新和删除（CRUD）操作，同时支持用户认证和简单的评论功能。

### 二、相关依赖
本系统使用依赖如下：
1. gopkg.in/yaml.v3            读取yaml配置文件
2. go.uber.org/zap             日志模块
3. lumberjack                  日志轮换
4. gorm                        数据库框架
5. gorm.io/driver/sqlite       sqlite数据库驱动
6. gin                         web框架
7. jwt/v5                      JWT 中间件
8. bcrypt                      加密工具
9. modernc.org/sqlite          go实现的sqlite驱动

安装方式
1. ```go get -u gopkg.in/yaml.v3```
2. ```go get -u go.uber.org/zap```
3. ```go get -u gopkg.in/natefinch/lumberjack.v2```
4. ```go get -u gorm.io/gorm```
5. ```go get -u gorm.io/driver/sqlite```
6. ```go get -u github.com/gin-gonic/gin```
7. ```go get -u github.com/golang-jwt/jwt/v5```
8. ```go get -u golang.org/x/crypto/bcrypt```
9. ```go get -u modernc.org/sqlite```


### 三、编译命令
在main.go文件所在目录运行以下命令
```go run ./main.go```

### 四、注意事项
1. 因为使用了sqlite，编译会有问题，需要设置 CGO_ENABLED 为1，动态链接 C 库，还是按照GCC。 当值为0 时，使用静态GO库。跨平台、容器、无依赖时设为 0。或者切换modernc.org/sqlite 驱动。本项目切换为了 modernc.org/sqlite驱动。

