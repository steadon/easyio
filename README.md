# EasyIO

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![Docker Pulls](https://img.shields.io/docker/pulls/steadon/easyio?color=green)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/steadon/easyio)

一个简易的对象存储服务后台，可以轻易借助宿主机实现对象存储服务

## 1. 开始

### 1.1 启动容器

- win/linux系统大多数使用以下指令

```
docker run --name your-easyio -p 8000:8000 -d steadon/easyio:1.0-amd64
```

- m1/m2芯片mac系统使用以下指令

```
docker run --name your-easyio -p 8000:8000 -d steadon/easyio:1.0-arm64
```

### 1.2 其他准备

- 本项目所有图片资源都存放在 `images` 目录下，正式部署需要挂载数据卷到宿主机，否则容器退出将导致数据丢失

```
-v /local/images:/app/images    //前提是已经在本地创建了/local/images文件夹
```

- 本项目配置文件是位于 `config` 目录下的 `app.ini` 文件，涉及数据库等配置，如需要针对性管理也需要挂载出来

```
-v /local/config/app.ini:/config/app.ini    //前提是已经在本地创建了/local/config/app.ini文件
```

- 通过以下命令可下载 `app.ini` 文件到本地

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
