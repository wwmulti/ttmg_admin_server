package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RequestOption 请求配置选项
type RequestOption struct {
	Method      string            // 请求方法 GET/POST/PUT/DELETE等
	URL         string            // 请求URL
	Headers     map[string]string // 请求头
	QueryParams map[string]string // URL查询参数
	Body        interface{}       // 请求体，可以是string、[]byte、map、struct等
	JSON        interface{}       // JSON请求体（自动设置Content-Type为application/json）
	FormData    map[string]string // Form表单数据（自动设置Content-Type为application/x-www-form-urlencoded）
	Timeout     time.Duration     // 请求超时时间
	Cookies     []*http.Cookie    // Cookies
	Proxy       string            // 代理地址
	UserAgent   string            // User-Agent
	BasicAuth   *BasicAuth        // Basic认证
	RetryTimes  int               // 重试次数
	RetryDelay  time.Duration     // 重试延迟
}

// BasicAuth Basic认证信息
type BasicAuth struct {
	Username string
	Password string
}

// Response 响应结构体
type Response struct {
	StatusCode int            // HTTP状态码
	Body       string         // 响应体
	Header     http.Header    // 响应头
	Cookies    []*http.Cookie // Cookies
	RequestURL string         // 实际请求的URL（可能包含重定向）
	Duration   time.Duration  // 请求耗时
	Error      error          // 错误信息
}

// 默认请求配置
var defaultOptions = RequestOption{
	Method:     "GET",
	Timeout:    30 * time.Second,
	RetryTimes: 0,
	RetryDelay: 1 * time.Second,
	UserAgent:  "Go-Curl/1.0",
}

// NewRequest 创建新请求（最灵活的方式）
func NewRequest(opts RequestOption) *Response {
	startTime := time.Now()

	// 合并默认配置
	opts = mergeOptions(opts)

	var resp *Response

	// 支持重试机制
	for i := 0; i <= opts.RetryTimes; i++ {
		if i > 0 {
			time.Sleep(opts.RetryDelay)
		}

		resp = executeRequest(opts)
		if resp.Error == nil && resp.StatusCode < 500 {
			break
		}
	}

	if resp != nil {
		resp.Duration = time.Since(startTime)
	}

	return resp
}

// 执行请求
func executeRequest(opts RequestOption) *Response {
	// 准备请求体和Content-Type
	bodyReader, contentType := prepareRequestBody(opts)

	// 构建完整的URL（包含查询参数）
	fullURL := buildURL(opts.URL, opts.QueryParams)

	// 创建HTTP请求
	req, err := http.NewRequest(opts.Method, fullURL, bodyReader)
	if err != nil {
		return &Response{Error: fmt.Errorf("创建请求失败: %v", err)}
	}

	// 设置请求头
	setRequestHeaders(req, opts, contentType)

	// 设置Basic认证
	if opts.BasicAuth != nil {
		req.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}

	// 设置Cookies
	for _, cookie := range opts.Cookies {
		req.AddCookie(cookie)
	}

	// 创建HTTP客户端
	client := createHTTPClient(opts)

	// 发送请求
	httpResp, err := client.Do(req)
	if err != nil {
		return &Response{Error: fmt.Errorf("请求失败: %v", err)}
	}
	defer httpResp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return &Response{Error: fmt.Errorf("读取响应失败: %v", err)}
	}

	return &Response{
		StatusCode: httpResp.StatusCode,
		Body:       string(bodyBytes),
		Header:     httpResp.Header,
		Cookies:    httpResp.Cookies(),
		RequestURL: httpResp.Request.URL.String(),
	}
}

// 准备请求体
func prepareRequestBody(opts RequestOption) (io.Reader, string) {
	var body io.Reader
	contentType := ""

	// 优先级：JSON > FormData > Body
	if opts.JSON != nil {
		// JSON请求
		jsonData, err := json.Marshal(opts.JSON)
		if err != nil {
			return nil, ""
		}
		body = bytes.NewBuffer(jsonData)
		contentType = "application/json"

	} else if len(opts.FormData) > 0 {
		// Form表单请求
		formData := url.Values{}
		for key, value := range opts.FormData {
			formData.Set(key, value)
		}
		body = strings.NewReader(formData.Encode())
		contentType = "application/x-www-form-urlencoded"

	} else if opts.Body != nil {
		// 原始Body请求
		switch v := opts.Body.(type) {
		case string:
			body = strings.NewReader(v)
		case []byte:
			body = bytes.NewBuffer(v)
		default:
			// 尝试转换为JSON
			jsonData, err := json.Marshal(v)
			if err == nil {
				body = bytes.NewBuffer(jsonData)
				contentType = "application/json"
			}
		}
	}

	return body, contentType
}

// 构建完整URL
func buildURL(baseURL string, queryParams map[string]string) string {
	if len(queryParams) == 0 {
		return baseURL
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

// 设置请求头
func setRequestHeaders(req *http.Request, opts RequestOption, contentType string) {
	// 设置User-Agent
	if opts.UserAgent != "" {
		req.Header.Set("User-Agent", opts.UserAgent)
	}

	// 设置Content-Type
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// 设置自定义请求头
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}
}

// 创建HTTP客户端
func createHTTPClient(opts RequestOption) *http.Client {
	client := &http.Client{
		Timeout: opts.Timeout,
	}

	// 设置代理
	// opts.Proxy = "socks5://192.168.1.99:8080"
	if opts.Proxy != "" {
		proxyURL, err := url.Parse(opts.Proxy)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
		}
	}

	return client
}

// 合并配置
func mergeOptions(opts RequestOption) RequestOption {
	if opts.Method == "" {
		opts.Method = defaultOptions.Method
	}
	if opts.Timeout == 0 {
		opts.Timeout = defaultOptions.Timeout
	}
	if opts.UserAgent == "" {
		opts.UserAgent = defaultOptions.UserAgent
	}

	return opts
}

// ==================== 便捷方法 ====================

// Get 发送GET请求
func Get(url string, options ...RequestOption) *Response {
	opts := RequestOption{
		Method: "GET",
		URL:    url,
	}

	if len(options) > 0 {
		opts.Headers = options[0].Headers
		opts.QueryParams = options[0].QueryParams
		opts.Timeout = options[0].Timeout
		opts.Cookies = options[0].Cookies
		opts.UserAgent = options[0].UserAgent
		opts.BasicAuth = options[0].BasicAuth
	}

	return NewRequest(opts)
}

// Post 发送POST请求
func Post(url string, body interface{}, options ...RequestOption) *Response {
	opts := RequestOption{
		Method: "POST",
		URL:    url,
		Body:   body,
	}

	if len(options) > 0 {
		opts.Headers = options[0].Headers
		opts.Timeout = options[0].Timeout
		opts.Cookies = options[0].Cookies
		opts.UserAgent = options[0].UserAgent
		opts.BasicAuth = options[0].BasicAuth
	}

	return NewRequest(opts)
}

// PostJSON 发送JSON POST请求
func PostJSON(url string, jsonData interface{}, options ...RequestOption) *Response {
	opts := RequestOption{
		Method: "POST",
		URL:    url,
		JSON:   jsonData,
	}

	if len(options) > 0 {
		opts.Headers = options[0].Headers
		opts.Timeout = options[0].Timeout
		opts.Cookies = options[0].Cookies
		opts.UserAgent = options[0].UserAgent
		opts.BasicAuth = options[0].BasicAuth
		opts.Proxy = options[0].Proxy
	}

	return NewRequest(opts)
}

// PostForm 发送Form表单请求
func PostForm(url string, formData map[string]string, options ...RequestOption) *Response {
	opts := RequestOption{
		Method:   "POST",
		URL:      url,
		FormData: formData,
	}

	if len(options) > 0 {
		opts.Headers = options[0].Headers
		opts.Timeout = options[0].Timeout
		opts.Cookies = options[0].Cookies
		opts.UserAgent = options[0].UserAgent
		opts.BasicAuth = options[0].BasicAuth
	}

	return NewRequest(opts)
}

// Put 发送PUT请求
func Put(url string, body interface{}, options ...RequestOption) *Response {
	opts := RequestOption{
		Method: "PUT",
		URL:    url,
		Body:   body,
	}

	if len(options) > 0 {
		opts.Headers = options[0].Headers
		opts.Timeout = options[0].Timeout
		opts.Cookies = options[0].Cookies
		opts.UserAgent = options[0].UserAgent
		opts.BasicAuth = options[0].BasicAuth
	}

	return NewRequest(opts)
}

// Delete 发送DELETE请求
func Delete(url string, options ...RequestOption) *Response {
	opts := RequestOption{
		Method: "DELETE",
		URL:    url,
	}

	if len(options) > 0 {
		opts.Headers = options[0].Headers
		opts.Body = options[0].Body
		opts.Timeout = options[0].Timeout
		opts.Cookies = options[0].Cookies
		opts.UserAgent = options[0].UserAgent
		opts.BasicAuth = options[0].BasicAuth
	}

	return NewRequest(opts)
}

// ==================== 响应方法 ====================

// JSON 将响应体解析为JSON
func (r *Response) JSON(v interface{}) error {
	if r.Error != nil {
		return r.Error
	}
	return json.Unmarshal([]byte(r.Body), v)
}

// Bytes 获取响应体的字节数组
func (r *Response) Bytes() []byte {
	if r.Error != nil {
		return nil
	}
	return []byte(r.Body)
}

// OK 判断请求是否成功（状态码2xx）
func (r *Response) OK() bool {
	return r.Error == nil && r.StatusCode >= 200 && r.StatusCode < 300
}

// String 获取响应体字符串
func (r *Response) String() string {
	if r.Error != nil {
		return r.Error.Error()
	}
	return r.Body
}

// ==================== 快捷使用示例 ====================
// func quickExamples() {
// 	// 最简单的GET请求
// 	resp := curl.Get("https://api.example.com/data")
// 	if resp.OK() {
// 		log.Println(resp.Body)
// 	}

// 	// 带查询参数的GET请求
// 	resp = curl.Get("https://api.example.com/users",
// 		curl.RequestOption{
// 			QueryParams: map[string]string{
// 				"page":  "1",
// 				"limit": "20",
// 				"sort":  "name",
// 			},
// 		})

// 	// POST JSON数据
// 	resp = curl.PostJSON("https://api.example.com/users",
// 		map[string]interface{}{
// 			"name":  "Alice",
// 			"email": "alice@example.com",
// 		},
// 		curl.RequestOption{
// 			Headers: map[string]string{
// 				"Authorization": "Bearer your-token",
// 			},
// 		})

// 	// 获取并解析JSON响应
// 	var userData struct {
// 		ID    int    `json:"id"`
// 		Name  string `json:"name"`
// 		Email string `json:"email"`
// 	}

// 	if resp.OK() {
// 		resp.JSON(&userData)
// 		fmt.Printf("用户ID: %d, 姓名: %s\n", userData.ID, userData.Name)
// 	}
// }
