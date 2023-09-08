# simple-demo

## 抖音项目服务端

具体功能内容参考飞书说明文档

## 部署文档

您需要依次完成以下步骤

1. 完善配置文件
    ```shell
    mv settings.yml.template settings.yml
    ```
    > 依照配置文件模板的说明填写配置文件，你至少需要以下环境的支持：
    > - MySQL
    > - Redis
    > - rabbitMQ
2. MySQL建表
    > 根据 `config/init.sql` 文件中的建表语句建立项目所需要的表格
3. 初始化项目
    ```shell
        go run main.go init
    ```
   > `init`命令将初始化角色权限，同时利用`gorm gen`生成安全可靠的DAO层。
   > 
   > 所以，在你改动MySQL中的表结构之后，也需要重新执行 `init`命令

4. 启动服务器程序
    ```shell
        go run main.go server -m=release -c="./config/settings.yml"
    ```
   > 不同的模式会有不同的日志输出方式
