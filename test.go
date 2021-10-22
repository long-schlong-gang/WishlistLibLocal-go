package main

import "fmt"

func main() {
	fmt.Println(" All Users:\n------------")

	ctx := DefaultContext()
	users, _ := ctx.GetAllUsers()

	for _, u := range users {
		fmt.Println(u)
	}

	user := users[0]

	fmt.Printf("\n All User [%v]'s Items:\n------------\n", user.ID)

	items, _ := ctx.GetAllItems(user)

	for _, i := range items {
		fmt.Println(i)
	}
}
