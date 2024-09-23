package httpUtil

import "net/http"

// requestParamsDto 请求相关参数
type requestParamsDto struct {
	method      string         // 请求方式
	Path        string         `json:"path"` // 接口路径
	Params      map[string]any // 请求参数
	timeout     uint32         // 接口超时时间，默认3000毫秒
	header      http.Header    // 请求header头
	contentType string
	basicAuth   *BasicAuth
}

type BasicAuth struct {
	Username string
	Password string
}

type ReqParamsOption func(*requestParamsDto)

func WithTimeout(millionSecond uint32) ReqParamsOption {
	return func(params *requestParamsDto) {
		params.timeout = millionSecond
	}
}
func WithReqHeader(header http.Header) ReqParamsOption {
	return func(params *requestParamsDto) {
		params.header = header
	}
}

func WithContentType(contentType string) ReqParamsOption {
	return func(params *requestParamsDto) {
		params.contentType = contentType
	}
}

func SetBasicAuth(username, password string) ReqParamsOption {
	return func(params *requestParamsDto) {
		params.basicAuth = &BasicAuth{
			Username: username,
			Password: password,
		}
	}
}
