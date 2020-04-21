package wechat

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

type requestVerify struct {
	Signature string `form:"signature" binding:"required"`
	EchoStr   string `form:"echostr"`
	Timestamp string `form:"timestamp"`
	Nonce     string `form:"nonce"`
}

func (w *API) verifier() gin.HandlerFunc {
	token := w.config.AppToken
	return func(context *gin.Context) {
		rv := requestVerify{}
		err := context.MustBindWith(&rv, binding.Query)
		if err != nil {
			log.Printf("verify: %v\n", err)
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		strArr := make([]string, 0)
		strArr = append(strArr, token)
		strArr = append(strArr, rv.Timestamp)
		strArr = append(strArr, rv.Nonce)
		sort.Strings(strArr)
		unsigned := strings.Join(strArr, "")
		s := sha1.New()
		s.Write([]byte(unsigned))
		signed := fmt.Sprintf("%x", s.Sum(nil))

		if rv.Signature != signed {
			context.AbortWithStatus(http.StatusBadRequest)
			log.Printf("verifier: request payload error")
			return
		}

		if context.Request.Method == "GET" {
			context.String(http.StatusOK, "%s", rv.EchoStr)
			context.Abort()
			log.Printf("verifier: request verified")
		}

		if w.config.Debug {
			log.Printf("verifier: check passed")
		}
	}
}

func (w *API) logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		if !w.config.Debug {
			return ""
		}
		sb := &strings.Builder{}
		_, _ = fmt.Fprintf(sb, "[%s] from %s \"%s %s %s %d\" [%s] %s\n",
			params.TimeStamp.Format("2006-01-02 15:04:05"),
			params.ClientIP,
			params.Method,
			params.Path,
			params.Request.Proto,
			params.StatusCode,
			params.Latency,
			params.ErrorMessage,
		)
		return sb.String()
	})
}

type bodyLoggerWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (b *bodyLoggerWriter) Write(bin []byte) (int, error) {
	_, _ = b.body.Write(bin)
	return b.ResponseWriter.Write(bin)
}

func (w *API) debugger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if w.config.Debug {
			if c.Request.Method == "POST" {
				body, _ := c.GetRawData()
				log.Printf("receive: %s", body)
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
			}
			blg := &bodyLoggerWriter{
				ResponseWriter: c.Writer,
				body:           bytes.NewBuffer(nil),
			}
			c.Writer = blg
			//control flow
			c.Next()

			//log response body
			log.Printf("send: %s\n", blg.body.String())
		}
	}
}
