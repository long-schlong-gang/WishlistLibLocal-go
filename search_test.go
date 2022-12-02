package wishlistlib

import "testing"

func TestSearch(t *testing.T) {

	wc, err := DefaultWishClient(DATA_FILE)
	t.Log("  Error: ", err)

	users, err := wc.SearchUsers("jOSEF")
	t.Log("  Error: ", err)
	t.Logf("Search Results: %v\n", users)

	err = wc.Close()
	t.Log("  Error: ", err)
}
