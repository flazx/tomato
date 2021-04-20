package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/test"
	_ "github.com/lib/pq" // postgres driver
	"gopkg.in/mgo.v2"
)

// OpenMongoDB 打开 MongoDB
func OpenMongoDB() *mgo.Database {
	// 此处仅用于测试
	if config.TConfig.DatabaseURI == "" {
		config.TConfig.DatabaseURI = test.MongoDBTestURL
	}

	session, err := mgo.Dial(config.TConfig.DatabaseURI)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	if err != nil {
		fmt.Println("mgo.Dial-error:", err)
		os.Exit(0)
	}
	session.SetMode(mgo.Eventual, true)
	myDB := session.DB("test") //这里的关键是连接mongodb后，选择admin数据库，然后登录，确保账号密码无误之后，该连接就一直能用了
	//出现server returned error on SASL authentication step: Authentication failed. 这个错也是因为没有在admin数据库下登录
	err = myDB.Login(config.TConfig.DatabaseUserName, config.TConfig.DatabaseUserPassword)
	if err != nil {
		fmt.Println("Login-error:", err)
		os.Exit(0)
	}
	//myDB = session.DB(mDBName) //如果要在这里就选择数据库，这个myDB可以定义为全局变量
	session.SetPoolLimit(10)

	return session.DB("test")
}

// OpenPostgreSQL 打开 PostgreSQL
func OpenPostgreSQL() *sql.DB {
	db, err := sql.Open("postgres", config.TConfig.DatabaseURI)
	if err != nil {
		panic(err)
	}
	return db
}
