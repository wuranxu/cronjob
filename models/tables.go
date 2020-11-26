package models

var Tables = make([]interface{}, 0, 10)

func AutoMigrate(table ...interface{}) {
	Tables = append(Tables, table...)
}

func init() {
	AutoMigrate(new(Job))
}
