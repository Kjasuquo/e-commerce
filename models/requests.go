package models

type PasswordResetReq struct {
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}
type ForgotPassword struct {
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type RatingRequest struct {
	Rating uint `json:"rating"`
	Id     uint `json:"Id"`
}
