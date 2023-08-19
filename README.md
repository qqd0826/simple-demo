# simple-demo

## 原架构和运行

linux：

```bash
go build && ./simple-demo
```

windows系统默认终端：

```powershell
go build; .\simple-demo
```

### 功能说明

接口功能不完善，仅作为示例

* 用户登录数据保存在内存中，单次运行过程中有效
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可

### 测试

test 目录下为不同场景的功能测试case，可用于验证功能实现正确性

其中 common.go 中的 _serverAddr_ 为服务部署的地址，默认为本机地址，可以根据实际情况修改

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试

## 数据库

数据库配置与连接方法在db文件夹的**db.go**中，连接自己本地数据库请修改以下配置（默认端口3306）

```go
connect := "用户名:密码@tcp(localhost:3306)/数据库"
```

AutoMigrate会将结构体和数据库表相匹配，若不存在对应表则会新生成

其他controller用到其他表请在db.go中添加对应语句

数据库对应表名称：结构体名s

对应列名为common.go中对应结构体的属性的json后的值

```go
DB.AutoMigrate(&model.User{})
```

controller中导入db包并调用类似语句即可对数据库进行查询等功能

```go
db.DB.Where("id=?", 1).First(&user1)
```

