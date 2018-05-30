package types

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Error struct {
	Status int `json:"-"`
	ErrMsg string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("(%d)%s", e.Status, e.ErrMsg)
}

func ParseErr(resp *http.Response) *Error {
	if resp.StatusCode >= http.StatusBadRequest {
		msg, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		var gitErr Error
		gitErr.Status = resp.StatusCode
		if len(msg) > 0 {
			json.Unmarshal(msg, &gitErr)
			if gitErr.ErrMsg == "" {
				gitErr.ErrMsg = string(msg)
			}
		}

		return &gitErr
	}

	return nil
}

func IsNotFoundErr(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Status == http.StatusNotFound
	}

	return false
}

func isResourceNotFoundErr(resource string, err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Status == http.StatusNotFound && e.ErrMsg == fmt.Sprintf("404 %s Not Found", resource)
	}

	return false
}

func IsFileNotFoundErr(err error) bool {
	return isResourceNotFoundErr("File", err)
}

func IsProjectNotFoundErr(err error) bool {
	return isResourceNotFoundErr("Project", err)
}