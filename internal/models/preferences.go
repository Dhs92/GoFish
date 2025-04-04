package models

type Preferences struct {
	Theme        string `json:"theme" gorm:"default:'light'"`
	Language     string `json:"language" gorm:"default:'en'"`
	DefaultUnits string `json:"default_units" gorm:"default:'imperial'"`
	Timezone     string `json:"timezone" gorm:"default:'UTC'"`
	Notification bool   `json:"notification" gorm:"default:true"`
}

func DefaultPreferences() Preferences {
	return Preferences{
		Theme:        "light",
		Language:     "en",
		DefaultUnits: "imperial",
		Timezone:     "UTC",
		Notification: true,
	}
}
