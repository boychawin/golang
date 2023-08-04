package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID        uint   `json:"id"`
	UserName  string `gorm:"type:varchar(100);default:null" json:"UserName"`
	PassWord  string `gorm:"type:varchar(255);default:null"  json:"PassWord"`
	FirstName string `gorm:"type:varchar(255);default:null" json:"FirstName"`
	LastName  string `gorm:"type:varchar(255);default:null"  json:"LastName"`
	Role      string         `gorm:"size:20;default:client;comment:client,admin" json:"Role"`
	Status    int            `gorm:"size:100;comment:0/1;default:1" json:"Status"`
	CreatedAt time.Time      `gorm:"comment:วันที่สร้าง;" json:"CreatedAt"`
	UpdatedAt time.Time      `gorm:"comment:วันที่แก้ไข;default:null" json:"UpdatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
