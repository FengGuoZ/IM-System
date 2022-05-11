# IM-System
基于Golang的即时通信系统
Server服务器采用读写分离模型
从user读信息位置：Server.Handler() Go程中
往user写信息位置：User.ListenMessage() Go程中
