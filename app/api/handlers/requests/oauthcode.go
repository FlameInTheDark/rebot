package requests

//OauthCode contains code to authorize person using discord's oauth2 api
type OauthCode struct {
	Code string `json:"code" form:"code" validate:"required"`
}
