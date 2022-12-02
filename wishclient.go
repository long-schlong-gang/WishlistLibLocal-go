package wishlistlib

import (
	"encoding/json"
	"io/ioutil"
)

type WishClient struct {
	datafile   string
	Users      map[uint64]JSONUser   `json:"users"`
	NextUserID uint64                `json:"new_user_id"`
	Items      map[uint64]JSONItem   `json:"items"`
	NextItemID uint64                `json:"new_item_id"`
	Links      map[uint64]JSONLink   `json:"links"`
	NextLinkID uint64                `json:"new_link_id"`
	Statuses   map[uint64]JSONStatus `json:"statuses"`
	token      Token
}

type JSONUser struct {
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	PasswordHash string   `json:"phash"`
	List         []uint64 `json:"list"`
}

type JSONItem struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Price          uint64   `json:"price"`
	ReservedByUser int64    `json:"reserved_by"`
	Status         uint64   `json:"status"`
	Links          []uint64 `json:"links"`
}

type JSONLink struct {
	Text string `json:"text"`
	URL  string `json:"hyperlink"`
}

type JSONStatus struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Returns the default context
func DefaultWishClient(DataFile string) (*WishClient, error) {
	data, err := ioutil.ReadFile(DataFile)
	if err != nil {
		return nil, err
	}

	var wc WishClient
	wc.datafile = DataFile
	err = json.Unmarshal(data, &wc)
	return &wc, err
}

func (wc *WishClient) Close() error {
	data, err := json.MarshalIndent(wc, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(wc.datafile, data, 0650)
}
