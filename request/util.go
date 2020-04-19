package request

import (
	"encoding/json"
	"fmt"
	"runtime"
)

func CheckCommonError(body []byte) error {
	ret := &CommonError{}
	err := json.Unmarshal(body, ret)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if err != nil {
		return fmt.Errorf("in %s at %d func %s: unmarshal: %v", file, line, f.Name(), err)
	}
	if ret.ErrCode != 0 {
		return fmt.Errorf("in %s at %d func %s: %d %s", file, line, f.Name(), ret.ErrCode, ret.ErrMsg)
	}
	return nil
}
