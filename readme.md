# eat-and-go

你是否还在为中午吃什么而烦恼不已？不妨问问神奇的小海螺！
### 简介

`eat-and-go`是基于gin-gonic框架的后端服务，服务于slack APP（slack机器人），处理来自slack APP的请求，给用户随机推荐一家餐厅，俗称“今天中午吃什么”。目前该服务支持：

1. 关键字检索，检索`店铺名称`、`标签1`和`推荐菜`中包含该关键字的所有餐厅，随机返回一家餐厅的详情链接；
2. 如果没有关键字，或没有餐厅符合条件，随机返回一家餐厅的详情链接；
3. 每次推荐后，会推送问卷调查，收集用户满意度。

需要注意的是，若想测试该服务，最佳的做法是将其部署到服务器，通过slack测试。

### Swagger user guide
项目使用swagger管理API
* 安装 `swagger`
```
   $ sudo mkdir -p $GOPATH/src/github.com/swaggo
   $ cd $GOPATH/src/github.com/swaggo
   $ git clone https://github.com/swaggo/swag
   $ cd swag/cmd/swag/
   $ go install -v
```
如果不知道gopath在哪里，或者gopath不在`PATH`内，可以运行：
```
   $ go env
```
该指令会显示gopath
* 下载 `gin-swagger`
```
   $ cd $GOPATH/src/github.com/swaggo
   $ git clone https://github.com/swaggo/gin-swagger
```
* 生成swagger文档
```
   $ cd xxx/pathto/gopa-server/
   $ swag init
```
* API注释示例
```
   // @Summary       api
   // @Description   Add a new user
   // @Tags          user
   // @Accept        json
   // @Produce       json
   // @Param         env   path     string          true "dev/fat/uat/pro"
   // @Param         user  body     model.UserInfo  true "Create a new user"
   // @Success       200   {object} handler.Response 
   // @Router        /user/{env} [post]
   func Create(c *gin.Context) {
       ...
   }
```


详情见[swagger文档](https://github.com/swaggo/swag/blob/master/README.md)

### How to start


