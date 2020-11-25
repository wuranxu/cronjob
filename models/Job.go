package models

import "github.com/jinzhu/gorm"

type Job struct {
	gorm.Model
	Name     string `gorm:"type:varchar(32);not null" json:"name"`
	CronExpr string `gorm:"type:varchar(16);not null" json:"cron_expr"`
	Command  string `gorm:"type:text;not null" json:"command"`
}
