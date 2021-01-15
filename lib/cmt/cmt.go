package cmt

type AuthorizationCodeRequest struct {
	AccountID string `json:"account_id"`
}

type AuthorizationCodeResponse struct {
	AuthorizationCode string `json:"auth_code"`
}
