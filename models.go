package main

import "gorm.io/gorm"

const (
	PhotoKindOriginal   = "ORIGINAL"
	PhotoKindRoughTuned = "ROUGH"
	PhotoKindFineTuned  = "FINE"
)

type Event struct {
	gorm.Model
	Name  string `gorm:""`
	Girls []Girl `gorm:"foreignKey:EventID"`
}

type Girl struct {
	gorm.Model
	EventID   uint `gorm:"index"`
	AvatarURL string
	Token     string
	Photos    []Photo `gorm:"foreignKey:GirlID"`
}

type Photo struct {
	gorm.Model
	GirlID      uint   `gorm:"index"`
	Kind        string `gorm:"index"`
	URL         string
	DownloadURL string
}
