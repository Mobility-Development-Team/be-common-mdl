package common

import (
	"time"

	"github.com/go-resty/resty/v2"
)

func NewResty() *resty.Client {
	client := resty.New()
	client.SetRetryCount(3). // 设置最大重试次数为 3 次
					SetRetryWaitTime(2 * time.Second).    // 设置初始等待时间为 2 秒
					SetRetryMaxWaitTime(15 * time.Second) // 设置最大等待时间为 15 秒 (指数退避)
	client.AddRetryCondition(func(r *resty.Response, err error) bool {
		// 自定义重试条件：如果状态码是 500-599 或 429 (限流)，则重试
		if r == nil {
			return true // 如果响应为空（如网络错误），也重试
		}
		return r.StatusCode() >= 400
	})

	return client
}
