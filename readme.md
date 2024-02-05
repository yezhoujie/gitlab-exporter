# gitlab-exporter

A tool to batch export all projects in gitlab to aliyun OSS bucket

## [中文文档](https://github.com/yezhoujie/gitlab-exporter/blob/main/readme-cn.md)

## usage

```
go build .

./gitlab-backup -f config.yaml
```

## config
config.yaml

| name | description |
|--|--|
|gitlab.url| your gitlab website url|
|gitlab.token| your gitlab access token |
|oss.accessKeyId|your aliyun accessKeyId with oss authentication|
|oss.accessSecret|you aliyun accessKeyId secret|
|oss.bucketName|the bucketName of aliyun oss which you want to store the export files|
|oss.endpoint|the endpoint of your aliyun oss|
|keepLocalBackup| if keep the export file in your local storage|

## flags with command


| flags | optional | description | default value|
|--|--|--|--|
|-h|yes|show help doc| Na|
|-f|yes|custom your config file path| ./config.yaml|
|-fromId|yes| only export projects which id smaller than the value, to avoid export duplicate projects| -1|


## tips
When stucking at a project exporting, it probably because the gitlab exporter worker is stucking there due to the gitlab server's capacity.
you can visit your gitlab project page then "Settings -> General -> Export project(Expand)" and manual click the export button and then wait for seconds to see it will come back.