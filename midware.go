package wechat

import (
	"bytes"
	"crypto/sha1"
	"encoding/xml"
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
			return
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

//encryptedXMLMsg 安全模式下的消息体
type encryptedXMLMsg struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	ToUserName   string   `xml:"ToUserName" json:"ToUserName"`
	EncryptedMsg string   `xml:"Encrypt"    json:"Encrypt"`
}
type encryptorWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (e *encryptorWriter) Write(p []byte) (int, error) {
	return e.body.Write(p)
}

type encryptedXmlReply struct {
	XMLName    struct{} `xml:"xml" json:"-"`
	EncryptMsg string   `xml:"Encrypt"`
	TimeStamp  string   `xml:"TimeStamp"`
	Nonce      string   `xml:"Nonce"`
}

func (w *API) encryptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("encrypt_type") == "" {
			return
		}
		msgSignature := c.Query("msg_signature")
		if msgSignature == "" {
			log.Printf("encryptor: msg_signature is empty")
			c.AbortWithStatus(http.StatusBadRequest)
		}
		encXml := &encryptedXMLMsg{}
		if err := xml.NewDecoder(c.Request.Body).Decode(encXml); err != nil {
			log.Printf("err: %v", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		sigGen := signature(w.config.AppToken, timestamp, nonce, encXml.EncryptedMsg)
		if sigGen != msgSignature {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("signature check failed"))
			return
		}
		random, rawXmlBytes, err := decryptMsg(w.config.AppID, encXml.EncryptedMsg, w.config.AesEncodeKey)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		//替换body为解密后的内容
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(rawXmlBytes))
		//同时，截取body，将之加密后发送
		bufWriter := &encryptorWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = bufWriter
		//control flow
		c.Next()
		rawXmlBytes, _ = ioutil.ReadAll(bufWriter.body)
		encXmlMsg, err := encryptMsg(random, rawXmlBytes, w.config.AppID, w.config.AesEncodeKey)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		_, err = bufWriter.ResponseWriter.Write(encXmlMsg)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
