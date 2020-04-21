package midware

import (
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

func Verifier(token string) gin.HandlerFunc {
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
			log.Printf("verify: request payload error")
			return
		}
		if context.Request.Method == "GET" {
			context.String(http.StatusOK, "%s", rv.EchoStr)
			context.Abort()
			log.Printf("verify: request verified\n")
		}
		//flow control
		//context.Next()
		context.Set("openid", context.DefaultQuery("openid", ""))
	}
}

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		sb := &strings.Builder{}
		_, _ = fmt.Fprintf(sb, "%s - [%s] \"%s %s %s %d\" [%s] %s\n",
			params.ClientIP,
			params.TimeStamp.Format("2006-01-02 15:04:05"),
			params.Method,
			params.Path,
			params.Request.Proto,
			params.StatusCode,
			params.Latency,
			params.ErrorMessage,
		)
		if params.Method == "POST" {
			if body, err := ioutil.ReadAll(params.Request.Body); err == nil {
				_, _ = fmt.Fprintf(sb, "%s\n", body)
			} else {
				_, _ = fmt.Fprintf(sb, "body: %v\n", err)
			}
		}
		return sb.String()
	})
}
