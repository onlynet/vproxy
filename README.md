# vproxy
golang 代理服务器，go Proxy Server
<br />
列表：
-----------------------------------
	const defaultDataBufioSize    = 1<<20                             		// 默认数据缓冲1MB
	type Config struct {												// 配置
	    DataBufioSize int                                                   // 缓冲区大小
	}
	type Proxy struct {													// 代理
	    Addr        string                                                  // 代理IP地址
	    Server      http.Server                                         	// 服务器
	    Transport   http.RoundTripper                                   	// 代理
	    *Config                                                         	// 配置
	    l           net.Listener                                        	// 连接对象
	}
		func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request)// 处理
		func (p *Proxy) ListenAndServ() error 								// 监听
		func (p *Proxy) Serve(l net.Listener) error							// 监听
		func (p *Proxy) Close() error										// 关闭代理
	
