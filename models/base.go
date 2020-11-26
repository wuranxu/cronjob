package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Model struct {
	ID        uint  `gorm:"primary_key" json:"id"`
	CreatedAt Time  `json:"created_at"`
	UpdatedAt Time  `json:"updated_at"`
	DeletedAt *Time `sql:"index" json:"deleted_at"`
}

const formatTime = "2006-01-02 15:04:05"

type Time struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t Time) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(formatTime))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan value of time.Time
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *Time) UnmarshalJSON(p []byte) error {
	var timeOrigin string
	err := json.Unmarshal(p, &timeOrigin)
	if err != nil {
		return err
	}
	tm, err := time.Parse(formatTime, timeOrigin)
	if err != nil {
		return err
	}
	*t = Time{Time: tm}
	return nil
}
