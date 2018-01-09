package main

import (
	"github.com/elazarl/goproxy"
	"net/http"
	"log"
	"io/ioutil"
	"bytes"
	"wangzhe/lib"
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

var (
	proxy *goproxy.ProxyHttpServer
)

func init() {
	proxy = goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.Verbose = false

	//请求拦截
	requestHandle := func(request *http.Request, ctx *goproxy.ProxyCtx) (req *http.Request, resp *http.Response) {
		req = request

		if ctx.Req.URL.Path == `/question/fight/findQuiz` || ctx.Req.URL.Path == `/question/fight/choose` {
			requestBody, e := ioutil.ReadAll(req.Body)
			if check(e) {
				req.Body = ioutil.NopCloser(bytes.NewReader(lib.Injection(requestBody, ctx)))
			}
		}
		return
	}
	//返回拦截
	responseHandle := func(response *http.Response, ctx *goproxy.ProxyCtx) (resp *http.Response) {
		resp = response

		if ctx.Req.URL.Path == `/question/fight/findQuiz` || ctx.Req.URL.Path == `/question/fight/choose` {
			responseBody, e := ioutil.ReadAll(resp.Body)
			if check(e) {
				resp.Body = ioutil.NopCloser(bytes.NewReader(lib.Injection(responseBody, ctx)))
			}
		}
		return
	}

	proxy.OnRequest().DoFunc(requestHandle)
	proxy.OnResponse().DoFunc(responseHandle)

}

func check(e error) (pass bool) {
	if e != nil {
		log.Println(e)
		return false
	}
	return true
}

func main() {
	go Run("8989")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println("exiting")
}

func Run(port string) {

	go func() {
		log.Println("server port:", port)
		e := http.ListenAndServe(":"+port, proxy)
		check(e)
	}()

	go func() {
		crtSever := http.NewServeMux()
		crtSever.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Disposition", "attachment; filename=ca.crt")
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(goproxy.CA_CERT)
		})
		e := http.ListenAndServe(":8080", crtSever)
		check(e)
	}()
}
