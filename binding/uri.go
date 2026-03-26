package binding

import (
	"net/http"
	"strings"
)

type uriBinding struct{}

func (uriBinding) Name() string {
	return "uri"
}

func (uriBinding) Bind(req *http.Request, obj any) error {
	m := make(map[string][]string)
	pat := req.Pattern
	url := req.RequestURI
	_, pvMap := cleanPatternAndParams(pat, url)
	for k, v := range pvMap {
		m[k] = []string{v}
	}

	if err := mapURI(obj, m); err != nil {
		return err
	}

	return nil
}

func cleanPatternAndParams(pattern, urlPath string) (cleaned string, params map[string]string) {
	params = make(map[string]string)

	// 如果 pattern 带方法前缀（如 "GET /users/{id}"），先去掉方法部分
	if idx := strings.Index(pattern, " "); idx != -1 {
		pattern = pattern[idx+1:] // 只保留路径部分
	}

	// 按 / 分割 pattern 和 urlPath
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	urlParts := strings.Split(strings.Trim(urlPath, "/"), "/")

	cleanParts := make([]string, 0, len(urlParts))

	for i := range patternParts {
		if i >= len(urlParts) {
			break
		}

		p := patternParts[i]
		u := urlParts[i]

		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
			// 提取参数名，去掉前后 {}
			paramName := strings.Trim(p, "{}")
			params[paramName] = u
			cleanParts = append(cleanParts, u) // 用真实值替换
		} else {
			// 静态部分，直接保留
			cleanParts = append(cleanParts, u)
		}
	}

	// 处理剩余的 url 部分（例如通配符 * 的情况，这里简化处理）
	for i := len(patternParts); i < len(urlParts); i++ {
		cleanParts = append(cleanParts, urlParts[i])
	}

	cleaned = "/" + strings.Join(cleanParts, "/")
	if urlPath == "/" || cleaned == "" {
		cleaned = "/"
	}

	return cleaned, params
}
