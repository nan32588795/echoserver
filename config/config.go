package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type PostgreSQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type MySQLConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DBName    string `json:"dbname"`
	ParseTime string `json:"parseTime"`
}

type Config struct {
	ActiveDB         string           `json:"activeDB"`
	PostgreSQLConfig PostgreSQLConfig `json:"postgreSQL"`
	MySQLConfig      MySQLConfig      `json:"mySQL"`
}

func init() {
	configPath := "config/config.json"
	var err error
	GlobalConfig, err = LoadConfig(configPath)
	if err != nil {
		log.Fatal("設定ファイルの読み込み失敗:", err)
	}
	// JSONで整形して出力
	prettyCfg, err := json.MarshalIndent(maskConfig(GlobalConfig), "", "  ")
	if err != nil {
		log.Println("設定の出力整形に失敗:", err)
	} else {
		fmt.Println("Loaded config:")
		fmt.Println(string(prettyCfg))
	}
}

var GlobalConfig *Config

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	return &cfg, err
}

func (c *Config) DriverName() string {
	switch c.ActiveDB {
	case "PostgreSQL":
		return "postgres"
	case "MySQL":
		return "mysql"
	default:
		panic("unsupported database type")
	}
}

func (c *Config) ConnString() string {
	switch c.ActiveDB {
	case "PostgreSQL":
		return c.PostgreSQLConfig.connString()
	case "MySQL":
		return c.MySQLConfig.connString()
	default:
		panic("unsupported database type")
	}
}

func (c *PostgreSQLConfig) connString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func (c *MySQLConfig) connString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.ParseTime)
}

func maskConfig(cfg *Config) any {
	return struct {
		ActiveDB         string `json:"activeDB"`
		PostgreSQLConfig struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			User     string `json:"user"`
			Password string `json:"password"`
			DBName   string `json:"dbname"`
			SSLMode  string `json:"sslmode"`
		} `json:"postgreSQL"`
		MySQLConfig struct {
			Host      string `json:"host"`
			Port      int    `json:"port"`
			User      string `json:"user"`
			Password  string `json:"password"`
			DBName    string `json:"dbname"`
			ParseTime string `json:"parseTime"`
		} `json:"mySQL"`
	}{
		ActiveDB: cfg.ActiveDB,
		PostgreSQLConfig: struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			User     string `json:"user"`
			Password string `json:"password"`
			DBName   string `json:"dbname"`
			SSLMode  string `json:"sslmode"`
		}{
			Host:     cfg.PostgreSQLConfig.Host,
			Port:     cfg.PostgreSQLConfig.Port,
			User:     cfg.PostgreSQLConfig.User,
			Password: "****",
			DBName:   cfg.PostgreSQLConfig.DBName,
			SSLMode:  cfg.PostgreSQLConfig.SSLMode,
		},
		MySQLConfig: struct {
			Host      string `json:"host"`
			Port      int    `json:"port"`
			User      string `json:"user"`
			Password  string `json:"password"`
			DBName    string `json:"dbname"`
			ParseTime string `json:"parseTime"`
		}{
			Host:      cfg.MySQLConfig.Host,
			Port:      cfg.MySQLConfig.Port,
			User:      cfg.MySQLConfig.User,
			Password:  "****",
			DBName:    cfg.MySQLConfig.DBName,
			ParseTime: cfg.MySQLConfig.ParseTime,
		},
	}
}
