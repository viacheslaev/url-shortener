package account

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	AccountID string `json:"account_id"` // public_id UUID
}

func createRegistrationResponse(publicId string) RegisterResponse {
	return RegisterResponse{
		AccountID: publicId,
	}
}
