package wechat

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type requestVerify struct {
	Signature string `form:"signature" binding:"required"`
	EchoStr   string `form:"echostr"`
	Timestamp string `form:"timestamp"`
	Nonce     string `form:"nonce"`
}

func (w *API) verifier() gin.HandlerFunc {
	return func(c *gin.Context) {
		rv := requestVerify{}
		if err := c.ShouldBindQuery(&rv); err != nil {
			log.Printf("verifier: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		sigGen := signature(w.config.AppToken, rv.Timestamp, rv.Nonce)
		if rv.Signature != sigGen {
			c.AbortWithStatus(http.StatusBadRequest)
			log.Printf("verifier: check failed")
			return
		}

		if c.Request.Method == "GET" {
			c.String(http.StatusOK, "%s", rv.EchoStr)
			c.Abort()
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

//额外拷贝一份输出流,实现logger的记录
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

//encryptedXmlRecv 安全模式下的消息体 接受方面
type encryptedXmlRecv struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	ToUserName   string   `xml:"ToUserName" json:"ToUserName"`
	EncryptedMsg string   `xml:"Encrypt"    json:"Encrypt"`
}

//劫持发送的消息到Buffer里面
type encryptorWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (e *encryptorWriter) Write(p []byte) (int, error) {
	return e.body.Write(p)
}

//encryptedXmlReply 安全模式下的发送消息结构体
type encryptedXmlReply struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	EncryptMsg   string   `xml:"Encrypt"`
	MsgSignature string   `xml:"MsgSignature"`
	TimeStamp    string   `xml:"TimeStamp"`
	Nonce        string   `xml:"Nonce"`
}

type requestEncrypt struct {
	requestVerify
	EncryptedType string `form:"encrypt_type"`
	MsgSignature  string `form:"msg_signature"`
}

//处理安全模式下的消息解密和加密过程
func (w *API) encryptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		//检测消息是否加密
		reqEnc := requestEncrypt{}
		_ = c.ShouldBindQuery(&reqEnc)
		if reqEnc.EncryptedType == "" {
			return
		}
		if w.config.AesEncodeKey == "" {
			panic(errors.New("config: AesEncodeKey is empty"))
		}
		//对收到的消息进行解密，和再封装
		eXmlRecv := &encryptedXmlRecv{}
		if err := xml.NewDecoder(c.Request.Body).Decode(eXmlRecv); err != nil {
			log.Printf("%#v", errors.Wrap(err, "decode xml message error"))
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		//进行消息体的验证
		sigGen := signature(w.config.AppToken, reqEnc.Timestamp, reqEnc.Nonce, eXmlRecv.EncryptedMsg)
		if sigGen != reqEnc.MsgSignature {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("signature check failed"))
			return
		}
		//解密消息
		random, rawXmlBytes, err := decryptMsg(w.config.AppID, eXmlRecv.EncryptedMsg, w.config.AesEncodeKey)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		//替换body为解密后的内容,实现透明代理
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(rawXmlBytes))
		//同时，截取Writer，便于后面加密消息
		bufWriter := &encryptorWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = bufWriter
		//control flow
		c.Next()
		rawXmlBytes, _ = ioutil.ReadAll(bufWriter.body)
		encRawReply, err := encryptMsg(random, rawXmlBytes, w.config.AppID, w.config.AesEncodeKey)
		if err != nil {
			log.Printf("%#v", errors.Wrap(err, "encrypt message"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		replySignature := signature(w.config.AppToken, reqEnc.Timestamp, reqEnc.Nonce, string(encRawReply))
		xmlReply := encryptedXmlReply{
			EncryptMsg:   string(encRawReply),
			MsgSignature: replySignature,
			TimeStamp:    reqEnc.Timestamp,
			Nonce:        reqEnc.Nonce,
		}
		//发送加密后的消息
		c.Writer = bufWriter.ResponseWriter
		c.XML(http.StatusOK, &xmlReply)
	}
}
