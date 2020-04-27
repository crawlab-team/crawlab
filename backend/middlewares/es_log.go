package middlewares

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

func EsLog(ctx context.Context, esClient *elastic.Client) gin.HandlerFunc {

	return func(c *gin.Context) {
		// 开始时间
		crawlabIndex := viper.GetString("setting.crawlabLogIndex")
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := strconv.FormatInt(int64(end.Sub(start).Milliseconds()), 10)
		path := c.Request.URL.Path

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := strconv.Itoa(c.Writer.Status())
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		b := buf.String()
		accessLog := "costTime:" + latency + "ms--" + "StatusCode:" + statusCode + "--" + "Method:" + method + "--" + "ClientIp:" + clientIP + "--" +
			"RequestURI:" + path + "--" + "Host:" + c.Request.Host + "--" + "UserAgent--" + c.Request.UserAgent() + "--RequestBody:" +
			string(b)
		WriteMsg(ctx, crawlabIndex, esClient, time.Now(), accessLog)
	}

}

// WriteMsg will write the msg and level into es
func WriteMsg(ctx context.Context, crawlabIndex string, es *elastic.Client, when time.Time, msg string) error {
	vals := make(map[string]interface{})
	vals["@timestamp"] = when.Format(time.RFC3339)
	vals["@msg"] = msg
	uid := uuid.NewV4().String()
	_, err := es.Index().Index(crawlabIndex).Id(uid).BodyJson(vals).Refresh("wait_for").Do(ctx)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
