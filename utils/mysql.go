package utils

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"gorm.io/gorm/clause"
)

const (
	Limit       = 30
	MaxInterval = time.Hour * 24 * 10
)

type options struct {
	Tx         *gorm.DB
	Lock       bool
	FromMaster bool
}

type Options func(*options)

func defaultClause() options {
	return options{
		Tx:         nil,
		Lock:       false,
		FromMaster: false,
	}
}

func ConfigureDB(db *gorm.DB, clauses ...Options) *gorm.DB {
	q := defaultClause()

	for _, fn := range clauses {
		fn(&q)
	}

	chain := db
	if q.Tx != nil {
		chain = q.Tx
	}

	if q.FromMaster {
		chain = chain.Clauses(dbresolver.Write)
	}

	if q.Lock {
		chain = chain.Clauses(clause.Locking{
			Strength: "UPDATE",
			Table:    clause.Table{Name: clause.CurrentTable},
		})
	}

	return chain
}
