package mysql

import (
	"context"
	"fmt"
	"os"

	mysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
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

	db.AutoMigrate(&model.Rating{}, &model.RatingEvent{})

	return &Repository{db}, nil
}

func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	var ratings []model.Rating
	result := r.db.Where("record_id = ? AND record_type = ?", recordID, recordType, &ratings)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}

	return ratings, nil
}

func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	ratingToSave := model.Rating{
		RecordType: string(recordType),
		RecordID: string(recordID),
		UserID: rating.UserID,
		Value: rating.Value,
	}
	result := r.db.Create(&ratingToSave)
	return result.Error
}
