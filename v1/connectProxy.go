package vproxy

import (
	"net/http"
    //"fmt"
    "net"
)
var resultStatus200 = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")

type connectProxy struct{
    config      *Config
    transport   http.RoundTripper
}
func (cp *connectProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request){
    dial := net.Dial
    if tr, ok := cp.transport.(*http.Transport); ok && tr.Dial != nil {
        dial = tr.Dial
    }
    netConn, err := dial("tcp", req.Host)
    if err != nil {
		http.Error(rw, err.Error(), http.StatusBadGateway)
		return
    }

    hj, ok := rw.(http.Hijacker)
	if !ok {
		http.Error(rw, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

    conn.Write(resultStatus200)

    var bufSize int = defaultDataBufioSize
    if cp.config != nil {
        bufSize = cp.config.DataBufioSize
    }

    go copyDate(netConn, conn, bufSize)
    go copyDate(conn, netConn, bufSize)
}

