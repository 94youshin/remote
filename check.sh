#!/bin/bash
#
# 节点健康巡检
#

echo "开始巡检---------------------------------------------"
#1. 输出当前时间
echo "当前时间---------------------------------------------"
date -d "1970-01-01 1234567890 seconds" +"%Y-%m-%d %H:%M:%S"

#2. 查看管理网卡IP地址
echo "管理网卡信息-----------------------------------------"
ip addr show eth0
if [ $? -ne 0 ];then
    echo "查看管理网卡信息失败"
	exit 1
fi

#3. 检查mysql服务是否正常
echo "mysqld服务状态---------------------------------------"
systemctl status mysqld
if [ $? -ne 0 ];then
    echo "mysqld 服务状态异常."
	exit 1
fi

#4. 检查Docker服务状态
echo "docker服务状态--------------------------------------"
systemctl status docker
if [ $? -ne 0 ];then
    echo "docker 服务状态异常."
	exit 1
fi

#5. 查询异常容器列表
echo "已停止容器列表--------------------------------------"
docker ps -a | grep Exited

echo "巡检结束,目标服务状态均正常-------------------------"
