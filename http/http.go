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
	"github.com/basicfu/gf/gconv"
	"github.com/basicfu/gf/json"
	"github.com/basicfu/gf/os/gfile"
	"github.com/basicfu/gf/text/gstr"
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
		SkipVerifyTLS      bool
	}
	Response struct {
		Success    bool   //网络级别，400，500为成功
		ErrorMsg   string //错误文本
		Data       []byte //请求完就获取body，虽然影响性能，但是不用在使用此http时主动随时释放
		StatusCode int    //状态码
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
	//TODO response里可设置详细错误信息，比如超时错误等
	c := &fasthttp.Client{}
	if h.SkipVerifyTLS {
		c.ConfigureClient = func(c *fasthttp.HostClient) error {
			c.TLSConfig = &tls.Config{InsecureSkipVerify: true}
			return nil
		}
	}
	if h.Proxy != "" {
		if gstr.HasPrefix(h.Proxy, "socks5://") {
			c.Dial = fasthttpproxy.FasthttpSocksDialer(h.Proxy)
		} else {
			c.Dial = fasthttpproxy.FasthttpHTTPDialerTimeout(h.Proxy, h.Timeout)
		}
	}
	c.TLSConfig = h.TLSConfig
	//TODO 错误重试
	//for i := 1; i <= 3; i++ {
	//	response, err = c.Client.Do(req) // nolint
	//	if err == nil {
	//		break
	//	}
	//	time.Sleep(time.Duration(i*100) * time.Millisecond)
	//}
	data := []byte("")
	if err := c.DoTimeout(req, resp, h.Timeout); err != nil { //分请求超时(如主机不通)和代理超时
		return Response{Success: false, ErrorMsg: err.Error(), Data: data, Header: &resp.Header, StatusCode: fasthttp.StatusGatewayTimeout}
	}
	if string(resp.Header.Peek("content-encoding")) == "gzip" { //是否忽略大小写
		gunzip, _ := resp.BodyGunzip()
		data = gunzip
	} else {
		data = resp.Body()
	}
	if strings.Contains(strings.ToLower(string(resp.Header.Peek("content-type"))), "gbk") {
		reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
		data, _ = ioutil.ReadAll(reader)
	}
	return Response{Success: true, ErrorMsg: "", Data: data, Header: &resp.Header, StatusCode: resp.StatusCode()}
}
func setRequest(req *fasthttp.Request, h H) {
	if isConflict(h) {
		panic("请求体form json file raw只能有一个")
	}
	if h.Headers != nil {
		for k, v := range h.Headers {
			req.Header.Set(strings.ToLower(k), gconv.String(v))
		}
	}
	if h.RandomUa {
		//TODO 随机UA选择
		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
	}
	if h.Params != nil {
		args := req.URI().QueryArgs()
		for k, v := range h.Params {
			args.Add(k, gconv.String(v))
		}
	}
	if h.Cookies != nil {
		for k, v := range h.Cookies {
			req.Header.SetCookie(k, gconv.String(v))
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
			arr = append(arr, fmt.Sprintf("%s=%s", k, url.QueryEscape(gconv.String(v))))
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
					fv := gconv.String(f.Value)
					if strings.HasPrefix(fv, "http://") || strings.HasPrefix(fv, "https://") {
						imgResp := GetUrl(fv)
						reader = imgResp.Reader()
						u, err := url.Parse(fv)
						if err == nil {
							dataFileName = path.Base(u.Path)
						}
					} else if regexp.MustCompile(`^[A-Za-z0-9+/]+={0,2}$`).MatchString(fv) || gstr.HasPrefix(fv, "data:") {
						d := fv
						if gstr.HasPrefix(d, "data:") {
							arr := gstr.Split(d, ",")
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
				_ = writer.WriteField(k, gconv.String(v))
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
