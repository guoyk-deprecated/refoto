package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	PhotoKindOriginal   = "ORIGINAL"
	PhotoKindRoughTuned = "ROUGH"
	PhotoKindFineTuned  = "FINE"
)

type Event struct {
	gorm.Model
	Name  string
	Girls []Girl `gorm:"foreignKey:EventID"`
}

type Girl struct {
	gorm.Model
	EventID    uint `gorm:"index"`
	AvatarPath string
	Token      string
	Photos     []Photo `gorm:"foreignKey:GirlID"`
}

func (g Girl) AvatarURL() string {
	return ossCombineURL(g.AvatarPath, ossSuffixAvatar)
}

func (g Girl) PhotosWithKind(kind string) (output []Photo) {
	for _, p := range g.Photos {
		if p.Kind == kind {
			output = append(output, p)
		}
	}
	return
}

type Photo struct {
	gorm.Model
	GirlID uint   `gorm:"index"`
	Kind   string `gorm:"index"`
	Path   string
	Size   int64 `gorm:"index"`
}

func (g Girl) PreviewURL() string {
	return ossCombineURL(g.AvatarPath, ossSuffixPreview)
}

func (g Girl) URL() string {
	return ossCombineURL(g.AvatarPath, "")
}

func setupDB() (db *gorm.DB, err error) {
	if db, err = gorm.Open(mysql.Open(envMySQLDSN), &gorm.Config{}); err != nil {
		return
	}
	if envDebug {
		db.Logger = db.Logger.LogMode(logger.Info)
	}
	if err = db.AutoMigrate(&Event{}, &Girl{}, &Photo{}); err != nil {
		return
	}
	return
}
