# Gin Template project

此项目是一个Gin框架的模板项目,用于快速创建golang后端项目,项目内部预设了基础的身份验证功能,如`登录` 和 `创建账号`
此项目原生支持多语言,并且内置了系统错误,具有完善的错误处理。

## Usage

1. Config

此项目所有配置项都由`.env`文件开始,默认的配置项如下
```.env
config="data/config.json"
mode="development"
gin_mode="debug"
port=8080
app_name="Template"
env_mode=development
logger_lang="en"
logger_path="logs"
logger_name="{app}-{level}-{date}.log"
```
- **config** 指定了系统配置的json文件的路径
- **mode** 指定了当前的开发模式
- **gin_mode** 指定了gin的模式,默认debug
- **port** 服务监听端口
- **app_name** 服务的名称,由于当前是模板项目,所以为`Template`
- **logger_lang** 指定的默认日志语言, 所有的可支持的语言都可以在 `locale` 文件夹下查看到,每一个json就是详细的语言配置
- **logger_path** 日志保存的路径
- **logger_name** 日志文件的模板字符串 其中的 `app`为当前的 项目名称(`app_name`) `date` 为项目启动日期


2. Struct

项目的所有启动配置都从boot文件夹开始,比如初始化数据库,开始创建Gin路由等.整个模板采用MVC的理念进行分割代码.

```
├─api  # 所有的api存放位置
├─boot # 项目的所有引导配置
├─common # 一些工具函数
├─core # gin框架的一些封装
├─dao # 数据库操作层
├─data # 系统配置和其他配置存放位置
├─dto # 数据交换层,包括数据字段的校验
├─global # 一些需要全局访问的变量的初始化
├─i18n # 多语言的支持,以及语言配置的定义
├─internal # 一些内置的常量和方法,比如自定义的错误类型 和内置的错误常量
├─locale # 多语言的配置
├─logs # 日志存放位置
├─middleware # gin中间件
├─model # 数据模型的定义用于ORM的操作
├─router # 全局路由的注册
├─service # 服务层,隔离DAO层和API的操作,将内置的错误统一转换为自定义的错误
```
