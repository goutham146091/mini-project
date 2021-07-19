package dbconnect

import (
	"fmt"
	"time"

	"git.npcompete.com/OSPSM_Servers/src/api/encrypt_decrypt"
	"git.npcompete.com/OSPSM_Servers/src/view"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

/*DB variable is used as common DB connection variable*/
var DB *gorm.DB
var dbUser, dbPassword, dbHost, dbPort, dbName, adminUserName, adminPassword string

func init() {
	// amcDBUser = os.Getenv("AMC_DB_USER")
	// amcDBPassword = os.Getenv("AMC_DB_PASSWORD")
	// amcDBHost = os.Getenv("AMC_DB_HOST")
	// amcDBPort = os.Getenv("AMC_DB_PORT")
	// amcDBName = os.Getenv("AMC_DB_NAME")
	// amcAdminUserName = os.Getenv("AMC_ADMIN_USERNAME")
	// amcAdminPassword = os.Getenv("AMC_ADMIN_PASSWORD")

	dbUser = "root"
	dbPassword = "ameex"
	dbHost = "localhost"
	dbPort = "3306"
	dbName = "mp_db1"
	adminUserName = "gouthamminiproject@mailinator.com"
	adminPassword = "goutham@123"
}

/*Connection is used connect the Db*/
func DbConnection() (dbConn *gorm.DB, err error) {
	if dbUser == "" {
		return nil, errors.New("DB connection : DB Username cannot be empty")
	}
	if dbPassword == "" {
		return nil, errors.New("DB connection : DB Password cannot be empty")
	}
	if dbHost == "" {
		return nil, errors.New("DB connection : DB Host cannot be empty")
	}
	if dbPort == "" {
		return nil, errors.New("DB connection : DB Port cannot be empty")
	}
	if dbName == "" {
		return nil, errors.New("DB connection : DB Name cannot be empty")
	}
	if adminUserName == "" {
		return nil, errors.New("DB connection : DB Admin username cannot be empty")
	}
	if adminPassword == "" {
		return nil, errors.New("DB connection : DB Admin password cannot be empty")
	}
	fmt.Println("connecting local DB")

	dbConn, err = gorm.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("db failed")
		fmt.Println("Error connection:" + err.Error())
	}
	fmt.Println("DB connected Successfully")
	dbConn.AutoMigrate(&view.Users{})
	dbConn.AutoMigrate(&view.Role{})
	dbConn.AutoMigrate(&view.Client{})
	dbConn.AutoMigrate(&view.ContactAddress{})
	dbConn.AutoMigrate(&view.ContactNumber{})
	dbConn.AutoMigrate(&view.Branch{})
	dbConn.AutoMigrate(&view.Product{})
	dbConn.AutoMigrate(&view.Invoice{})
	dbConn.Model(&view.Client{}).AddForeignKey("user_id", "contact_numbers(user_id)", "CASCADE", "CASCADE")
	dbConn.Model(&view.Client{}).AddForeignKey("user_id", "contact_addresses(user_id)", "CASCADE", "CASCADE")
	fmt.Println("table created successfully")
	currentTime := time.Now()
	formatedTime := currentTime.Format(time.RFC1123)

	key := "b826f19a79c492a6821233fdf54c5e923a5b839a357b024c62726821a02ba9f9"

	encryptpassword := encrypt_decrypt.EncryptData(adminPassword, key)

	fmt.Println(encryptpassword)

	roleName := [3]string{"Admin", "Client", "Branch"}
	for _, n := range roleName {
		dbConn.Create(&view.Role{
			RoleName: n,
		})
	}

	dbConn.Create(&view.Users{
		Email:          adminUserName,
		Password:       encryptpassword,
		Role:           "Admin",
		CreateDateTime: formatedTime,
	})
	DB = dbConn
	DB.LogMode(true)
	return
}

// func GenerateHashPassword(password string) string {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	encryptedPassword := string(hash)
// 	return encryptedPassword
// }
