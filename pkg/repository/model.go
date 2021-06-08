package repository

import "gorm.io/gorm"

type ScoreORM struct {
	gorm.Model
	ClientID string  `gorm:"uniqueIndex:uni_idx__client_id;type:varchar(256)"`
	Score    float64 `gorm:"index:idx__score,sort:desc;type:numeric(10,1)"`
}

func (ScoreORM) TableName() string {
	return "scores"
}
