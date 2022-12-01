package wishlistlib

import (
	"testing"
)

const (
	DATA_FILE = "testdata.json"

	TEST_USER_NAME  = "Jim Test"
	TEST_USER_EMAIL = "jim.test@example.com"
	TEST_USER_PASS  = "beowulf"
	TEST_ITEM_NAME  = "Test Item"
	TEST_ITEM_DESC  = "An item to test the API"
	TEST_ITEM_PRICE = 10
)

func TestCRUD(t *testing.T) {

	var err error
	var user User
	var users []User
	var item Item
	var items []Item

	wc, err := DefaultWishClient(DATA_FILE)
	t.Log("  Error:", err)

	t.Log("Getting all users...")
	users, err = wc.GetAllUsers()
	t.Log("  Error:", err)
	logUserList(t, users)

	t.Log("Creating new user...")
	user, err = wc.AddNewUser(User{
		Name:  TEST_USER_NAME,
		Email: TEST_USER_EMAIL,
	}, TEST_USER_PASS)
	t.Log("  Error:", err)

	t.Log("Getting new guy by email...")
	user, err = wc.GetUserByEmail(TEST_USER_EMAIL)
	t.Log("  Error:", err)

	t.Log("Authenticating...")
	err = wc.Authenticate(user.Email, TEST_USER_PASS)
	t.Log("  Error:", err)

	t.Log("Getting all of Guy's items...")
	items, err = wc.GetAllItemsOfUser(user)
	t.Log("  Error:", err)
	logItemList(t, items)

	t.Log("Add New Item...")
	item, err = wc.AddItemOfUser(Item{
		Name:        TEST_ITEM_NAME,
		Description: TEST_ITEM_DESC,
		Price:       TEST_ITEM_PRICE,
	}, user)
	t.Log("  Error:", err)
	t.Log("Last Added Item:", item)

	items, _ = wc.GetAllItemsOfUser(user)
	logItemList(t, items)

	t.Log("Reserving new item...")
	err = wc.ReserveItemOfUser(item, user)
	t.Log("  Error:", err)

	items, _ = wc.GetAllItemsOfUser(user)
	logItemList(t, items)

	item, _ = wc.GetItemByID(item.ItemID)
	t.Log("Reserved By: ", item.ReservedByUser.Name)

	t.Log("Un-Reserving new item...")
	err = wc.UnreserveItemOfUser(item, user)
	t.Log("  Error:", err)

	items, _ = wc.GetAllItemsOfUser(user)
	logItemList(t, items)

	t.Log("Deleting the new item...")
	wc.DeleteItemOfUser(item, user)
	t.Log("  Error:", err)

	items, _ = wc.GetAllItemsOfUser(user)
	logItemList(t, items)

	t.Log("Changing user's name...")
	err = wc.ChangeUser(user, "Fred Test", "", "")
	t.Log("  Error:", err)

	users, _ = wc.GetAllUsers()
	logUserList(t, users)

	t.Log("Deleting user...")
	err = wc.DeleteUser(user)
	t.Log("  Error:", err)

	users, _ = wc.GetAllUsers()
	logUserList(t, users)

	err = wc.Close()
	t.Log("  Error:", err)
}

// util funcs

func logUserList(t *testing.T, us []User) {
	t.Log("  All Users:")
	for _, u := range us {
		t.Log("   - " + u.String())
	}
}

func logItemList(t *testing.T, is []Item) {
	t.Log("  All Items:")
	for _, i := range is {
		t.Log("   - " + i.String())
	}
}
