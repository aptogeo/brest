package brest_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/aptogeo/brest"
	"github.com/stretchr/testify/assert"
)

func TestDeserialize(t *testing.T) {
	db, config := initTests(t)
	defer db.Close()
	engine := brest.NewEngine(config)

	var err error
	var content []byte
	var book *Book

	resource := engine.Config().GetResource("Book")

	book = &Book{}
	err = engine.Deserialize(&brest.RestQuery{Action: brest.Post, Resource: "Book", ContentType: "application/json", Content: []byte("{\"Title\":\"json title\",\"NbPages\":520,\"UnknownField\":\"UnknownField\"}")}, resource, book)
	assert.Nil(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "json title", book.Title)
	assert.Equal(t, 520, book.NbPages)
	err = engine.Deserialize(&brest.RestQuery{Action: brest.Post, Resource: "Book", ContentType: "application/json", Content: []byte("{\"NbPages\":230}")}, resource, book)
	assert.Nil(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "json title", book.Title)
	assert.Equal(t, 230, book.NbPages)

	book = &Book{}
	err = engine.Deserialize(&brest.RestQuery{Action: brest.Post, Resource: "Book", ContentType: "application/x-www-form-urlencoded", Content: []byte("Title=form title&NbPages=310&UnknownField=UnknownField")}, resource, book)
	assert.Nil(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "form title", book.Title)
	assert.Equal(t, 310, book.NbPages)
	err = engine.Deserialize(&brest.RestQuery{Action: brest.Post, Resource: "Book", ContentType: "application/x-www-form-urlencoded", Content: []byte("NbPages=450")}, resource, book)
	assert.Nil(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "form title", book.Title)
	assert.Equal(t, 450, book.NbPages)

	content, err = msgpack.Marshal(&Book{Title: "msgpack title", NbPages: 480})
	assert.Nil(t, err)
	book = &Book{}
	err = engine.Deserialize(&brest.RestQuery{Action: brest.Post, Resource: "Book", ContentType: "application/x-msgpack", Content: content}, resource, book)
	assert.Nil(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "msgpack title", book.Title)
	assert.Equal(t, 480, book.NbPages)
	content, err = msgpack.Marshal(&PageOnly{NbPages: 110})
	assert.Nil(t, err)
	err = engine.Deserialize(&brest.RestQuery{Action: brest.Post, Resource: "Book", ContentType: "application/x-msgpack", Content: content}, resource, book)
	assert.Nil(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "msgpack title", book.Title)
	assert.Equal(t, 110, book.NbPages)
}

func TestPostPatchGetDelete(t *testing.T) {
	db, config := initTests(t)
	defer db.Close()
	engine := brest.NewEngine(config)

	var err error
	var content []byte
	var res interface{}
	var page brest.Page
	var resTodo *Todo
	var resAuthor *Author
	var resAuthors []Author
	var resBook *Book

	for _, todo := range todos {
		content, err = json.Marshal(todo)
		assert.Nil(t, err)
		res, err = engine.Execute(&brest.RestQuery{Action: brest.Post, Resource: "Todo", ContentType: "application/json", Content: content})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		resTodo = res.(*Todo)
		assert.NotEqual(t, resTodo.Text, "")
	}

	for _, author := range authors {
		content, err = json.Marshal(author)
		assert.Nil(t, err)
		res, err = engine.Execute(&brest.RestQuery{Action: brest.Post, Resource: "Author", ContentType: "application/json", Content: content})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		resAuthor = res.(*Author)
		assert.NotEqual(t, resAuthor.ID, 0)
		assert.Equal(t, resAuthor.Firstname, author.Firstname)
		assert.Equal(t, resAuthor.Lastname, author.Lastname)
	}

	for _, book := range books {
		content, err = json.Marshal(book)
		assert.Nil(t, err)
		res, err = engine.Execute(&brest.RestQuery{Action: brest.Post, Resource: "Book", ContentType: "application/json", Content: content})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		resBook = res.(*Book)
		assert.NotEqual(t, resBook.ID, 0)
		assert.NotEqual(t, resBook.AuthorID, 0)
		assert.Equal(t, resBook.Title, book.Title)
		assert.Equal(t, resBook.NbPages, 0)

		res, err = engine.Execute(&brest.RestQuery{Action: brest.Patch, Resource: "Book", Key: strconv.Itoa(resBook.ID), ContentType: "application/x-www-form-urlencoded", Content: []byte("NbPages=200")})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		resBook = res.(*Book)
		assert.NotEqual(t, resBook.ID, 0)
		assert.NotEqual(t, resBook.AuthorID, 0)
		assert.Equal(t, resBook.Title, book.Title)
		assert.Equal(t, resBook.NbPages, 200)
	}

	res, err = engine.Execute(&brest.RestQuery{Action: brest.Get, Resource: "Author"})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	page = *res.(*brest.Page)
	assert.Equal(t, page.Count, 3)
	resAuthors = *page.Slice.(*[]Author)
	assert.Equal(t, len(resAuthors), 3)
	for _, author := range resAuthors {
		res, err = engine.Execute(&brest.RestQuery{Action: brest.Get, Resource: "Author", Key: strconv.Itoa(author.ID), Fields: []*brest.Field{{Name: "*"}}, Relations: []*brest.Relation{{Name: "Books"}}})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		resAuthor = res.(*Author)
		assert.Equal(t, resAuthor.ID, author.ID)
		assert.Equal(t, resAuthor.Firstname, author.Firstname)
		assert.True(t, len(resAuthor.Books) > 0)
		for _, resBook = range resAuthor.Books {
			assert.NotNil(t, resBook.Title)
			assert.Equal(t, resBook.NbPages, 200)

			res, err = engine.Execute(&brest.RestQuery{Action: brest.Put, Resource: "Book", Key: strconv.Itoa(resBook.ID), ContentType: "application/x-www-form-urlencoded", Content: []byte("Title=" + resBook.Title + "_1&AuthorID=" + strconv.Itoa(resBook.AuthorID))})
			assert.Nil(t, err)
			assert.NotNil(t, res)
			resBook2 := res.(*Book)
			assert.NotEqual(t, resBook2.ID, 0)
			assert.NotEqual(t, resBook2.AuthorID, 0)
			assert.Equal(t, resBook2.Title, resBook.Title+"_1")
			assert.Equal(t, resBook2.NbPages, 0)
		}
	}

	res, err = engine.Execute(&brest.RestQuery{Action: brest.Get, Resource: "Author", Key: "1"})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	resAuthor = res.(*Author)
	assert.Equal(t, resAuthor.Lastname, "de Saint Exupéry")

	res, err = engine.Execute(&brest.RestQuery{Action: brest.Get, Resource: "Book", Filter: &brest.Filter{Op: brest.Or, Filters: []*brest.Filter{{Op: brest.Llk, Attr: "book.title", Value: "%lo%"}, {Op: brest.Llk, Attr: "book.title", Value: "%ta%"}}}})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	page = *res.(*brest.Page)
	assert.Equal(t, 4, page.Count)

	res, err = engine.Execute(&brest.RestQuery{Action: brest.Get, Resource: "Author", Filter: &brest.Filter{Op: brest.In, Attr: "author.firstname", Value: []string{"Antoine", "Franz"}}})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	page = *res.(*brest.Page)
	assert.Equal(t, 2, page.Count)

	_, err = engine.Execute(&brest.RestQuery{Action: brest.Delete, Resource: "Author", Key: "1"})
	assert.Nil(t, err)
	res, err = engine.Execute(&brest.RestQuery{Action: brest.Get, Resource: "Author"})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	page = *res.(*brest.Page)
	assert.Equal(t, page.Count, 2)
	resAuthors = *page.Slice.(*[]Author)
	assert.Equal(t, len(resAuthors), 2)

	_, err = engine.Execute(&brest.RestQuery{Action: brest.Delete, Resource: "Author", Key: "3"})
	assert.Nil(t, err)
	res, err = engine.Execute(&brest.RestQuery{Action: brest.Get, Resource: "Author"})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	page = *res.(*brest.Page)
	assert.Equal(t, page.Count, 1)
	resAuthors = *page.Slice.(*[]Author)
	assert.Equal(t, len(resAuthors), 1)
}

func TestFormUrlencoded(t *testing.T) {
	db, config := initTests(t)
	defer db.Close()
	engine := brest.NewEngine(config)

	var err error
	var res interface{}
	var resAuthor *Author

	res, err = engine.Execute(&brest.RestQuery{Action: brest.Post, Resource: "Author", ContentType: "application/x-www-form-urlencoded", Content: []byte("Firstname=Firstname&Lastname=Lastname")})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	resAuthor = res.(*Author)
	assert.NotEqual(t, resAuthor.ID, 0)
	assert.Equal(t, resAuthor.Firstname, "Firstname")
	assert.Equal(t, resAuthor.Lastname, "Lastname")
}

func TestMsgpack(t *testing.T) {
	db, config := initTests(t)
	defer db.Close()
	engine := brest.NewEngine(config)

	var err error
	var content []byte
	var res interface{}
	var resAuthor *Author

	resAuthor = &Author{Firstname: "MsgpackFirstname", Lastname: "MsgpackLastname", Picture: []byte{187, 163, 35, 30}}
	content, err = msgpack.Marshal(resAuthor)
	assert.Nil(t, err)

	res, err = engine.Execute(&brest.RestQuery{Action: brest.Post, Resource: "Author", ContentType: "application/x-msgpack", Content: content})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	resAuthor = res.(*Author)
	assert.NotEqual(t, resAuthor.ID, 0)
	assert.Equal(t, resAuthor.Firstname, "MsgpackFirstname")
	assert.Equal(t, resAuthor.Lastname, "MsgpackLastname")
	assert.Equal(t, resAuthor.Picture, []byte{187, 163, 35, 30})
}
