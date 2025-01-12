package db

type DatabaseCon interface {
	CreateUser(name string, email string, password string) error
	ModifyUser( /* TODO */ ) error
	DeleteUser(uuid string) error
	AddTank( /* TODO */ ) error
	ModifyTank( /* TODO */ ) error
	DeleteTank( /* TODO */ ) error
	AddLivestock( /* TODO */ ) error
	ModifyLivestock( /* TODO */ ) error
	DeleteLivestock( /* TODO */ ) error
	AddInventoryItem( /* TODO */ ) error
	ModifyIntenvoryItem( /* TODO */ ) error
	DeleteInventoryItem( /* TODO */ ) error
}
