package models

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Ad struct {
	gorm.Model
	Description          string `gorm:"index"`
	Price                uint   `gorm:"index"`
	Title                string `gorm:"index"`
	UserID               uint   `gorm:"index"`
	FrameSizeID          uint   `gorm:"index"`
	FrameSizeDescription string
	BikeTypeID           uint    `gorm:"index"`
	Weight               float64 `gorm:"index"`
	Condition            bool    `gorm:"index"`
	FrameMaterialID      uint    `gorm:"index"`
	AdTypeID             uint    `gorm:"index"`
	AdDirectionID        uint    `gorm:"index"`
	LocationID           uint    `gorm:"index"`
	LocationDescription  string
	Open                 bool      `gorm:"index"`
	OpenUntil            time.Time `gorm:"index"`
	PhoneNo              string
	ValidationErrors     map[string]string `gorm:"-"`
}

func (ad *Ad) GetAdValueById(input map[uint]string, AdValueID uint) string {
	return input[AdValueID]
}

func (ad *Ad) Validate() bool {
	ad.ValidationErrors = make(map[string]string)

	if strings.TrimSpace(ad.Title) == "" {
		ad.ValidationErrors["Title"] = "Palun sisestage kuulutuse pealkiri"
	}

	if ad.AdDirectionID == 0 {
		ad.ValidationErrors["AdDirectionID"] = "Palun valige kuulutuse tüüp"
	}

	if ad.AdTypeID == 0 {
		ad.ValidationErrors["AdTypeID"] = "Palun valige kuulutuse kategooria"
	}

	if ad.AdTypeID == 1 && ad.BikeTypeID == 0 {
		ad.ValidationErrors["BikeTypeID"] = "Palun valige ratta tüüp"
	}

	if ad.AdTypeID == 1 && ad.FrameSizeID == 0 {
		ad.ValidationErrors["FrameSizeID"] = "Palun valige raami suurus"
	}

	if ad.AdTypeID == 1 && ad.FrameMaterialID == 0 {
		ad.ValidationErrors["FrameMaterialID"] = "Palun valige raami materjal"
	}

	if ad.Price == 0 {
		ad.ValidationErrors["Price"] = "Palun sisestage hind"
	}

	if ad.LocationID == 0 {
		ad.ValidationErrors["LocationID"] = "Palun valige asukoht"
	}

	return len(ad.ValidationErrors) == 0
}
