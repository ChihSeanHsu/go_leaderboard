package repository

import "gorm.io/gorm"

type ScoreORM struct {
	gorm.Model
	ClientID string `gorm:"uniqueIndex:uni_idx__client_id;type:varchar(256)"`
	Score    int64  `gorm:"index:idx__score,sort:desc"`
}

func (ScoreORM) TableName() string {
	return "scores"
}