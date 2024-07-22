# proxy
autodl 单端口，转发多端口

# 可自行源码编译生成，linux 系统 
GOOS=linux GOARCH=amd64 go build -o proxy_in_instance


# 运行
```
  # 如果没有添加可执行权限，先添加权限
  chmod +x proxy_in_instance
  
  # 启动
  ./proxy_in_instance
```


# 使用：
config.yaml 配置转发。
通过 pattern 后缀进行转发
https://u219252-b96d-cb24cdf.west.seeacloud.com:8443/ollama  转发到配置的端口 http://127.0.0.1:11434

如原来配置的 ollama 的地址是
https://u219252-b96d-cb24cdf.west.seeacloud.com:8443
通过 proxy 转发后地址需要增加后缀如
https://u219252-b96d-cb24cdf.west.seeacloud.com:8443/ollama
