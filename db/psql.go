package db

import "fmt"

type PostgresCon struct{}

// AddInventoryItem implements DatabaseCon.
func (psql *PostgresCon) AddInventoryItem() error {
	panic("unimplemented")
}

// AddLivestock implements DatabaseCon.
func (psql *PostgresCon) AddLivestock() error {
	panic("unimplemented")
}

// AddTank implements DatabaseCon.
func (psql *PostgresCon) AddTank() error {
	panic("unimplemented")
}

// DeleteInventoryItem implements DatabaseCon.
func (psql *PostgresCon) DeleteInventoryItem() error {
	panic("unimplemented")
}

// DeleteLivestock implements DatabaseCon.
func (psql *PostgresCon) DeleteLivestock() error {
	panic("unimplemented")
}

// DeleteTank implements DatabaseCon.
func (psql *PostgresCon) DeleteTank() error {
	panic("unimplemented")
}

// DeleteUser implements DatabaseCon.
func (psql *PostgresCon) DeleteUser(uuid string) error {
	panic("unimplemented")
}

// ModifyIntenvoryItem implements DatabaseCon.
func (psql *PostgresCon) ModifyIntenvoryItem() error {
	panic("unimplemented")
}

// ModifyLivestock implements DatabaseCon.
func (psql *PostgresCon) ModifyLivestock() error {
	panic("unimplemented")
}

// ModifyTank implements DatabaseCon.
func (psql *PostgresCon) ModifyTank() error {
	panic("unimplemented")
}

// ModifyUser implements DatabaseCon.
func (psql *PostgresCon) ModifyUser() error {
	panic("unimplemented")
}

func (psql *PostgresCon) CreateUser(name string, email string, password string) error {
	fmt.Printf("%s %s %s", name, email, password)
	return nil
}
