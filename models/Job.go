package models

import "github.com/jinzhu/gorm"

type Job struct {
	gorm.Model
	Name     string `gorm:"type:varchar(32);not null" json:"name"`
	Desc     string `gorm:"type:varchar(200);not null;comment: '描述'" json:"desc"`
	CronExpr string `gorm:"type:varchar(16);not null" json:"cron_expr"`
	Command  string `gorm:"type:text;not null" json:"command"`
	Retry    uint   `gorm:"type:int;not null;default 0" json:"retry"`
	Type     uint   `gorm:"type:int;not null;default 0;comment:'0: shell'" json:"type"`
}
