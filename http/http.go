package http

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/basicfu/gf/g"
	"github.com/basicfu/gf/gfile"
	"github.com/basicfu/gf/json"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type (
	H struct {
		Method  string
		Headers g.Map
		Params  g.Map
		//form json file raw 4选1
		Form               g.Map
		Json               any
		File               g.Map //value中允许有file类型
		Raw                string
		ContentType        string
		Cookies            g.Map
		Proxy              string        //代理,不带http前缀，ip+端口
		Timeout            time.Duration //超时时间
		TLSConfig          *tls.Config   //证书设置
		Auth               g.Map
		RandomUa           bool
		AllowRedirect      bool
		Chunked            bool
		DisableKeepAlives  bool
		DisableCompression bool
		SkipVerifyTLS      bool //跳过证书验证
		NoRetry            bool //关闭重试，默认开启重试
		RetryDelay         []time.Duration
	}
	Response struct {
		Success    bool   //2xx为成功，如有特殊需求，可自行根据状态码是否成功
		StatusCode int    //状态码，请求未成功为0
		ErrorMsg   string //错误文本
		Data       []byte //请求完就获取body，虽然影响性能，但是不用在使用此http时主动随时释放
		Header     *fasthttp.ResponseHeader
	}
)
type File struct {
	FileName string //文件名，可选
	Value    any    //值，必填
}

/*
*
默认要全部随机UA
*/
func (resp Response) String() string {
	return string(resp.Data)
}
func (resp Response) Reader() *bytes.Reader {
	return bytes.NewReader(resp.Data)
}
func (resp Response) Json() *json.Result {
	return json.Parse(string(resp.Data))
}
func (resp Response) GetHeader(key string) string {
	if resp.Header == nil {
		return ""
	}
	return string(resp.Header.Peek(key))
}
func (resp Response) AllCookie() string {
	var str []string
	resp.Header.VisitAllCookie(func(key, value []byte) {
		parts := strings.Split(string(value), ";")
		if len(parts) != 0 && parts[0] != "" {
			str = append(str, parts[0])
		}
	})
	return strings.Join(str, "; ")
}
func GetUrl(url string) Response {
	return Get(url, H{})
}
func Get(url string, h H) Response {
	h.Method = http.MethodGet
	return Do(url, h)
}
func Post(url string, h H) Response {
	h.Method = http.MethodPost
	return Do(url, h)
}
func Put(url string, h H) Response {
	h.Method = http.MethodPut
	return Do(url, h)
}
func Delete(url string, h H) Response {
	h.Method = http.MethodDelete
	return Do(url, h)
}
func Do(url string, h H) Response {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	if resp == nil {
		panic(resp)
	}
	defer func() { //离开此方法就会释放
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	req.SetRequestURI(url)
	setRequest(req, h)
	req.Header.SetMethod(h.Method)
	if h.Timeout == 0 {
		h.Timeout = 10 * time.Second //默认时间
	}
	c := &fasthttp.Client{
		TLSConfig: h.TLSConfig,
	}
	if h.SkipVerifyTLS {
		c.ConfigureClient = func(c *fasthttp.HostClient) error {
			c.TLSConfig = &tls.Config{InsecureSkipVerify: true}
			return nil
		}
	}
	if h.Proxy != "" {
		if strings.HasPrefix(h.Proxy, "socks5://") {
			c.Dial = fasthttpproxy.FasthttpSocksDialer(h.Proxy)
		} else {
			c.Dial = fasthttpproxy.FasthttpHTTPDialerTimeout(h.Proxy, h.Timeout)
		}
	}
	runN := 1
	if !h.NoRetry {
		if len(h.RetryDelay) == 0 { //没有给重试间隔使用默认3次重试
			h.RetryDelay = []time.Duration{
				1 * time.Second,
				1 * time.Second,
				2 * time.Second,
			}
		}
		runN += len(h.RetryDelay)
	}
	r := Response{}
	for i := 0; i < runN; i++ {
		if i != 0 {
			resp.Reset()
		}
		r = do(c, req, resp, h)
		if r.Success {
			break
		}
		if r.StatusCode != 0 && r.StatusCode != 429 && (r.StatusCode < 500 || r.StatusCode > 599) { //只有这些状态码才重试
			break
		}
		if i < runN-1 { //需要重试
			time.Sleep(h.RetryDelay[i])
		}
	}
	return r
}
func do(c *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response, h H) Response {
	data := []byte("")
	err := c.DoTimeout(req, resp, h.Timeout)
	if err != nil {
		//if errors.Is(err, fasthttp.ErrTimeout) {
		//	//request timeout
		//}
		//var netErr net.Error
		//if errors.As(err, &netErr) && netErr.Timeout() {
		//	//network timeout
		//}
		return Response{Success: false, StatusCode: 0, ErrorMsg: err.Error(), Data: data, Header: &resp.Header}
	}
	code := resp.StatusCode()
	if code < 200 || code >= 300 {
		return Response{Success: false, StatusCode: code, ErrorMsg: fasthttp.StatusMessage(code), Data: resp.Body(), Header: &resp.Header}
	}
	if strings.ToLower(string(resp.Header.Peek("content-encoding"))) == "gzip" {
		gunzip, _ := resp.BodyGunzip()
		data = gunzip
	} else {
		data = resp.Body()
	}
	if strings.Contains(strings.ToLower(string(resp.Header.Peek("content-type"))), "gbk") {
		reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
		data, _ = ioutil.ReadAll(reader)
	}
	return Response{Success: true, StatusCode: code, ErrorMsg: "", Data: data, Header: &resp.Header}
}
func setRequest(req *fasthttp.Request, h H) {
	if isConflict(h) {
		panic("请求体form json file raw只能有一个")
	}
	if h.Headers != nil {
		for k, v := range h.Headers {
			req.Header.Set(strings.ToLower(k), g.String(v))
		}
	}
	if h.RandomUa {
		//TODO 随机UA选择
		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
	}
	if h.Params != nil {
		args := req.URI().QueryArgs()
		for k, v := range h.Params {
			args.Add(k, g.String(v))
		}
	}
	if h.Cookies != nil {
		for k, v := range h.Cookies {
			req.Header.SetCookie(k, g.String(v))
		}
	}
	//请求体
	if h.Json != nil {
		req.Header.Set("Content-Type", "application/json")
		req.SetBodyString(json.String(h.Json))
		return
	}
	if h.Form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		var arr []string
		for k, v := range h.Form {
			arr = append(arr, fmt.Sprintf("%s=%s", k, url.QueryEscape(g.String(v))))
		}
		req.SetBodyString(strings.Join(arr, "&"))
		return
	}
	if h.Raw != "" {
		req.SetBodyString(h.Raw)
		if h.ContentType != "" {
			req.Header.Set("Content-Type", h.ContentType)
		}
		return
	}
	if h.File != nil {
		var b bytes.Buffer
		writer := multipart.NewWriter(&b)
		for k, v := range h.File {
			switch v.(type) {
			case File:
				f := v.(File)
				var reader io.Reader
				dataFileName := ""
				if fb, ok := f.Value.([]byte); ok {
					reader = bytes.NewReader(fb)
				} else {
					fv := g.String(f.Value)
					if strings.HasPrefix(fv, "http://") || strings.HasPrefix(fv, "https://") {
						imgResp := GetUrl(fv)
						reader = imgResp.Reader()
						u, err := url.Parse(fv)
						if err == nil {
							dataFileName = path.Base(u.Path)
						}
					} else if regexp.MustCompile(`^[A-Za-z0-9+/]+={0,2}$`).MatchString(fv) || strings.HasPrefix(fv, "data:") {
						d := fv
						if strings.HasPrefix(d, "data:") {
							arr := strings.Split(d, ",")
							if len(arr) == 2 {
								d = arr[1]
							}
						}
						imageData, err := base64.StdEncoding.DecodeString(d)
						if err != nil {
							panic(err)
						}
						reader = bytes.NewReader(imageData)
					} else {
						u, err := url.Parse(fv)
						if err == nil && u.Scheme != "" {
							panic("file字段中的value格式错误")
						}
						file, err := os.Open(fv)
						if err != nil {
							panic(err)
						}
						reader = file
						dataFileName = gfile.Basename(fv)
					}
				}
				filename := f.FileName
				if filename == "" {
					filename = dataFileName //如果还为空就空着
				}
				part, _ := writer.CreateFormFile(k, filename)
				_, _ = io.Copy(part, reader)
			default:
				_ = writer.WriteField(k, g.String(v))
			}
		}
		err := writer.Close()
		if err != nil {
			panic(err)
		}
		req.Header.SetContentType(writer.FormDataContentType())
		req.SetBody(b.Bytes())
	}
}
func isConflict(h H) bool {
	count := 0
	if h.Form != nil {
		count++
	}
	if h.Raw != "" {
		count++
	}
	//if h.File != nil {
	//	count++
	//}
	if h.Json != nil {
		count++
	}
	return count > 1
}
