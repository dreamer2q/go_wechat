package wechat

type ErrorWrapper []func() error

func (e *ErrorWrapper) Add(f func() error) *ErrorWrapper {
	*e = append(*e, f)
	return e
}

func (e *ErrorWrapper) Run() error {
	for _, f := range *e {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
