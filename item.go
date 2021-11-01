package wishlistlib

import (
	"fmt"
	"strconv"
)

type Item struct {
	ID          uint64 `json:"item_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Status      Status `json:"status"`
	Links       []Link `json:"links"`
}

type Link struct {
	ID   uint64 `json:"link_id"`
	Text string `json:"text"`
	URL  string `json:"hyperlink"`
}

type Status struct {
	ID          uint64 `json:"status_id"`
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

	return items, nil
}

// Converts the item to a short string for debugging (Doesn't contain all info)
func (i Item) String() string {
	return fmt.Sprintf("[%v](%v) %v - CHF %v", i.ID, i.Status.Name, i.Name, i.Price)
}
