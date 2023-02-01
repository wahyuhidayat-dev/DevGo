package user
//! POST REGISTER USER
type RegisterUserInput struct {
	Name       string `json:"name" binding:"required"`
	Occupation string `json:"occupation" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}
//! LOGIN USER
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
//! CHECK EMAIL ALREADY OR NOT
type CheckEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}