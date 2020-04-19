package message

import "../request"

type Once struct {
	req *request.Request
}

func newOnce(r *request.Request) *Once {
	return &Once{r}
}

//TODO implement once message
