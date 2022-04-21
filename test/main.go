package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

var (
	URL  = "http://127.0.0.1:8080"
	addr = "V6Avp9KqLfwGUaRFPV27a8VZhwxiYofoU"
	msg  = "123456"
)

// 响应结果
type Response struct {
	Code int    `json:"code"` // 错误码
	Msg  string `json:"msg"`  // 信息提示
	Data []byte `json:"data"` // 返回数据
}

// 集成测试
func main() {
	fmt.Println("==========  测试: /ping ========== ")
	pingResult, err := request(URL+"/ping", "GET", nil)
	if err != nil {
		fmt.Println("ping:", err)
		return
	}
	fmt.Println(pingResult)

	fmt.Println("==========  测试: /create ========== ")
	createResult, err := request(URL+"/create", "GET", nil)
	if err != nil {
		fmt.Println("create:", err)
		return
	}
	fmt.Println(createResult)

	fmt.Println("==========  测试: /sign ========== ")
	signArgs := map[string]interface{}{
		"address": addr,
		"msg":     []byte(msg),
	}
	signResult, err := request(URL+"/sign", "POST", signArgs)
	if err != nil {
		fmt.Println("create:", err)
		return
	}
	fmt.Println(signResult)

	fmt.Println("========== 测试: /verify ==========")
	verifyArgs := map[string]interface{}{
		"address": addr,
		"sign":    signResult.Data,
		"msg":     msg,
	}
	verifyResult, err := request(URL+"/verify", "POST", verifyArgs)
	if err != nil {
		fmt.Println("create:", err)
		return
	}
	fmt.Println(string(verifyResult.Data))
}

// 请求服务封装
func request(url, method string, args map[string]interface{}) (*Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源
	// 默认是application/x-www-form-urlencoded
	req.Header.SetContentType("application/json")
	req.Header.SetMethod(method)

	req.SetRequestURI(url)

	// 请求体
	if args != nil {
		requestBody, err := json.Marshal(&args)
		if err != nil {
			return nil, err
		}
		req.SetBody(requestBody)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, err
	}

	// 处理返回结果
	res := &Response{}
	err := json.Unmarshal(resp.Body(), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
