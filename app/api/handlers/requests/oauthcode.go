package requests

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var _ render.Binder = (*OauthCode)(nil)

//OauthCode contains code to authorize person using discord's oauth2 api
type OauthCode struct {
	Code string `json:"code"`
}

func (a *OauthCode) Bind(_ *http.Request) error {
	if a.Code == "" {
		return errors.New("\"code\" must be non empty")
	}
	return nil
}
