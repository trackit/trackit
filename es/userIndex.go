package es

import (
	"fmt"

	"gopkg.in/olivere/elastic.v5"

	"github.com/trackit/trackit2/users"
)

const (
	IndexPrefixLineItems = "lineitems"
)

var Client *elastic.Client

func init() {
	var err error
	Client, err = elastic.NewClient(
		elastic.SetBasicAuth("elastic", "changeme"),
	)
	if err != nil {
		panic(err)
	}
}

func IndexNameForUser(u users.User, p string) string {
	return IndexNameForUserId(u.Id, p)
}

func IndexNameForUserId(i int, p string) string {
	return fmt.Sprintf("%06d-%s", i, p)
}
