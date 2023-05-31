# EasyIO

一个简易的对象存储服务后台，可以轻易借助宿主机实现对象存储服务

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![Docker Pulls](https://img.shields.io/docker/pulls/steadon/easyio?color=green)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/steadon/easyio)

## 1. 开始

### 1.1 启动容器

win/linux系统大多数使用以下指令

```
docker run --name your-easyio -p 8000:8000 -d steadon/easyio:1.0-amd64
```

m1/m2芯片mac系统使用以下指令

```
docker run --name your-easyio -p 8000:8000 -d steadon/easyio:1.0-arm64
```

### 1.2 其他准备

- 本项目所有图片资源都存放在 `images` 目录下，正式部署需要挂载数据卷到宿主机，否则容器退出将导致数据丢失

- 本项目配置文件是位于 `config` 目录下的 `app.ini` 文件，涉及数据库等配置，如需要针对性管理也需要挂载出来

## 2. 使用

本项目支持接口调用以及配套客户端操作两种方式使用...

该项目正在加急研发中...
