package wishlistlib

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type User struct {
	ID    uint64 `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Retrieves all the users in the database
func (wc *WishClient) GetAllUsers() ([]User, error) {
	users := []User{}
	resBody, err := wc.executeRequest("GET", "/user", nil, nil, false)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Searches users based on a given search string
func (wc *WishClient) SearchUsers(query string) ([]User, error) {
	users := []User{}
	resBody, err := wc.executeRequest("GET", "/user/search", map[string]string{
		"search": query,
	}, nil, false)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Retrieves a user by their email, returns the user from the server with its assigned ID
func (wc *WishClient) GetUserByEmail(email string) (User, error) {
	user := User{}
	resBody, err := wc.executeRequest("GET", "/user/email", map[string]string{
		"email": email,
	}, nil, false)
	if err != nil {
		return User{}, err
	}

	err = json.Unmarshal(resBody, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Retrieves a user by their email, returns the user from the server with its assigned ID
func (wc *WishClient) GetUserByID(id uint64) (User, error) {
	user := User{}
	resBody, err := wc.executeRequest("GET", "/user/"+strconv.FormatUint(id, 10), nil, nil, false)
	if err != nil {
		return User{}, err
	}

	err = json.Unmarshal(resBody, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Adds the given user to the server and returns the user with its new ID
func (wc *WishClient) AddNewUser(user User, password string) (User, error) {

	// Marshal user object
	userReq := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Name:     user.Name,
		Email:    user.Email,
		Password: password,
	}
	reqBody, err := json.Marshal(userReq)
	if err != nil {
		return User{}, err
	}

	_, err = wc.executeRequest("POST", "/user", nil, reqBody, false)
	if err != nil {
		return User{}, err
	}

	newUser, err := wc.GetUserByEmail(user.Email)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

// Alters a user based on the provided arguments (leave argument string empty to leave unchanged) Returns the user with changed fields
func (wc *WishClient) ChangeUser(user User, name, email, password string) error {
	userInfo := make(map[string]string)

	if name != "" {
		userInfo["name"] = name
	}

	if email != "" {
		userInfo["email"] = email
	}

	if password != "" {
		userInfo["password"] = password
	}

	// Marshal user data
	reqBody, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}

	_, err = wc.executeRequest("PUT", fmt.Sprint("/user/", strconv.FormatUint(user.ID, 10)), nil, reqBody, true)
	if err != nil {
		return err
	}

	return nil
}

// Deletes the given user from the database (WARNING: PERMANENT)
func (wc *WishClient) DeleteUser(user User) error {
	_, err := wc.executeRequest("DELETE", fmt.Sprint("/user/", strconv.FormatUint(user.ID, 10)), nil, nil, true)
	return err
}

// Converts the user to a string for debugging
func (u *User) String() string {
	return fmt.Sprintf("[%03d] %20s (%s)", u.ID, u.Name, u.Email)
}
