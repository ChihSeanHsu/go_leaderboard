package model

import (
	"gorm.io/gorm"
)

type ScoreModel struct {
	gorm.Model
	ClientID string `gorm:"uniqueIndex:idx_name;type:varchar(256)"`
	Score    uint64 `gorm:"index"`
}

//type ScoreModel struct {
//
//	Score
//}
