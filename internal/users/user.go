package users

// User of user model
type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username,omitempty"`
	FullName  string `json:"full_name,omitempty"`
	Messages []User `json:"following,omitempty"`
}

// UserRepository represent repositories of the User
// Expect implementation by the infrastructure layer
type UserRepository interface {
	Get(id int) (*User, error)
	GetAll() ([]*User, error)
	Save(*User) error
}
