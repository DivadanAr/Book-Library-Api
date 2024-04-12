package migration

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"main.go/models/entity"
	"os"
)

func Migrate() {
	log.Println("Start Migration Database ...")
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOSTNAME")+":3306)/"+os.Getenv("DB_NAME")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:  "",
			NameReplacer: nil,
			NoLowerCase:  true,
		},
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(
		&entity.Users{},
		&entity.Roles{},
		&entity.UserRoles{},
		&entity.Buku{},
		&entity.Kategori{},
		&entity.KategoriBuku{},
		&entity.UlasanBuku{},
		&entity.KoleksiPribadi{},
		&entity.Peminjaman{},
	)

	// &entities.Site{}
	if err != nil {
		log.Println(err)
	}
	log.Println("Database Migrated")

	defer func() {
		dbInstance, _ := db.DB()
		dbInstance.Close()
	}()
}
