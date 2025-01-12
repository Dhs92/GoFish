package main

import (
	"github.com/Dhs92/GoFish/db"
)

func main() {
	test(&db.PostgresCon{})
}

func test(db db.DatabaseCon) error {
	db.CreateUser("Dude", "Dude@Dude.com", "Hunter2")

	return nil
}
