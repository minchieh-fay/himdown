#!/bin/bash
# 华创密信下载页面Docker镜像构建和运行脚本

# 设置变量
IMAGE_NAME="himdown"
CONTAINER_NAME="himdown-container"
PORT=80

# 停止并删除已有容器（如果存在）
echo "停止并删除已有容器（如果存在）"
docker stop $CONTAINER_NAME >/dev/null 2>&1
docker rm $CONTAINER_NAME >/dev/null 2>&1

# 构建镜像
echo "开始构建Docker镜像..."
docker build -t $IMAGE_NAME .

# 检查构建结果
if [ $? -ne 0 ]; then
    echo "镜像构建失败！"
    exit 1
fi

# 运行容器
echo "启动容器..."
docker run -d --name $CONTAINER_NAME -p $PORT:80 $IMAGE_NAME

# 检查运行结果
if [ $? -ne 0 ]; then
    echo "容器启动失败！"
    exit 1
fi

# 显示容器信息
echo "容器已成功启动"
echo "访问地址: http://localhost:$PORT"
echo ""
echo "容器信息:"
docker ps | grep $CONTAINER_NAME

echo ""
echo "查看日志:"
echo "docker logs $CONTAINER_NAME"

echo ""
echo "停止容器:"
echo "docker stop $CONTAINER_NAME" 