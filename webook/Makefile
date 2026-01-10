# 该文件就是 编译成linux文件 + 打包成docker镜像

.PHONY: docker
docker:
	@del /f /q webook 2>nul
	@go mod tidy
	@set GOOS=linux&& set GOARCH=amd64&& go build -tags=k8s -o webook .
	@docker rmi -f  dabiao22/webook:v0.0.1
	@docker build -t dabiao22/webook:v0.0.1 .