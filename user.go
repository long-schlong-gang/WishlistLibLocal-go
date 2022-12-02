package wishlistlib

import (
	"fmt"
	"regexp"
	"strings"
)

type User struct {
	ID    uint64 `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Retrieves all the users in the database
func (wc *WishClient) GetAllUsers() ([]User, error) {
	users := []User{}
	for id, ju := range wc.Users {
		users = append(users, User{
			ID:    id,
			Name:  ju.Name,
			Email: ju.Email,
		})
	}
	return users, nil
}

// Searches users based on a given search string
func (wc *WishClient) SearchUsers(query string) ([]User, error) {
	users := []User{}

	// Filter by Name or Email
	filter := regexp.MustCompile(fmt.Sprintf("(?i)%s", query))
	for id, ju := range wc.Users {
		if filter.Match([]byte(ju.Name)) || filter.Match([]byte(ju.Email)) {
			users = append(users, User{
				ID:    id,
				Name:  ju.Name,
				Email: ju.Email,
			})
		}
	}

	return users, nil
}

// Retrieves a user by their email, returns the user from the server with its assigned ID
func (wc *WishClient) GetUserByEmail(email string) (User, error) {
	var user *User
	for id, ju := range wc.Users {
		if strings.EqualFold(ju.Email, email) {
			user = &User{
				ID:    id,
				Name:  ju.Name,
				Email: ju.Email,
			}
		}
	}
	if user == nil {
		return User{}, NotFoundError(fmt.Sprintf("User with email '%s' not found", email))
	}

	return *user, nil
}

// Retrieves a user by their email, returns the user from the server with its assigned ID
func (wc *WishClient) GetUserByID(id uint64) (User, error) {
	ju, exists := wc.Users[id]
	if !exists {
		return User{}, NotFoundError(fmt.Sprintf("User with ID '%d' not found", id))
	}
	return User{
		ID:    id,
		Name:  ju.Name,
		Email: ju.Email,
	}, nil
}

// Adds the given user to the server and returns the user with its new ID
func (wc *WishClient) AddNewUser(user User, password string) (User, error) {
	password, err := hashPassword(password)

	// Check user doesn't already exist
	if _, err := wc.GetUserByEmail(user.Email); err == nil {
		return User{}, EmailExistsError(fmt.Sprintf("User with email '%s' already exists", user.Email))
	}

	if err != nil {
		return User{}, err
	}
	wc.Users[wc.NextUserID] = JSONUser{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: password,
	}
	wc.NextUserID++

	newUser, err := wc.GetUserByEmail(user.Email)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

// Alters a user based on the provided arguments (leave argument string empty to leave unchanged) Returns the user with changed fields
func (wc *WishClient) ChangeUser(user User, name, email, password string) error {

	// Check user exists
	ju, exists := wc.Users[user.ID]
	if !exists {
		return NotFoundError(fmt.Sprintf("User with ID '%d' not found", user.ID))
	}

	if name != "" {
		ju.Name = name
	}

	if email != "" {
		ju.Email = email
	}

	if password != "" {
		hash, err := hashPassword(password)
		if err != nil {
			return err
		}
		ju.PasswordHash = hash
	}

	wc.Users[user.ID] = ju
	return nil
}

// Deletes the given user from the database (WARNING: PERMANENT)
func (wc *WishClient) DeleteUser(user User) error {
	if _, exists := wc.Users[user.ID]; !exists {
		return NotFoundError(fmt.Sprintf("User with ID '%d' not found", user.ID))
	}

	delete(wc.Users, user.ID)
	return nil
}

// Converts the user to a string for debugging
func (u *User) String() string {
	return fmt.Sprintf("[%05d] %20s (%s)", u.ID, u.Name, u.Email)
}
