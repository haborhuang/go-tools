package types

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Error struct {
	Message string `json:"message"`
	Status int `json:"-"`
}

func (e *Error) Error() string {
	if nil == e {
		return ""
	}

	return fmt.Sprintf("Kong response error: %d - %s", e.Status, e.Message)
}

func IsNotFoundErr(err error) bool {
	kongErr, ok := err.(*Error)
	return ok && kongErr.Status == http.StatusNotFound
}

func ParseErr(resp *http.Response) *Error {
	if resp.StatusCode >= http.StatusBadRequest {
		msg, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		var gitErr Error
		gitErr.Status = resp.StatusCode
		if len(msg) > 0 {
			json.Unmarshal(msg, &gitErr)
			if gitErr.Message == "" {
				gitErr.Message = string(msg)
			}
		}

		return &gitErr
	}

	return nil
}