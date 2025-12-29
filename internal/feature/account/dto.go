package account

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"` // public_id UUID
}

func createRegistrationResponse(publicId string) RegisterResponse {
	return RegisterResponse{
		UserID: publicId,
	}
}
