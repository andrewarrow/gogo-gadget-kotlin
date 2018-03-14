package main

import "fmt"
import "github.com/codegangsta/cli"
import "io/ioutil"

func TableAction(c *cli.Context) {
	table := c.Args().Get(0)
	model := c.Args().Get(1)
	prefix := c.Args().Get(2)

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
	ioutil.WriteFile(fmt.Sprintf("sql/999-create-%s.sql", table), d1, 0644)

	model_template := `package %s.model
			
		import com.fasterxml.jackson.annotation.JsonIgnore
		import com.fasterxml.jackson.annotation.JsonProperty
		import java.time.OffsetDateTime
		import java.util.UUID
		import org.postgis.Point

		@TableName("%s")
		data class %s(
			val id: UUID = UUID.randomUUID(),
			val otherId: UUID? = null,
			val createdAt: OffsetDateTime = OffsetDateTime.now(),
			val updatedAt: OffsetDateTime = OffsetDateTime.now()
		)
		`

	modelfile := fmt.Sprintf(model_template, prefix, table, model)
	d1 = []byte(modelfile)
	ioutil.WriteFile(fmt.Sprintf("scb/model/%s.kt", model), d1, 0644)
}
