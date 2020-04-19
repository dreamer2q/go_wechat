package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"runtime"
)

func newFileForm(fieldname string, filename string, in io.Reader, out io.Writer) (*multipart.Writer, error) {
	m := multipart.NewWriter(out)
	w, err := m.CreateFormFile(fieldname, filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(w, in)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return m, nil
}

func DownFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func msgJsonBuilder(msg *Message) io.Reader {
	var (
		sb  = bytes.NewBuffer(nil)
		tmp []byte
	)
	sb.WriteByte('{')
	if msg.ToWxName != "" {
		_, _ = fmt.Fprintf(sb, `"towxname":%q,`, msg.ToWxName)
	}
	if msg.ToUser != nil {
		tmp, _ = json.Marshal(msg.ToUser)
		_, _ = fmt.Fprintf(sb, `"touser":%s,`, tmp)
	}
	if msg.Filter != nil {
		tmp, _ = json.Marshal(msg.Filter)
		_, _ = fmt.Fprintf(sb, `"filter":%s,`, tmp)
	}
	tmp, _ = json.Marshal(msg.MsgWrapper)
	if msg.MsgWrapper.Type() == "image" {
		_, _ = fmt.Fprintf(sb, `"images":%s,`, tmp)
	} else {
		_, _ = fmt.Fprintf(sb, "%q:%s,", msg.MsgWrapper.Type(), tmp)
	}
	_, _ = fmt.Fprintf(sb, `"send_ignore_reprint":%d,`, msg.IgnoreReprint)
	if msg.ClientMsgId != "" {
		_, _ = fmt.Fprintf(sb, `"clientmsgid":%q,`, msg.ClientMsgId)
	}
	_, _ = fmt.Fprintf(sb, `"msgtype":%q`, msg.MsgWrapper.Type())
	sb.WriteByte('}')
	return sb
}

func checkResult(body []byte) (*Result, error) {
	ret := &Result{}
	err := json.Unmarshal(body, ret)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if err != nil {
		return nil, fmt.Errorf("in %s at %d func %s: unmarshal: %v", file, line, f.Name(), err)
	}
	if ret.ErrCode != 0 {
		return nil, fmt.Errorf("in %s at %d func %s: %d %s", file, line, f.Name(), ret.ErrCode, ret.ErrMsg)
	}
	return ret, nil
}
