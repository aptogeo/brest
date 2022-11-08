package brest_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aptogeo/brest"
	"github.com/stretchr/testify/assert"
)

func decodeHandler(prefix string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		restQuery := brest.RequestDecoder(r, brest.NewConfig("/rest/", nil))
		if restQuery != nil {
			w.Write([]byte(restQuery.String()))
		}
	})
}

var requestDecoderTests = []struct {
	uri      string
	method   string
	expected *brest.RestQuery
}{
	{"/rest/User/1", "GET", &brest.RestQuery{Action: brest.Get, Resource: "User", Key: "1"}},
	{"/rest/User/445cf124-f5e6-4fd3-9f0d-d22bd6c90d40", "GET", &brest.RestQuery{Action: brest.Get, Resource: "User", Key: "445cf124-f5e6-4fd3-9f0d-d22bd6c90d40"}},
	{"/rest/User/1?fields=*,Roles", "GET", &brest.RestQuery{Action: brest.Get, Resource: "User", Key: "1", Offset: 0, Limit: 10, Fields: []*brest.Field{{Name: "*"}, {Name: "Roles"}}}},
	{"/rest/User", "GET", &brest.RestQuery{Action: brest.Get, Resource: "User", Offset: 0, Limit: 10, Fields: []*brest.Field{}, Sorts: []*brest.Sort{}, Filter: &brest.Filter{}}},
	{"/rest/User?offset=50&limit=10&sort=lastname,-firstname", "GET", &brest.RestQuery{Action: brest.Get, Resource: "User", Offset: 50, Limit: 10, Fields: []*brest.Field{}, Sorts: []*brest.Sort{{Name: "lastname", Asc: true}, {Name: "firstname", Asc: false}}, Filter: &brest.Filter{}}},
	{"/rest/User?offset=60&limit=10&sort=lastname&fields=user.*,user.roles", "GET", &brest.RestQuery{Action: brest.Get, Resource: "User", Offset: 60, Limit: 10, Fields: []*brest.Field{{Name: "user.*"}, {Name: "user.roles"}}, Sorts: []*brest.Sort{{Name: "lastname", Asc: true}}, Filter: &brest.Filter{}}},
	{"/rest/User?filter=%7B%22Op%22%3A%22ilk%22%2C%22Attr%22%3A%22title%22%2C%22Value%22%3A%22%25lo%25%22%7D", "GET", &brest.RestQuery{Action: brest.Get, Resource: "User", Offset: 0, Limit: 10, Fields: []*brest.Field{}, Sorts: []*brest.Sort{}, Filter: &brest.Filter{Op: brest.Llk, Attr: "title", Value: "%lo%"}}},
	{"/rest/User?filter=%7B%22Op%22%3A%22in%22%2C%22Attr%22%3A%22title%22%2C%22Value%22%3A%5B%22Titre+1%22%2C%22Titre+2%22%5D%7D", "GET", &brest.RestQuery{Action: brest.Get, Resource: "User", Offset: 0, Limit: 10, Fields: []*brest.Field{}, Sorts: []*brest.Sort{}, Filter: &brest.Filter{Op: brest.In, Attr: "title", Value: []string{"Titre 1", "Titre 2"}}}},
	{"/rest/User", "POST", &brest.RestQuery{Action: brest.Post, Resource: "User", ContentType: "application/json"}},
	{"/rest/User/1", "PUT", &brest.RestQuery{Action: brest.Put, Resource: "User", Key: "1", ContentType: "application/json"}},
	{"/rest/User/1", "PATCH", &brest.RestQuery{Action: brest.Patch, Resource: "User", Key: "1", ContentType: "application/json"}},
	{"/rest/User/1", "DELETE", &brest.RestQuery{Action: brest.Delete, Resource: "User", Key: "1"}},
	{"/rest/User/specific/otherservice", "GET", nil},
	{"/rest", "GET", nil},
	{"/", "GET", nil},
}

func TestRequestDecoder(t *testing.T) {
	ts := httptest.NewServer(decodeHandler("/rest/"))
	defer ts.Close()

	for _, rt := range requestDecoderTests {
		req, err := http.NewRequest(rt.method, ts.URL+rt.uri, bytes.NewBufferString(""))
		assert.Nil(t, err)
		res, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		body, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		err = res.Body.Close()
		assert.Nil(t, err)
		if rt.expected != nil {
			assert.Equal(t, string(body), rt.expected.String())
		} else {
			assert.Equal(t, string(body), "")
		}
	}
}
