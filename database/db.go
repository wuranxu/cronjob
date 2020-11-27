package database

import (
	"cronjob/config"
	"cronjob/models"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"reflect"
)

var (
	UnSupportedDatabase = errors.New("database is not supported now")
	UpdateParamsError   = errors.New("you must provide update column and value")
	StructError         = errors.New("you must provide the same struct")
)

var (
	Conn *Cursor
)

//type Columns map[string]interface{}

func Use(cfg config.DbConfig) {
	var err error
	if Conn, err = NewConnect(cfg); err != nil {
		panic(err)
	}
	Conn.LogMode(cfg.LogMode)
	if cfg.AutoMigrate {
		Conn.AutoMigrate(models.Tables...)
	}
}

func NewConnect(cfg config.DbConfig) (cur *Cursor, err error) {
	var (
		args string
		db   *gorm.DB
		like string
	)
	switch cfg.Name {
	case "postgres":
		// postgres sql
		args = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password)
		like = "ilike"
	case "mysql":
		args = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		like = "like"
	default:
		err = UnSupportedDatabase
		return
	}
	db, err = gorm.Open(cfg.Name, args)
	if err != nil {
		return
	}
	return &Cursor{db, like}, nil
}

type Cursor struct {
	*gorm.DB
	like string
}

func (c *Cursor) Like(field, name string) *Cursor {
	c.DB.Error = nil
	c.DB = c.Where(fmt.Sprintf("%s %s '%%%s%%'", field, c.like, name))
	return c
}

func (c *Cursor) Llike(field, name string) *Cursor {
	c.DB.Error = nil
	c.DB = c.Where(fmt.Sprintf("%s %s '%%%s'", field, c.like, name))
	return c
}
func (c *Cursor) Rlike(field, name string) *Cursor {
	c.DB.Error = nil
	c.DB = c.Where(fmt.Sprintf("%s %s '%s%%'", field, c.like, name))
	return c
}

func (c *Cursor) Find(out interface{}, where ...interface{}) *Cursor {
	c.DB.Error = nil
	c.DB = c.DB.Find(out, where...)
	return c
}

// find with pagination
func (c *Cursor) FindPagination(page, pageSize int, out interface{}, where ...interface{}) (int, error) {
	c.DB.Error = nil
	var total int
	c.DB = c.Find(out, where...).Count(&total)
	err := c.Page(page, pageSize).Error
	return total, err
}

// find with order
func (c *Cursor) FindPaginationAndOrder(page, pageSize int, order string, out interface{}, where ...interface{}) (int, error) {
	var total int
	c.DB.Error = nil
	c.DB = c.Find(out, where...).Count(&total)
	err := c.Page(page, pageSize).Order(order).Error
	return total, err
}

func (c *Cursor) Table(name string) *Cursor {
	c.DB.Error = nil
	c.DB = c.DB.Table(name)
	return c
}

func (c *Cursor) Select(query interface{}, args ...interface{}) *Cursor {
	c.DB.Error = nil
	c.DB = c.DB.Select(query, args...)
	return c
}

func (c *Cursor) Sql(v interface{}, sql string, params ...interface{}) error {
	c.DB.Error = nil
	c.DB = c.Raw(sql, params...).Scan(v)
	return c.DB.Error
}

func (c *Cursor) Insert(v interface{}) error {
	c.DB.Error = nil
	c.DB = c.DB.Create(v)
	return c.DB.Error
}

func (c *Cursor) Delete(v interface{}, where ...interface{}) error {
	c.DB.Error = nil
	c.DB = c.DB.Delete(v, where...)
	return c.Error
}

func (c *Cursor) Updates(v interface{}, attrs ...interface{}) (int64, error) {
	c.DB.Error = nil
	switch len(attrs) {
	case 0:
		return 0, UpdateParamsError
	case 1:
		switch to := attrs[0].(type) {
		// receive map and struct
		case map[string]interface{}:
			c.DB = c.DB.Model(v).Updates(to)
		default:
			var dist string
			vType := reflect.ValueOf(v).Elem().Type().String()
			if value := reflect.ValueOf(to); value.Kind() == reflect.Ptr {
				// 指针模式
				dist = value.Elem().Type().String()
			} else {
				dist = value.Type().String()
			}
			if dist != vType {
				return 0, StructError
			}
			c.DB = c.DB.Model(v).Updates(to)
		}
	default:
		c.DB = c.DB.Model(v).Update(attrs...)
	}
	return c.DB.RowsAffected, c.DB.Error
}

func (c *Cursor) Page(current, pageSize int) *Cursor {
	c.DB.Error = nil
	c.DB = c.Offset((current - 1) * pageSize).Limit(pageSize)
	return c
}
