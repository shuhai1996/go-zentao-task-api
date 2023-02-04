# 基础镜像 golang
FROM golang
# 运行环境
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=1
# 工作目录
WORKDIR /go/src/go-zentao-task-api
# 将文件copy到镜像相同目录中
COPY . .
# 暴露端口
EXPOSE 8899
# 执行命令
RUN go mod tidy && go build main.go && ./main
