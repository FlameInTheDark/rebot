package auth

type Service struct {
	AuthURL  string
	TokenURL string
	Scopes   []string
}

//func (s *Service) Auth(ctx context.Context, id, secret string) error {
//	conf := &oauth2.Config{
//		ClientID:     id,
//		ClientSecret: secret,
//		Scopes:       s.Scopes,
//		Endpoint: oauth2.Endpoint{
//			AuthURL:  s.AuthURL,
//			TokenURL: s.TokenURL,
//		},
//	}
//
//	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
//
//	var code string
//	if _, err := fmt.Scan(&code); err != nil {
//		return err
//	}
//
//	tok, err := conf.Exchange(ctx, code)
//	if err != nil {
//		return err
//	}
//
//	client := conf.Client(ctx, tok)
//	client.Get()
//}
