package api

// NewUser contains all the information needed to create a new user in the system
type NewUser struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// OK is a success response object
type OK struct {
}

// Error represnts an expected error from your call
type Error struct {
	ErrorCode string `json:"error_code"`
}
