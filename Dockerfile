FROM nginx:alpine

# 设置工作目录
WORKDIR /usr/share/nginx/html

# 拷贝HTML文件
COPY index.html .
COPY versions.json .

# 拷贝图片文件
COPY images/ ./images/

# 创建下载目录
RUN mkdir -p downloads

# 如果有实际的下载包，可以取消下面注释并添加
COPY downloads/ ./downloads/

# 拷贝自定义的nginx配置
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

# 暴露80端口
EXPOSE 80

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s CMD wget --quiet --tries=1 --spider http://localhost/ || exit 1

# 使用nginx的默认启动命令
CMD ["nginx", "-g", "daemon off;"] 