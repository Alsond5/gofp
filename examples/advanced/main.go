package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Alsond5/gofp"
	"github.com/Alsond5/gofp/option"
	"github.com/Alsond5/gofp/result"
)

type User struct {
	ID    uint64
	Email gofp.Option[string]
}

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidEmail = errors.New("invalid email")
)

func getUser(db map[uint64]User, id uint64) gofp.Result[User] {
	return gofp.FromZero(db[id]).OkOr(ErrNotFound)
}

func validateEmail(email string) gofp.Result[string] {
	if strings.Contains(email, "@") {
		return gofp.Ok(email)
	}
	return gofp.Err[string](ErrInvalidEmail)
}

func domainAccess(email string) gofp.Option[uint8] {
	switch {
	case strings.HasSuffix(email, "@admin.com"):
		return gofp.Some[uint8](10)
	case strings.HasSuffix(email, "@user.com"):
		return gofp.Some[uint8](1)
	default:
		return gofp.None[uint8]()
	}
}

func computeAccess(db map[uint64]User, userID uint64) gofp.Result[gofp.Option[uint8]] {
	return result.FlatMap(
		getUser(db, userID),
		func(user User) gofp.Result[gofp.Option[uint8]] {
			email := user.Email
			return option.MapOr(
				email,
				gofp.Ok(gofp.None[uint8]()),
				func(email string) gofp.Result[gofp.Option[uint8]] {
					return result.Map(
						validateEmail(email),
						domainAccess,
					)
				},
			)
		},
	)
}

func main() {
	db := map[uint64]User{
		0: {ID: 0, Email: gofp.Some("")},
		1: {ID: 1, Email: gofp.Some("alice@admin.com")},
		2: {ID: 2, Email: gofp.None[string]()},
		3: {ID: 3, Email: gofp.Some("bob@unknown.com")},
	}

	for _, id := range []uint64{0, 1, 2, 3, 99} {
		computeAccess(db, id).
			IfOk(func(opt gofp.Option[uint8]) {
				opt.Match(
					func(level uint8) { fmt.Printf("id=%d → access level: %d\n", id, level) },
					func() { fmt.Printf("id=%d → Email not found or domain does not match\n", id) },
				)
			}).
			IfErr(func(err error) {
				fmt.Printf("id=%d → error: %v\n", id, err)
			})
	}
}
