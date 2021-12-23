package unit

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"strings"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func Get(uri string, engine *gin.Engine) []byte {
	// 构造GET请求
	r := httptest.NewRequest("GET", uri, nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, r)

	// 获取响应
	result := w.Result()
	defer result.Body.Close()
	resp, _ := ioutil.ReadAll(result.Body)
	return resp
}

func PostJson(uri string, params map[string]interface{}, engine *gin.Engine) []byte {
	// 构造POST请求
	body, _ := json.Marshal(params)
	r := httptest.NewRequest("POST", uri, bytes.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, r)

	// 获取响应
	result := w.Result()
	defer result.Body.Close()
	resp, _ := ioutil.ReadAll(result.Body)
	return resp
}

func PostJsonWithHeaders(uri string, params map[string]interface{}, engine *gin.Engine, headers map[string]string) []byte {
	// 构造POST请求
	body, _ := json.Marshal(params)
	r := httptest.NewRequest("POST", uri, bytes.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		r.Header.Add(k, v)
	}
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, r)

	// 获取响应
	result := w.Result()
	defer result.Body.Close()
	resp, _ := ioutil.ReadAll(result.Body)
	return resp
}

func PostForm(uri string, params map[string]string, engine *gin.Engine) []byte {
	// 构造post form的参数
	var s []string
	for k, v := range params {
		s = append(s, k+"="+v)
	}
	body := strings.Join(s, "&")
	// 构造请求
	r := httptest.NewRequest("POST", uri, bytes.NewReader([]byte(body)))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, r)

	// 获取响应
	result := w.Result()
	defer result.Body.Close()
	resp, _ := ioutil.ReadAll(result.Body)
	return resp
}

func PostFile(uri string, params map[string]string, f io.Reader, engine *gin.Engine) []byte {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, err := writer.CreateFormFile("file", "filename")
	if err != nil {
		return nil
	}
	if _, err = io.Copy(file, f); err != nil {
		return nil
	}

	for k, v := range params {
		_ = writer.WriteField(k, v)
	}

	contentType := writer.FormDataContentType()
	if err = writer.Close(); err != nil {
		return nil
	}

	// 构造请求
	r := httptest.NewRequest("POST", uri, body)
	r.Header.Add("Content-Type", contentType)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, r)

	// 获取响应
	result := w.Result()
	defer result.Body.Close()
	resp, _ := ioutil.ReadAll(result.Body)
	return resp
}
