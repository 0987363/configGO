package middleware

import (
	//	"net/http"

	//	"golang.org/x/net/context"

	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"
)

const dbKey = "Db"

var db *mgo.Session

// ConnectDB connects to the database and store the db pointer for DbConnector
// middleware.
func ConnectDB(dataURL string) (err error) {
	db, err = mgo.Dial(dataURL)
	if err == nil {
		//		db.SetMode(mgo.Eventual, true)
	}
	return
}

// DBConnector middleware stores a mongo db handler in context
func DBConnector() gin.HandlerFunc {
	return func(c *gin.Context) {
		d := db.Copy()
		d.SetMode(mgo.Monotonic, true)
		c.Set(dbKey, d)
		c.Next()
		defer d.Close()
	}
}

// GetDB returns the db pointer from context or nil if db has not been connected
func GetDB(c *gin.Context) *mgo.Session {
	if db, ok := c.Get(dbKey); ok {
		return db.(*mgo.Session)
	}

	return nil
}
