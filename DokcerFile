# 基础镜像 golang
FROM golang
# 工作目录
WORKDIR /go/go-zentao-task-api
# 将文件copy到镜像相同目录中
COPY . .
# 暴露端口
EXPOSE 8899
# 执行命令
RUN go build main.go && go run ./main
