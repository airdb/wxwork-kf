# 微信客服

## 快速开始

### 依赖检查

- 操作系统：Linux, MacOS
- 开发语言：Golang
- 部署工具：[Serverless](https://github.com/serverless/serverless)
- 正常访问外网

### 环境文件

1. .env

serverless 登录腾讯云的凭据配置，以及针对 serverless 部署到腾讯云必要环境变量(SERVERLESS_PLATFORM_VENDOR=tencent)配置。

其中 SERVERLESS_PLATFORM_VENDOR=tencent 需要手动添加。

凭据为 serverless 工具自动追加的。

2. .env.(test|release)

对应环境下的针对应用的配置信息，具体参考 .env.example 文件。

### 构建

1. 代码包下载

```
$ git clone https://gitee.com/airdb/wxwork-kf
```

2. 编译

```bash
$ make build
```

### 部署

1. 测试环境

必要配置文件清单：.env, .env.test

```bash
$ make deploy
```

2. 正式环境

必要配置文件清单：.env, .env.release

```bash
$ make release
```
