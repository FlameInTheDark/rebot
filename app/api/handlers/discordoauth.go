package handlers

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/FlameInTheDark/rebot/app/api/handlers/requests"
	"github.com/FlameInTheDark/rebot/business/service/auth"
	"github.com/FlameInTheDark/rebot/foundation/berrors"
	"github.com/FlameInTheDark/rebot/foundation/web"
)

type AuthHandlersGroup struct {
}

func NewAuthHandlersGroup() *AuthHandlersGroup {
	return &AuthHandlersGroup{}
}

//OAuthCallbackHandler accept discord oauth requests
func (a *AuthHandlersGroup) OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	var req requests.OauthCode
	if err := render.Bind(r, &req); err != nil {
		web.RespondError(w, r, berrors.WrapWithError(auth.ErrInvalidInput, err))
		return
	}
}
