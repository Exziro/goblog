package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePage(t *testing.T) {
	baseURL := "http://localhost:3000/"

	//请求 模仿用户数据
	var (
		resp *http.Response
		err  error
	)
	resp, err = http.Get(baseURL + "/")

	//检测是否无错且200

	assert.NoError(t, err, "有错误发生，err 不为空")
	assert.Equal(t, 200, resp.StatusCode, "应返回状态码200")

}
