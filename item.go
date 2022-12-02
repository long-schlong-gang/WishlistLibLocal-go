package wishlistlib

import (
	"fmt"
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
	user, err := wc.GetUserByID(user.ID)
	if err != nil {
		return nil, err
	}

	items := []Item{}
	for _, iid := range wc.Users[user.ID].List {
		if item, err := wc.GetItemByID(iid); err == nil {
			items = append(items, item)
		}
		fmt.Println()
	}

	return items, nil
}

// Retrieves an item with a certain ID
func (wc *WishClient) GetItemByID(id uint64) (Item, error) {
	var item *Item
	if ji, exists := wc.Items[id]; exists {
		// Check if reserved
		reservedby := User{}
		if ji.ReservedByUser >= 0 {
			if u, err := wc.GetUserByID(uint64(ji.ReservedByUser)); err == nil {
				reservedby = u
			}
		}

		status := Status{}
		if ji.Status > 0 {
			if s, exists := wc.Statuses[uint64(ji.Status)]; exists {
				status = Status{
					StatusID:    uint64(ji.Status),
					Name:        s.Name,
					Description: s.Description,
				}
			}
		}

		links := []Link{}
		for _, lid := range ji.Links {
			if l, exists := wc.Links[uint64(lid)]; exists {
				links = append(links, Link{
					LinkID: uint64(lid),
					Text:   l.Text,
					URL:    l.URL,
				})
			}
		}

		item = &Item{
			ItemID:         id,
			Name:           ji.Name,
			Description:    ji.Description,
			Price:          float32(ji.Price / 100), // Convert price from cents
			ReservedByUser: reservedby,
			Status:         status,
			Links:          links,
		}
	}
	if item == nil {
		return Item{}, NotFoundError(fmt.Sprintf("Item with ID '%d' not found", id))
	}

	return *item, nil
}

// Adds an item to a user's list and returns the item with all its info
func (wc *WishClient) AddItemOfUser(item Item, user User) (Item, error) {

	ju, exists := wc.Users[user.ID]
	if !exists {
		return Item{}, NotFoundError(fmt.Sprintf("User with ID '%d' not found", user.ID))
	}

	// Check price is within limits
	intPrice := uint32(item.Price * 100)
	if uint64(intPrice) != uint64(item.Price*100) {
		return Item{}, PriceOutOfRangeError(item.Price * 100)
	}

	// Set status to available if not provided
	if item.Status.StatusID == 0 || item.Status.StatusID > 3 {
		item.Status.StatusID = 1
	}

	// Add new item
	reservedByID := int64(-1)
	if item.ReservedByUser == (User{}) {
		reservedByID = int64(item.ReservedByUser.ID)
	}

	links := []uint64{}
	for _, link := range item.Links {
		links = append(links, link.LinkID)
	}

	ji := JSONItem{
		Name:           item.Name,
		Description:    item.Description,
		Price:          uint64(intPrice),
		ReservedByUser: reservedByID,
		Status:         item.Status.StatusID,
		Links:          links,
	}
	nid := wc.NextItemID
	wc.NextItemID++

	wc.Items[nid] = ji
	ju.List = append(ju.List, nid)
	wc.Users[user.ID] = ju

	return wc.GetItemByID(nid)
}

// Delete an item from the database (WARNING! PERMANENT!)
func (wc *WishClient) DeleteItemOfUser(item Item, user User) error {
	if _, exists := wc.Items[item.ItemID]; !exists {
		return NotFoundError(fmt.Sprintf("Item with ID '%d' not found", item.ItemID))
	}

	delete(wc.Items, item.ItemID)
	return nil
}

// Sets the status of the provided item
func (wc *WishClient) SetItemStatusOfUser(item Item, user User, status Status) error {
	ji, exists := wc.Items[item.ItemID]
	if !exists {
		return NotFoundError(fmt.Sprintf("Item with ID '%d' not found", item.ItemID))
	}

	ji.Status = status.StatusID
	wc.Items[item.ItemID] = ji

	return nil
}

// Reserves an item using the reserve endpoint
func (wc *WishClient) ReserveItemOfUser(item Item, user User) error {
	ji, exists := wc.Items[item.ItemID]
	if !exists {
		return NotFoundError(fmt.Sprintf("Item with ID '%d' not found", item.ItemID))
	}

	ji.ReservedByUser = int64(user.ID)
	wc.Items[item.ItemID] = ji

	return nil
}

// Reserves an item using the reserve endpoint
func (wc *WishClient) UnreserveItemOfUser(item Item, user User) error {
	ji, exists := wc.Items[item.ItemID]
	if !exists {
		return NotFoundError(fmt.Sprintf("Item with ID '%d' not found", item.ItemID))
	}

	ji.ReservedByUser = -1
	wc.Items[item.ItemID] = ji

	return nil
}

// Converts the item to a short string for debugging (Doesn't contain all info)
func (i Item) String() string {
	return fmt.Sprintf("[%05d](%v) %v - CHF %v", i.ItemID, i.Status.Name, i.Name, i.Price)
}
