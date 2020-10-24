## 一个微型的web框架
### features:
#### 1,路由
将请求转发到对应的处理函数
支持动态路由,例如:  
/hello/:name/  
/static/*filepath

#### 2,context
对原生的request和response进行封装
#### 3,分组控制
对url按照规则分组(url或者key),方便插入中间件
#### 4,中间件
框架不可能去理解所有的业务,允许用户定义自己的功能
#### 5,错误处理
能够对错误recover