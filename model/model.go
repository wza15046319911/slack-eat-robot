package model

import (
	"sync"
	"time"
)

// BaseModel config
type BaseModel struct {
	ID        uint64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"-"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"-"`
	DeletedAt *time.Time `gorm:"column:deletedAt" sql:"index" json:"-"`
}

// UserInfo config
type UserInfo struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	SayHello  string `json:"sayHello"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// UserList config
type UserList struct {
	Lock  *sync.Mutex
	IDMap map[uint64]*UserInfo
}

// Token represents a JSON web token.
type Token struct {
	Token string `json:"token"`
}

// KV key and valuse
type KV struct {
	Key   string
	Value string
}

// KVL key and values
type KVL struct {
	Key    string
	Values []string
}

// TsAndVersion config
type TsAndVersion struct {
	CreateTs int64 `xorm:"created"`
	UpdateTs int64 `xorm:"updated"`
	Version  int64 `xorm:"version"`
}

//DianpingItem 一条大众点评商铺信息结构体
type DianpingItem struct {
	ShopID        string `bson:"店铺id"`
	ShopName      string `bson:"店铺名"`
	TotalComments string `bson:"评论总数"`
	AvgPrice      string `bson:"人均价格"`
	TagOne        string `bson:"标签1"`
	TagTwo        string `bson:"标签2"`
	ShopAddress   string `bson:"店铺地址"`
	DetailLink    string `bson:"详情链接"`
	ImageLink     string `bson:"图片链接"`
	AvgRatings    struct {
		Appetite    string `bson:"口味"`
		Environment string `bson:"环境"`
		Service     string `bson:"服务"`
	} `bson:"店铺均分"`
	Recommend   string `bson:"推荐菜"`
	TotalRating string `bson:"店铺总分"`
	Phone       string `bson:"店铺电话"`
	Others      string `bson:"其他信息"`
	BonusInfo   string `bson:"优惠券信息"`
}

type SlackRecommendation struct {
	BaseModel
	UserId            string `gorm:"column:userId"`
	UserName          string `gorm:"column:userName"`
	RecommendShopId   string `gorm:"column:recommendShopId"`
	RecommendShopName string `gorm:"column:recommendShopName"`
	RecommendShopLink string `gorm:"column:recommendShopLink"`
	Preference        string `gorm:"column:preference"`
}

func (SlackRecommendation) TableName() string {
	return "slack_recommendation"
}
