# IM-System
#### 1 基于Golang的即时通信系统



#### 2 服务器读写模型

Server服务器采用读写分离模型
从user读信息位置：Server.Handler() Go程中
往user写信息位置：User.ListenMessage() Go程中



#### 3 支持的用户消息格式

##### 3.1 who

`who`：查询当前在线用户

##### 3.2 rename|NEWNAME

`rename|foo`：将用户名改为foo

##### 3.3 to|NAME|CONTENT

`to|foo|Hello`：对foo说Hello

##### 3.4 *

默认将消息广播给所有在线用户（包括自己）
