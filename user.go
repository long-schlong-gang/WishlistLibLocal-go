package wishlistlib

import "fmt"

type User struct {
	ID     uint64 `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Domain string `json:"domain"`
}

// Retrieves all the users in the database
func (ctx *Context) GetAllUsers() ([]User, error) {
	users := []User{}
	err := ctx.parseObjectFromServer("/user", "GET", &users, nil)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Searches users based on a given search string
func (ctx *Context) SearchUsers(query string) ([]User, error) {
	users := []User{}
	err := ctx.parseObjectFromServer("/user/search", "GET", &users, map[string]string{
		"search": query,
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Converts the user to a string for debugging
func (u User) String() string {
	return fmt.Sprintf("[%v] %v (%v)", u.ID, u.Name, u.Email)
}
