package wishlistlib

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Item struct {
	ItemID         uint64  `json:"item_id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          float32 `json:"price"`
	ReservedByUser User    `json:"reserved_by"`
	Status         Status  `json:"status"`
	Links          []Link  `json:"links"`
}

type Link struct {
	LinkID uint64 `json:"link_id"`
	Text   string `json:"text"`
	URL    string `json:"hyperlink"`
}

type Status struct {
	StatusID    uint64 `json:"status_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Retrieves all the items from a user's list
func (wc *WishClient) GetAllItemsOfUser(user User) ([]Item, error) {
	items := []Item{}
	resBody, err := wc.executeRequest("GET", "/user/"+strconv.FormatUint(user.ID, 10)+"/list", nil, nil, false)
	if err != nil {
		return nil, err
	}

	// Parse objects
	err = json.Unmarshal(resBody, &items)
	if err != nil {
		return nil, err
	}

	// Convert prices
	if len(items) > 0 {
		for k, i := range items {
			items[k].Price = float32(i.Price) / 100
		}
	}

	return items, nil
}

// Retrieves an item with a certain ID
func (wc *WishClient) GetItemByID(id uint64) (Item, error) {
	item := Item{}
	resBody, err := wc.executeRequest("GET", "/item/"+strconv.FormatUint(id, 10), nil, nil, false)
	if err != nil {
		return Item{}, err
	}

	// Parse objects
	err = json.Unmarshal(resBody, &item)
	if err != nil {
		return Item{}, err
	}

	// Convert price
	item.Price = float32(item.Price) / 100

	return item, nil
}

// Adds an item to a user's list and returns the item with all its info
func (wc *WishClient) AddItemOfUser(item Item, user User) (Item, error) {

	// Check price is within limits
	intPrice := uint32(item.Price * 100)
	if uint64(intPrice) != uint64(item.Price*100) {
		return Item{}, PriceOutOfRangeError(item.Price * 100)
	}

	obj := struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		Price          uint32 `json:"price"`
		ReservedByUser *User  `json:"reserved_by,omitempty"`
		Status         struct {
			StatusID uint64 `json:"status_id"`
		} `json:"status"`
		Links []struct {
			Text      string `json:"text"`
			Hyperlink string `json:"hyperlink"`
		} `json:"links"`
	}{
		Name:           item.Name,
		Description:    item.Description,
		Price:          intPrice,
		ReservedByUser: &item.ReservedByUser,
		Status: struct {
			StatusID uint64 `json:"status_id"`
		}{StatusID: item.Status.StatusID},
		Links: make([]struct {
			Text      string `json:"text"`
			Hyperlink string `json:"hyperlink"`
		}, 0),
	}
	for _, l := range item.Links {
		obj.Links = append(obj.Links, struct {
			Text      string `json:"text"`
			Hyperlink string `json:"hyperlink"`
		}{
			Text:      l.Text,
			Hyperlink: l.URL,
		})
	}

	// Set status to available if not provided
	if obj.Status.StatusID == 0 {
		obj.Status.StatusID = 1
	}

	// Marshal object
	reqBody, err := json.Marshal(obj)
	if err != nil {
		return Item{}, err
	}

	_, err = wc.executeRequest("POST", fmt.Sprint("/user/", strconv.FormatUint(user.ID, 10), "/list"), nil, reqBody, true)
	if err != nil {
		return Item{}, err
	}

	items, err := wc.GetAllItemsOfUser(user)
	if err != nil {
		return Item{}, err
	} else if len(items) < 1 {
		return Item{}, AddingItemFailed("No Items were inserted to user's list, but request succeeded")
	}
	lastItem := items[0]
	for _, i := range items {
		if i.ItemID > lastItem.ItemID {
			lastItem = i
		}
	}

	return lastItem, nil
}

// Delete an item from the database (WARNING! PERMANENT!)
func (wc *WishClient) DeleteItemOfUser(item Item, user User) error {
	_, err := wc.executeRequest("DELETE", fmt.Sprint("/user/", strconv.FormatUint(user.ID, 10), "/list/", strconv.FormatUint(item.ItemID, 10)), nil, nil, true)
	return err
}

// Sets the status of the provided item
func (wc *WishClient) SetItemStatusOfUser(item Item, user User, status Status) error {
	obj := struct {
		Status `json:"status"`
	}{
		Status: status,
	}

	// Marshal object
	reqBody, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = wc.executeRequest("PUT", fmt.Sprint("/user/", strconv.FormatUint(user.ID, 10), "/list/", strconv.FormatUint(item.ItemID, 10)), nil, reqBody, true)
	return err
}

// Reserves an item using the reserve endpoint
func (wc *WishClient) ReserveItemOfUser(item Item, user User) error {
	_, err := wc.executeRequest("PUT", fmt.Sprint("/user/", strconv.FormatUint(user.ID, 10), "/list/", strconv.FormatUint(item.ItemID, 10), "/reserve"), nil, nil, true)
	return err
}

// Reserves an item using the reserve endpoint
func (wc *WishClient) UnreserveItemOfUser(item Item, user User) error {
	_, err := wc.executeRequest("PUT", fmt.Sprint("/user/", strconv.FormatUint(user.ID, 10), "/list/", strconv.FormatUint(item.ItemID, 10), "/unreserve"), nil, nil, true)
	return err
}

// Converts the item to a short string for debugging (Doesn't contain all info)
func (i Item) String() string {
	return fmt.Sprintf("[%v](%v) %v - CHF %v", i.ItemID, i.Status.Name, i.Name, i.Price)
}
