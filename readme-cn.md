# gitlab-exporter

一个能按照project,批量导出gitlab里所有project（你所在账号权限内）到阿里云OSS上的工具。

## 原理
- 调用gitlab API导出
- 下载本地
- 上传OSS

## 用法

```
go build .

./gitlab-backup -f config.yaml
```

## 配置
config.yaml

| 配置项 | 描述 |
|--|--|
|gitlab.url| gitlab的网址|
|gitlab.token| 提供接口调用的gitlab access token |
|oss.accessKeyId|拥有阿里云OSS写权限的accessKeyId|
|oss.accessSecret|阿里云 accessKeyId secret|
|oss.bucketName|希望存储的OSS bucket 名称|
|oss.endpoint|对应OSS bucket 的 endpoint 地址|
|keepLocalBackup|是否保留本地导出文件|

## 命令行参数


| 参数 | 缺省 | 描述 | 默认值|
|--|--|--|--|
|-h|yes|显示使用帮助| 无|
|-f|yes|自定义的配置文件路径| ./config.yaml|
|-fromId|yes|当设置了该参数，程序会只导出小于等于该ID的project, 避免导出过程中程序中断等其他原因重复导出| -1|


## 提示
如果导出过程中出现程序卡住不动，可能是由于gitlab服务器资源配置问题，自己的导出任务卡住了，这个时候你需要访问到卡住的项目页，然后在"Settings -> General -> Export project(Expand)"中手动点击Export按钮，让gitlab重新做一次导出操作， 这个时候再稍等一会看程序，应该就顺利往下执行了.