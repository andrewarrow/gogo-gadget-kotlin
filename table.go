package main

import "fmt"
import "github.com/codegangsta/cli"
import "io/ioutil"

func TableAction(c *cli.Context) {
	table := c.Args().Get(0)

	sql := `CREATE TABLE %s (
						id            uuid NOT NULL DEFAULT uuid_generate_v4(),
						other_id       uuid,
						created_at    timestamp with time zone NOT NULL DEFAULT now(),
						updated_at    timestamp with time zone NOT NULL DEFAULT now()
					);

					ALTER TABLE ONLY %s
						ADD CONSTRAINT %s_pkey PRIMARY KEY (id);
					CREATE INDEX index_%s_by_other ON %s(other_id);
					`

	sqlfile := fmt.Sprintf(sql, table, table, table, table, table)
	d1 := []byte(sqlfile)
	ioutil.WriteFile("sql/999.sql", d1, 0644)
}