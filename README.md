# EasyIO

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![Docker Pulls](https://img.shields.io/docker/pulls/steadon/easyio?color=green)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/steadon/easyio)

一款轻量级对象存储服务器软件，它可以轻易地部署到任何服务器上并提供稳定快速的对象存储服务。
- Author：[steadon](https://github.com/steadon)

我想以下几点可能会成为你使用EasyIO的原因：
- 使用 `Go` 和 `Gin` 开发，同时摒弃了数据库带来的厚重依赖，性能及独立性表现卓越
- 通过 `Docker` 部署，支持 `Linux/amd` 和 `Linux/arm` 两种主流系统架构：[Docker Hub - steadon/easyio](https://hub.docker.com/repository/docker/steadon/easyio/general)
- 提供了一套基于 `Java17` 和 `SpringBoot 3.1.0` 的SDK：[GitHub - steadon/easyio-sdk-java](https://github.com/steadon/easyio-sdk-java)
- 同相似产品一样也提供了一套勉强能看的可视化界面：<链接待补充>

## 1. 开始

### 1.1 启动容器

- 使用以下指令启动一个最简单的容器实例

```
docker run --name your-easyio -p 8000:8000 -d steadon/easyio:1.0
```

### 1.2 其他准备

- 图片资源存放在 `images` 目录下，正式部署需要挂载数据卷到宿主机，否则容器退出将导致数据丢失

```
-v /local/images:/app/images    //前提是已经在本地创建了/local/images文件夹
```

- 配置文件是位于 `config` 目录下的 `app.ini` 文件，涉及root账号密码等配置，正式部署需要挂载出来

```
-v /local/config/app.ini:/config/app.ini    //前提是已经在本地创建了/local/config/app.ini文件
```

- 通过以下命令可下载 `app.ini` 文件到本地并进行编辑

```
wget https://raw.githubusercontent.com/steadon/EasyIO/main/conf/app.ini
```

## 2. 使用

### 2.1 通过接口调用

- 检索根下目录 GET /action/show/root

```
无需参数
```

- 创建指定目录 POST /action/add/dir

```
{
    "name": "string"    //目录路径，例如 name/group
}
```

- 上传一张图片 POST /action/upload

```
form-data: file file       //上传的图片，支持常见类型如.jpeg .png .img
form-data: group string    //目录路径，例如 name/group
form-data: name string     //图片名，可带后缀，不传则用随机串代替
```

- 查看目录列表 GET /action/show/dir

```
param: group string    //目录路径，例如 name/group，将查看该路径下的所有目录
```

- 查看图片列表 GET /action/show/img

```
param: group string    //目录路径，例如 name/group，将查看该路径下的所有图片
```

- 删除指定图片 DEL /action/delete/img

```
param: path string    //文件路径，例如 name/group/xxx.png，将删除该图片
```

- 删除指定目录 DEL /action/delete/dir

```
param: path string    //目录路径，例如 name/group，将删除该目录及其所有子文件
```

- 用户登录 POST /user/login

```
{
    "username": "root",    //用户名
    "password": "123456"   //密码
}
// 返回token，访问所有/action的接口都需要在请求头带上 "Authorization":"token"
```

---

该项目正在加急研发中...
