## 简介

cmdctl是一个高仿kubectl的go语言的命令行框架，展示了些CLI的常用demo, 能够自动生成子命令。

## Features

+ 配置文件解析
+ 子命令
+ 参数
+ 自动补全（子命令和参数）
+ 各种编写子命令技巧
    + options
    + group
+ 生成子命令
+ 数据库操作
+ http调用
+ 版本信息


## 使用方式，详见`make help`

+ `make gotool` 运行go tool
+ `make cmd name=newctl` 在$GOPATH/src下新建命令
+ `make clean` 执行清理工作
+ `make install` 安装命令

## TODO
