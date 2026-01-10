# 将指定文件打包成镜像

# 指定基础镜像，我这个镜像即在这个镜像上添加的功能
FROM ubuntu:20.04

# 将编译后的文件，打包成的镜像，放到这个工作目录 /app 下
COPY webook /app/webook
WORKDIR /app

ENTRYPOINT ["/app/webook"]
