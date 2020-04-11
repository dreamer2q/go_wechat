package media

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func newFilePost(fieldname string, filename string, out io.Writer, description ...VideoDescription) (*multipart.Writer, error) {
	m := multipart.NewWriter(out)
	w, err := m.CreateFormFile(fieldname, filename)
	if err != nil {
		return nil, err
	}
	fw, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(w, fw)
	if err != nil {
		return nil, err
	}
	if description != nil {
		dw, err := m.CreateFormField("description")
		if err != nil {
			return nil, err
		}
		dwBytes, err := json.Marshal(description[0])
		if err != nil {
			return nil, err
		}
		_, _ = dw.Write(dwBytes)
	}
	//close
	_ = m.Close()
	_ = fw.Close()
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
