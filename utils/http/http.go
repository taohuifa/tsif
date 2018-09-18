package http

import (
	"bytes"
	"fmt"
	// "log"      //日志库

	"encoding/json"
	"net/http" // http
)

// 获取Http参数
func GetParamByGet(r *http.Request, key string, defaultvalue interface{}) interface{} {
	querys := r.URL.Query()
	values, ok := querys[key]
	if !ok {
		return defaultvalue
	}
	if len(values) <= 0 {
		return defaultvalue
	}
	return values[0]
}

// 获取Http参数
func GetParam(r *http.Request, key string, defaultvalue interface{}) interface{} {
	return GetParamByGet(r, key, defaultvalue)
}

// 解析参数
func GetParams(r *http.Request) map[string]interface{} {
	//解析参数
	r.ParseForm() //解析参数，默认是不会解析的
	// 常规解析方法
	var params = make(map[string]interface{})
	values := r.Form
	var vlen = len(values)
	if vlen > 0 {
		for k, v := range values {
			vsize := len(v)
			if vsize > 0 {
				params[k] = v[0]
			}
		}
	}
	// log.Println("GetHttpParams: ", r, params)
	return params
}

// map转http参数
func CreateParams(params map[string]interface{}) string {
	b := bytes.Buffer{}
	var index int = 0
	for k, v := range params {
		if index > 0 {
			b.WriteString("&")
		}
		index++
		//写入参数
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(fmt.Sprint(v))

	}
	return b.String()
}

func Response(data []byte, w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json;charset=utf-8")
	w.Write(data)
}

func ResponseStr(msg string, w http.ResponseWriter) {
	Response([]byte(msg), w)
}

func ResponseJson(obj interface{}, w http.ResponseWriter) {
	jdata, err := json.Marshal(obj)
	if err != nil {
		ResponseStr("json error: "+err.Error(), w)
		return
	}
	Response(jdata, w)
}

func ResponseResult(code int, msg string, w http.ResponseWriter) {
	retStr := fmt.Sprintf("{\"code\":%d,\"message\":\"%s\"}", code, msg)
	Response([]byte(retStr), w)
}
