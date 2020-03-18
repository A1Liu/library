package main

import (
	"errors"
	"flag"
	"github.com/A1Liu/webserver/users"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"log"
	"net/url"
)

type URL url.URL;

func (p *URL) Set(in string) error {
	u, err := url.Parse(in)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "psql", "postgresql":
	default:
		return errors.New("unexpected scheme in URL")
	}

	*p = URL(*u)
	return nil
}

func (p URL) String() string {
	return (*url.URL)(&p).String()
}


func main() {
	var u URL
	flag.Var(&u, "postgres-url", "URL formatted address of the postgres server to connect to")
	flag.Parse()

	if u.String() == "" {
		log.Fatal("Flag postgres-url is required")
	}

	c, err := pgx.ParseURI(u.String())

	db := stdlib.OpenDB(c)

	if err:=users.ValidateSchema(db); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
}
