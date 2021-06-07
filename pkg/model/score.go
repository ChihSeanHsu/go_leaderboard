package model

import (
	"gorm.io/gorm"
)

type Score struct {
	ClientID string `gorm:"varchar(256)"`
	Score    uint64 `gorm:"index""`
}

type ScoreModel struct {
	gorm.Model
	Score
}
