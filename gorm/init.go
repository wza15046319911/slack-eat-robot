package gorm

import (
	"context"
	"eat-and-go/config"
	"eat-and-go/model"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/mysql"
	"istio.io/pkg/log"
	"time"
	// MySQL driver.
	"gorm.io/gorm"
)

type Database struct {
	Self        *gorm.DB
	Docker      *gorm.DB
	MongoClient *mongo.Client
}
type MongoCollection struct {
	RegoCollection *mongo.Collection
}

var Collections *MongoCollection
var DB *Database

// NewMongoDB 创建新的MongoDB实例
func NewMongoDB() *mongo.Client {
	mongoAddress := config.GetConfig().Mongo.Address
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoAddress))
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	fmt.Println("Ping MongoDB: ", mongoAddress)
	return client
}
//NewGormDB 初始化新的MySQL数据库
func NewGormDB() *gorm.DB {
	databaseConfig := config.GetConfig().Gorm
	dsn := config.GetConfig().MySql.DSN()
	db, err := newDatabase(
		&model.GormConfig{
			Debug:        databaseConfig.Debug,
			DSN:          dsn,
			MaxLifetime:  databaseConfig.MaxLifetime,
			MaxOpenConns: databaseConfig.MaxOpenConns,
			MaxIdleConns: databaseConfig.MaxIdleConns,
			//TablePrefix:  databseConfig.TablePrefix,
		})
	if err != nil {
		panic(err)
		//log.Errorf(err, "Database connection failed. Database name: %s", db.Name())
		//return nil
	}
	err = db.AutoMigrate(&model.SlackRecommendation{})
	if err != nil {
		panic(err)
	}
	//err = db.AutoMigrate(&model.GopaRoles{})
	//if err != nil {
	//	panic(err)
	//}
	//err = db.AutoMigrate(&model.GopaMembers{})
	//if err != nil {
	//	panic(err)
	//}
	//err = db.AutoMigrate(&model.GopaProjectRoles{})
	//if err != nil {
	//	panic(err)
	//}
	//err = db.AutoMigrate(&model.GopaApplication{})
	//if err != nil {
	//	panic(err)
	//}
	//err = db.AutoMigrate(&model.ProjectResources{})
	//if err != nil {
	//	panic(err)
	//}
	return db
}

func newDatabase(c *model.GormConfig) (*gorm.DB, error) {
	fmt.Println("Pinging database: ", c.DSN)
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if c.Debug {
		db = db.Debug()
	}
	return db, nil
}

func GetSelfDB() *gorm.DB {
	return NewGormDB()
}

func GetCollection(database string, collectionName string) *mongo.Collection {
	mongoClient := DB.MongoClient
	collection := mongoClient.Database(database).Collection(collectionName)
	return collection
}

func (db *Database) Init() {
	log.Info("Init DB.")
	DB = &Database{
		Self:        GetSelfDB(),
		MongoClient: NewMongoDB(),
		//		Docker: GetDockerDB(),
	}
	var database = config.GetConfig().Mongo.Database
	var collection = config.GetConfig().Mongo.Collection
	Collections = &MongoCollection{
		RegoCollection: GetCollection(database, collection),
	}
}
