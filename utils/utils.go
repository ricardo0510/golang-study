package utils

import (
	"encoding/json"
	"net/http"
)

// 公共返回数据方法
func HandleResponse(w http.ResponseWriter, code int, data interface{}, msg string) {
	res := map[string]interface{}{
		"code": code,
		"data": data,
		"msg":  msg,
	}
	newRes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(newRes)
}
