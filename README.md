# vproxy [![Build Status](https://travis-ci.org/456vv/vproxy.svg?branch=master)](https://travis-ci.org/456vv/vproxy)
go/golang HTTP/HTTPS proxy server, HTTP/HTTPS 代理服务器
<br />
最近更新20160821：<a href="/v1/update.txt">update.txt</a>
<br/>
已编译好的二进制文件下载：暂不提供
<br />
列表：
-----------------------------------
    const defaultDataBufioSize    = 1<<20                                            // 默认数据缓冲1MB
    type Config struct {                                                     // 配置
        DataBufioSize int                                                            // 缓冲区大小
    }
    type Proxy struct {                                                      // 代理
        Addr        string                                                           // 代理IP地址
        Server      http.Server                                                      // 服务器
        Transport   http.RoundTripper                                                // 代理
        *Config                                                                      // 配置
        l           net.Listener                                                     // 连接对象
    }
        func (p *Proxy) setDefault()                                                 // 设置默认
        func (p *Proxy) initServer() *http.Server                                    // 初始化服务器
        func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request)         // 处理
        func (p *Proxy) ListenAndServ() error                                        // 监听
        func (p *Proxy) Serve(l net.Listener) error                                  // 监听
        func (p *Proxy) Close() error                                                // 关闭代理

