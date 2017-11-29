package es

import (
	"fmt"

	"github.com/trackit/trackit2/users"
)

const (
	IndexPrefixLineItems = "lineitems"
)

func IndexNameForUser(u users.User, p string) string {
	return IndexNameForUserId(u.Id, p)
}

func IndexNameForUserId(i int, p string) string {
	return fmt.Sprintf("%06d-%s", i, p)
}
