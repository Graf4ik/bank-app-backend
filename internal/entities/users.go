package entities

// User represents a user in the system.
// @Description User model
// @example { "id": 1, "email": "user@example.com", "username": "user1", password: "123456" }
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Username string
	Password string
}

// UserResponse represents the public view of a user, safe to be returned in API responses.
// @Description Public user information without sensitive fields like password.
// @example { "id": 1, "email": "user@example.com", "username": "user1" }
type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// UpdateUserRequest is used to update user fields.
// @Description Update user model
// @example { "email": "new@example.com", "username": "newusername" }
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
	}
}

// UsersToResponse converts a list of Users to a list of UserResponses.
func UsersToResponse(users []*User) []*UserResponse {
	responses := make([]*UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}
	return responses
}
