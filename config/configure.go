package config

type DbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Database string `json:"database"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

