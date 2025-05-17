package mysql

import (
	"context"
	"fmt"
	"os"

	mysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"movieexample.com/metadata/pkg/model"
)

type Repository struct {
	db *gorm.DB
}

func New() (*Repository, error) {
	mysqlUser := os.Getenv("MYSQL_USER")
    mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDb := os.Getenv("MYSQL_DATABASE")
	dbUrl := os.Getenv("DATABASE_URL")


	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlUser, mysqlPassword, dbUrl, mysqlDb)
	db, err := gorm.Open(mysqldriver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Metadata{})

	return &Repository{db}, nil
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var metadata model.Metadata
	result := r.db.Where("id = ?", id, &metadata)
	if result.Error != nil {
		return nil, result.Error
	}

	return &metadata, nil
}

func (r *Repository) Put(ctx context.Context, id string, m *model.Metadata) error {
	newMetadata := model.Metadata { 
		Title: m.Title, 
		Description: m.Description, 
		Director: m.Director,
	}
	
	result := r.db.Create(&newMetadata)
	return result.Error
}
