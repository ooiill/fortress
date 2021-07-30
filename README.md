## 简介

Mini fortress program

## 安装

```shell
$ go get github.com/ecccho/fortress
$ cp $GOPATH/src/github.com/ecccho/fortress/fortress.json ~/.fortress.json
$ go install github.com/ecccho/fortress
```

## 运行
```
➜  ~ fortress -h                   // 打印参数列表
➜  ~ fortress                      // 进入交互式执行
➜  ~ fortress --index=1,2,3
➜  ~ fortress --mode=1             // 半交互式
➜  ~ fortress --index=1 --mode=1   // 自动执行，后面是数字与交互式菜单索引保持一致
```

1. `fortress.json > host` 定义相关主机信息；  

2. `fortress.json > fortress` 定义相关跳板任务；  

3. 请将本机 `~/.ssh/id_rsa.pub` 追加到远程主机的 `~/.ssh/authorized_keys` 文件中来避免交互输入密码；  

4. `fortress.json > hosts > password` 若被设置，将使用 `sshpass` 软件执行命令，可用于替代上一条事项（前提是本机已安装该软件）；  

5. `fortress.json > fortress > type` 支持 `ssh`、`sync`、`normal(缺省)`，具体配置项见示例配置文件；  
 
> Tips: 如果使用 `windows` 系统建议使用 `git bash` 执行；  

## 样例

![操作台](./demo.png)

## 作者

jiangxilee@gmail.com