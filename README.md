# vproxy [![Build Status](https://travis-ci.org/456vv/vproxy.svg?branch=master)](https://travis-ci.org/456vv/vproxy)
go/golang HTTP/HTTPS proxy server, HTTP/HTTPS 代理服务器
<br />
最近更新20160920：<a href="/v1/update.txt">update.txt</a>
<br/>
已编译好的二进制文件下载：<a href="https://github.com/456vv/vproxy/raw/master/v1/test/bin/vproxy%20Alpha.zip">vproxy Alpha.zip</a>
<br />
命令行：
-----------------------------------
    命令行例子：vproxy -addr 0.0.0.0:8080
    -addr string
          代理服务器地 (format "0.0.0.0:8080")
    -proxy string
          代理服务器的上级代理IP地址 (format "11.22.33.44:8888")
    -dataBufioSize int
          代理数据交换缓冲区大小，单位字节 (default 10240)
    -disableCompression
          禁止传送数据时候进行压缩 (default false)
    -disableKeepAlives
          禁止长连接 (default false)
    -expectContinueTimeout int
          http1.1过度到http2的等待超时，单位毫秒 (default 1000)
    -idleConnTimeout int
          空闲连接超时时，单位毫秒 (default 0)
    -maxIdleConns int
          保持空闲连接(TCP)数量 (default 500)
    -maxIdleConnsPerHost int
          保持空闲连接(Host)数量 (default 500)
    -responseHeaderTimeout int
          读取服务器发来的文件标头超时，单位毫秒 (default 0)
    -tlsHandshakeTimeout int
          SSL握手超时，单位毫秒 (default 10000)
    -log string
          日志文件(默认留空在控制台显示日志)  (format "./vproxy.txt")
    -logLevel int
          日志级别，0)不记录 1)客户端IP 2)认证 3)访问的Host地址 4)路径 5)请求 6)响应 7)错误 (default 0)
    -timeout int
          转发连接请求超时，单位毫秒 (default 5000)
    -user string
          用户名
    -pwd string
          密码

<br />
列表：
-----------------------------------
    type LogLevel int                                                                // 日志级别
    const (
        OriginAddr LogLevel    = iota+1                                              // 客户端。
        Authenticate                                                                 // 认证
        Host                                                                         // 访问的Host地址
        URI                                                                          // 路径
        Request                                                                      // 请求
        Response                                                                     // 响应
        Error                                                                        // 错误
    )
    type Config struct {                                                     // 配置
        DataBufioSize   int                                                          // 缓冲区大小
        Auth            func(username, password string) bool                         // 认证
        Timeout         time.Duration                                                // 转发连接请求超时
        Deadline        time.Time                                                    // 转发连接请求超时
    }
    type Proxy struct {                                                      // 代理
        *Config                                                                      // 配置
        Addr        string                                                           // 代理IP地址
        Server      http.Server                                                      // 服务器
        Transport   http.RoundTripper                                                // 代理
        ErrorLogLevel LogLevel                                                       // 日志级别
        ErrorLog    *log.Logger                                                      // 日志
    }
        func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request)         // 处理
        func (p *Proxy) ListenAndServ() error                                        // 监听
        func (p *Proxy) Serve(l net.Listener) error                                  // 监听
        func (p *Proxy) Close() error                                                // 关闭代理

