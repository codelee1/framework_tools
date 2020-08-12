package orm

import (
	"time"
	"xorm.io/core"

	"github.com/go-xorm/xorm"
	// database driver
	_ "github.com/go-sql-driver/mysql"
)

type MyORM struct {
	DB *xorm.Engine
}

type Config struct {
	DSN         string        // data source name
	Active      int           // pool
	Idle        int           // pool
	IdleTimeout time.Duration // connect max life time
	Debug       bool
	Prefix      string // Prefix
}

func NewMySQL(c *Config) *MyORM {
	db, err := xorm.NewEngine("mysql", c.DSN)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	if c.Prefix != "" {
		pfMapper := core.NewPrefixMapper(core.SnakeMapper{}, c.Prefix)
		db.SetTableMapper(pfMapper)
	}
	db.SetConnMaxLifetime(c.IdleTimeout * time.Second)
	if c.Debug {
		db.ShowSQL(true)
	}
	orm := &MyORM{
		DB: db,
	}

	return orm
}

// Close close the Engine
func (m *MyORM) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}

func (m *MyORM) NewSession() *xorm.Session {
	return m.DB.NewSession()
}
