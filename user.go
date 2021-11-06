package wishlistlib

import (
	"fmt"
	"strconv"
)

type User struct {
	ID       uint64 `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Domain   string `json:"domain"`
	password string
}

func (u *User) SetPassword(p string) { u.password = p }

// Retrieves all the users in the database
func (ctx *Context) GetAllUsers() ([]User, error) {
	users := []User{}
	err := ctx.parseObjectFromServer("/user", "GET", &users, nil, false)
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
	}, false)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Retrieves a user by their email, returns the user from the server with its assigned ID
func (ctx *Context) GetUserByEmail(email string) (User, error) {
	user := User{}
	err := ctx.parseObjectFromServer("/user/email", "GET", &user, map[string]string{
		"email": email,
	}, false)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Retrieves a user by their email, returns the user from the server with its assigned ID
func (ctx *Context) GetUserByID(id uint64) (User, error) {
	user := User{}
	err := ctx.parseObjectFromServer("/user/"+strconv.FormatUint(id, 10), "GET", &user, nil, false)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Adds the given user to the server and returns the user with its new ID
func (ctx *Context) AddNewUser(user User) (User, error) {

	if user.password == "" {
		return User{}, NoPasswordProvidedError(0)
	}

	err := ctx.sendObjectToServer("/user", "POST", struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.password,
	}, false)
	if err != nil {
		return User{}, err
	}

	newUser, err := ctx.GetUserByEmail(user.Email)
	if err != nil {
		return User{}, err
	}
	newUser.password = user.password
	return newUser, nil
}

// Alters a user based on the provided arguments (leave argument string empty to leave unchanged) Returns the user with changed fields
func (ctx *Context) ChangeAuthenticatedUser(name, email, password string) error {
	if ctx.authUser == (User{}) {
		return NoAuthenticatedUserError(0)
	}

	userInfo := make(map[string]string)

	if name != "" {
		userInfo["name"] = name
	}

	userEmail := ctx.authUser.Email
	if email != "" {
		userInfo["email"] = email
		userEmail = email
	}

	userPword := ctx.authUser.password
	if password != "" {
		userInfo["password"] = password
		userPword = password
	}

	err := ctx.sendObjectToServer("/user/"+strconv.FormatUint(ctx.authUser.ID, 10), "PUT", userInfo, true)
	if err != nil {
		return err
	}

	user, err := ctx.GetUserByEmail(userEmail)
	if err != nil {
		return err
	}
	ctx.authUser = user
	ctx.authUser.password = userPword

	return nil
}

// Deletes the given user from the database (WARNING: PERMANENT)
func (ctx *Context) DeleteAuthenticatedUser() error {
	if ctx.authUser == (User{}) {
		return NoAuthenticatedUserError(0)
	}
	return ctx.simpleRequest("/user/"+strconv.FormatUint(ctx.authUser.ID, 10), "DELETE", true)
}

// Converts the user to a string for debugging
func (u *User) String() string {
	return fmt.Sprintf("[%v] %v (%v)", u.ID, u.Name, u.Email)
}
