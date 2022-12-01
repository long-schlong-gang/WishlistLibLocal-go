package wishlistlib

import "testing"

func TestSearch(t *testing.T) {

	wc, err := DefaultWishClient(DATA_FILE)
	t.Logf("Err: %v\n", err)

	users, err := wc.SearchUsers("jOSEF")
	t.Logf("Search Results: %v\n", users)

	err = wc.Close()
	t.Logf("Err: %v\n", err)
}
