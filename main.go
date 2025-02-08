/*
Copyright © 2025 Alex Forehand <forehand.alex@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"context"

	"github.com/Dhs92/GoFish/db"
)

func main() {
	ctx := context.Background()

	// Create a new database connection
	database, err := db.Connect(ctx, "mongodb://localhost:27017", "aquarium")

	if err != nil {
		panic(err)
	}

	// Create a new user
	user, err := db.NewUser("Test User", "test@email.com", "password")

	if err != nil {
		panic(err)
	}

	// Insert the user into the database
	_, err = database.Create(ctx, user)

	if err != nil {
		panic(err)
	}
}
