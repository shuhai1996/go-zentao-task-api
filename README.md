## go-zentao-task

#### 配置
根目录下创建conf.ini文件，添加如下配置
```
[development]
db.host =localhost
db.port =3306
db.username =root
db.password =
db.database =
```
#### 包管理
新的项目，复制完代码以后，全局替换go-zentao-task为你的项目名称，删除go.mod和go.sum文件，执行下列命令：
```
go mod init go-zentao-task
go mod vendor //将下载到GOPATH的包复制到当前项目的vendor目录下
```
#### 运行本项目
执行如下命令即可运行
```
go build main.go
./main
```

> 参考链接
>
> [https://github.com/go-gorm/gorm](https://github.com/go-gorm/gorm)
> 
> [GORM 2.0 使用教程(中文文档)](https://www.bookstack.cn/read/gorm-2.0/docs-index.md)
> 
> [https://gin-gonic.com/](https://gin-gonic.com/)
