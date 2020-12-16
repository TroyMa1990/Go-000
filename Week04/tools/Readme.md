## 使用
* 1. 检测接口  
```
cd tools
sh run.sh
```
* 2. 测试接口
```
cd tools/curl
sh showStatus.curl
etc...
```

## 接口调用样例：  

### 1. 查看指令(指定设备)：

**请求URL：** 
- `echo`
  
**请求方式：** 
- POST

**参数：** 
- 无

**请求示例** 

```
curl -s -d '
{
"DevId":"301"
}
' "http://127.0.0.1:8880/echo"
```

**返回示例** 

```
{
    "errcode":0,
    "msg":"Success",
    "ret":0
}
```



## 接口及控制内容:

```
/echo  检测系统状态

TODO
...
```


