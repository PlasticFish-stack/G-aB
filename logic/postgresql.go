package logic

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Gorm *gorm.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "dbuser"
	password = "asd$%^456"
	dbname   = "clientdb"
)

func ExpandUserGroup(db *gorm.DB) *gorm.DB {
	return db.Preload("Children", ExpandUserGroup)
}

func GetConn() {
	// GORM connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbname, port)

	var err error
	Gorm, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix:   "t_", // 表名前缀，`User`表为`t_users`
			// SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
			// NameReplacer:  strings.NewReplacer("CID", "Cid"), // 在转为数据库名称之前，使用NameReplacer更改结构/字段名称。
		},
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	sqlDB, err := Gorm.DB()
	if err != nil {
		log.Fatalf("Failed to get database: %v", err)
	}
	sqlDB.SetConnMaxLifetime(300)
	sqlDB.SetMaxIdleConns(10)
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
}

//	func GetConn() {
//		conn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", user, password, host, port, dbname)
//		var err error
//		DB, err = sql.Open("postgres", conn)
//		if err != nil {
//			fmt.Println(err.Error())
//		}
//		DB.SetConnMaxLifetime(100)
//		DB.SetMaxIdleConns(10)
//		err = DB.Ping()
//		if err != nil {
//			fmt.Println("Failed to connect to the database: ", err)
//		}
//	}
