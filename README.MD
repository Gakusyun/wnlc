# 武昌工学院校园网登录器命令行客户端

此工具适用于在没有浏览器的设备上登录校园网。使用 Openai 的 o1-preview 模型与Go语言开发。本程序Python语言版[Gitee](https://gitee.com/Gakusyun/wuit-autologin)

使用方法：

1. 编译此项目
   ```shell
   go build -o wnlc
   ```
2. 把可执行文件移动到bin目录下
   ```shell
   mv wnlc /usr/local/bin
   ```
3. 运行命令 `wnlc`