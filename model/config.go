package model

import "fmt"
type GormConfig struct {
	Debug        bool
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	//TablePrefix  string
}
type SysConfig struct {
	Service SysService   `yaml:"service"`
	Log     Logger       `yaml:"log"`
	Slack   Slack        `yaml:"slack" mapstructure:"slack"`
	Gorm       GormService  `yaml:"gorm" mapstructure:"gorm"`
	MySql      MySQL        `yaml:"mysql" mapstructure:"mysql"`
	Mongo   MongoService `yaml:"mongo" mapstructure:"mongo"`
}

type SysService struct {
	Runmode      string `yaml:"runmode"`
	Addr         string `yaml:"addr"`
	Name         string `yaml:"name"`
	Url          string `yaml:"url"`
	MaxPingCount int    `yaml:"maxPingCount"`
	JwtSecret    string `yaml:"jwtSecret"`
}

type Logger struct {
	Dir    string `yaml:"dir"`
	Remain int    `yaml:"remain"`
}

//Slack configuration for slack
type Slack struct {
	SlackAppToken string `yaml:"slackAppToken" mapstructure:"slackAppToken"`
	SlackBotToken string `yaml:"slackBotToken" mapstructure:"slackBotToken"`
	SigningSecret string `yaml:"signingSecret" mapstructure:"signingSecret"`
	AuthUser      string `yaml:"authUser"`
}

type MongoService struct {
	Address    string `yaml:"address" mapstructure:"address"`
	Database   string `yaml:"database" mapstructure:"database"`
	Collection string `yaml:"collection" mapstructure:"collection"`
}

// MySQL 读取MySQL部署的配置信息
type MySQL struct {
	Host       string `yaml:"host" mapstructure:"host"`
	Port       int    `yaml:"port" mapstructure:"port"`
	User       string `yaml:"user" mapstructure:"user"`
	Password   string `yaml:"password" mapstructure:"password"`
	DBName     string `yaml:"dbName" mapstructure:"dbName"`
	Parameters string `yaml:"parameters" mapstructure:"parameters"`
}

// DSN 数据库连接串
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

// GormService Gorm配置信息
type GormService struct {
	Debug             bool   `yaml:"debug" mapstructure:"debug"`
	DBType            string `yaml:"dbType" mapstructure:"dbType"`
	MaxLifetime       int    `yaml:"maxLifetime" mapstructure:"maxLifetime"`
	MaxOpenConns      int    `yaml:"maxOpenConns" mapstructure:"maxOpenConns"`
	MaxIdleConns      int    `yaml:"maxIdleConns" mapstructure:"maxIdleConns"`
	TablePrefix       string `yaml:"tablePrefix" mapstructure:"tablePrefix"`
	EnableAutoMigrate bool   `yaml:"enableAutoMigrate" mapstructure:"enableAutoMigrate"`
}