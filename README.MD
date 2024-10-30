# 武昌工学院校园网登录器命令行客户端

此工具适用于在没有浏览器的设备上登录校园网。使用 Openai 的 o1-preview 模型与Go语言开发。本程序Python语言版[Gitee](https://gitee.com/Gakusyun/wuit-autologin)

使用方法：

1. 编译此项目
   ```shell
   go build -o wnlc
   ```
2. 把可执行文件移动到 /usr/bin 目录下
   ```shell
   mv wnlc /usr/bin
   ```
3. 运行命令 `wnlc`

或者使用 make 命令一键安装

```shell
make install
```
也可使用 make uninstall 卸载

make 参数：

| 参数      | 说明 |
| --------- | ---- |
| install   | 安装 |
| uninstall | 卸载 |
| clean     | 清理 |

> TODO暂未实现:
>
> wnlc 参数
> 
> | 参数 | 说明         |
> | ---- | ------------ |
> | -s   | 获取当前状态 |
> | -l   | 仅登录       |
> | -o   | 退出登录     |