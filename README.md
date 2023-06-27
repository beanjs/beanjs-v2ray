# 快速搭建 v2ray 服务端

项目提供快速架设 v2ray 服务器的一种方法，您只需要提供一个 vultr 的 API KEY 就能全自动的搭建 v2ray 的服务器环境，本部署方法不需要您进行任何的服务器操作，仅需在 PC 上运行相应命令即可。

> 本项目仅限技术交流，请勿滥用。

## 编译与运行

```bash
go run . # 运行
go build . # 编译
```

## 使用方法

```bash
beanjs-v2ray -h # 查看帮助
beanjs-v2ray --vultrkey <API KEY> # 部署服务端
```
