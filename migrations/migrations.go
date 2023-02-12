package migrations

import (
	"learn-go/helpers"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance uint
	UserId  uint
}

func connectDB() *gorm.DB {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	dsn, _ := viper.Get("DB_URL").(string)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	helpers.HandleErr(err)
	return db
}

func createAccounts() {
	db := connectDB()
	users := [...]User{
		{Username: "test1", Email: "test1@test1.com"},
		{Username: "test2", Email: "test2@test2.com"},
	}

	for i := 0; i < len(users); i++ {
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))
		user := User{
			Username: users[i].Username,
			Email:    users[i].Email,
			Password: generatedPassword,
		}
		db.Create(&user)

		account := Account{
			Type:    "Daily Account",
			Name:    string(users[i].Username + "'s account"),
			Balance: uint(10000 * int(i+1)),
			UserId:  user.ID,
		}
		db.Create(&account)
	}
}

func Migrate() {
	db := connectDB()
	db.AutoMigrate(&User{}, &Account{})

	createAccounts()
}
