package media

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func newFileForm(fieldname string, filename string, in io.Reader, out io.Writer) (*multipart.Writer, error) {
	m := multipart.NewWriter(out)
	w, err := m.CreateFormFile(fieldname, filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(w, in)
	if err != nil || err != io.EOF {
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
