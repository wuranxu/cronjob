package models

type Job struct {
	Model
	Name     string `gorm:"type:varchar(32);not null;unique" binding:"required" json:"name"`
	Desc     string `gorm:"type:varchar(200);not null;comment: '描述'" binding:"required" json:"desc"`
	CronExpr string `gorm:"type:varchar(16);not null" binding:"required" json:"cron_expr"`
	Command  string `gorm:"type:text;not null" binding:"required" json:"command"`
	Retry    uint   `gorm:"type:int;not null;default 0;comment: '重试次数'" json:"retry"`
	Type     uint   `gorm:"type:int;not null;default 0;comment:'0: shell'" json:"type"`
}