package brest_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/aptogeo/brest"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type Todo struct {
	ID   int `bun:",pk,autoincrement"`
	Text string
}

type Book struct {
	ID             int `bun:",pk,autoincrement"`
	Title          string
	NbPages        int
	AuthorID       int
	Author         *Author `bun:"rel:belongs-to,join:author_id=id"`
	TransientField string  `bun:"-"`
}

type Author struct {
	ID             int `bun:",pk,autoincrement"`
	Firstname      string
	Lastname       string
	Picture        []byte  `bun:",type:bytea"`
	Books          []*Book `bun:"rel:has-many,join:id=author_id"`
	TransientField string  `bun:"-"`
}

func AuthorBeforeHook(ctx context.Context, restQuery *brest.RestQuery, entity interface{}) error {
	log.Println("AuthorBeforeHook", restQuery, entity)
	return nil
}

func AuthorAfterHook(ctx context.Context, restQuery *brest.RestQuery, entity interface{}) error {
	log.Println("AuthorAfterHook", restQuery, entity)
	return nil
}

type PageOnly struct {
	NbPages int
}

func initTests(t *testing.T) (*bun.DB, *brest.Config) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	config := brest.NewConfig("/rest/", db)
	config.AddResource(brest.NewResource("Todo", (*Todo)(nil), brest.Get|brest.Post))
	config.AddResource(brest.NewResourceWithHooks("Author", (*Author)(nil), brest.All, AuthorBeforeHook, AuthorAfterHook))
	config.AddResource(brest.NewResource("Book", (*Book)(nil), brest.All))
	db.ResetModel(context.Background(), (*Todo)(nil))
	db.ResetModel(context.Background(), (*Author)(nil))
	db.ResetModel(context.Background(), (*Book)(nil))
	return db, config
}

var todos = []Todo{{
	Text: "Todo1",
}, {
	Text: "Todo2",
}}

var authors = []Author{{
	Firstname: "Antoine",
	Lastname:  "de Saint Exupéry",
}, {
	Firstname: "Franz",
	Lastname:  "Kafka",
}, {
	Firstname: "Francis Scott Key",
	Lastname:  "Fitzgerald",
}}

var books = []Book{{
	Title:    "Courrier sud",
	AuthorID: 1,
}, {
	Title:    "Vol de nuit",
	AuthorID: 1,
}, {
	Title:    "Terre des hommes",
	AuthorID: 1,
}, {
	Title:    "Lettre à un otage",
	AuthorID: 1,
}, {
	Title:    "Pilote de guerre",
	AuthorID: 1,
}, {
	Title:    "Le Petit Prince",
	AuthorID: 1,
}, {
	Title:    "La Métamorphose",
	AuthorID: 2,
}, {
	Title:    "La Colonie pénitentiaire",
	AuthorID: 2,
}, {
	Title:    "Le Procès",
	AuthorID: 2,
}, {
	Title:    "Le Château",
	AuthorID: 2,
}, {
	Title:    "L'Amérique",
	AuthorID: 2,
}, {
	Title:    "Gatsby le Magnifique",
	AuthorID: 3,
}}
