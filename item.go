package wishlistlib

import (
	"fmt"
	"strconv"
)

type Item struct {
	ItemID      uint64  `json:"item_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Status      Status  `json:"status"`
	Links       []Link  `json:"links"`
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
func (ctx *Context) GetAllItems(user User) ([]Item, error) {
	items := []Item{}
	err := ctx.parseObjectFromServer("/user/"+strconv.FormatUint(user.ID, 10)+"/list", "GET", &items, nil, false)
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
func (ctx *Context) GetItemByID(id uint64) (Item, error) {
	item := Item{}
	err := ctx.parseObjectFromServer("/item/"+strconv.FormatUint(id, 10), "GET", &item, nil, false)
	if err != nil {
		return Item{}, err
	}

	// Convert price
	item.Price = float32(item.Price) / 100

	return item, nil
}

// Adds an item to a user's list and returns the item with all its info
func (ctx *Context) AddItemToAuthenticatedUserList(item Item) (Item, error) {
	obj := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       uint64 `json:"price"`
		Status      struct {
			StatusID uint64 `json:"status_id"`
		} `json:"status"`
		Links []struct {
			Text      string `json:"text"`
			Hyperlink string `json:"hyperlink"`
		} `json:"links"`
	}{
		Name:        item.Name,
		Description: item.Description,
		Price:       uint64(item.Price * 100), // Converts to integer cents for server
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

	err := ctx.sendObjectToServer("/user/"+strconv.FormatUint(ctx.authUser.ID, 10)+"/list", "POST", obj, true)
	if err != nil {
		return Item{}, err
	}

	items, err := ctx.GetAllItems(ctx.authUser)
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
func (ctx *Context) DeleteItem(item Item) error {
	return ctx.simpleRequest("/user/"+strconv.FormatUint(ctx.authUser.ID, 10)+"/list/"+strconv.FormatUint(item.ItemID, 10), "DELETE", true)
}

// Converts the item to a short string for debugging (Doesn't contain all info)
func (i Item) String() string {
	return fmt.Sprintf("[%v](%v) %v - CHF %v", i.ItemID, i.Status.Name, i.Name, i.Price)
}
