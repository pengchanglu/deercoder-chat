#!/usr/bin/env bash

# 开发模式 dev/prod
# 此处修改模式
# 执行该脚本
# dev模式
# 前后端api接口地址
devMode=dev
api=localhost:8006
httpProtocol=http
wsProtocol=ws

# prod模式
#devMode=prod
#api=chat.deercoder.com
#httpProtocol=https
#wsProtocol=wss

# 后端配置文件地址
# 修改各个模块下app.conf文件开发模式
confFiles=(api/conf/app.yaml chat-srv/conf/app.yaml user-srv/conf/app.yaml front-srv/conf/app.yaml)
# 前端配置文件地址
apiFile=front-srv/static/js/chat.js

# 后端conf配置修改
for conf in ${confFiles[@]}
#也可以写成for element in ${array[*]}
do
echo "配置文件: "${conf}
# 替换源文件第二行内容
sed -i '3c \  \devMode: '${devMode} "${conf}"
done

# 前端api修改
sed -i '2c var api = "'${api}'";' ${apiFile}
sed -i '3c var httpProtocol = "'${httpProtocol}'";' ${apiFile}
sed -i '4c var wsProtocol = "'${wsProtocol}'";' ${apiFile}