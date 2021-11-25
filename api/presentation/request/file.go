package request

import "github.com/pkg/errors"

type CreateFileRequest struct {
	Message string `json:"message"`
	Grade   int    `json:"grade"`
}

func (r *CreateFileRequest) Validation() error {
	if r.Message == "" {
		return errors.New("required field is empty")
	}
	if r.Grade < 1 {
		return errors.New("grande don't allow minas")
	}
	return nil
}
