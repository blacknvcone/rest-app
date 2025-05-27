package transaction

import (
	"context"

	"rest-app/config/db"
)

// GetTrxContext returns default db instance where context doesn't contains custom db instance
func GetTrxContext(c context.Context, defaultDB *db.GormDB) *db.GormDB {
	tx, ok := c.Value(KeyTransaction).(*db.GormDB)
	if !ok {
		return &db.GormDB{
			DB: defaultDB.DB.WithContext(c),
		}
	}

	return &db.GormDB{
		DB: tx.DB.WithContext(c),
	}
}
