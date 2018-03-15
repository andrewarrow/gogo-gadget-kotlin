package main

import "fmt"
import "github.com/codegangsta/cli"
import "strings"
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

	repo_template := `package %s.repository

import %s.core.DatabaseManager
import %s.model.%s
import org.postgis.Geometry
import org.postgis.Point
import java.time.OffsetDateTime
import java.util.UUID
import javax.inject.Inject

interface %sRepository {
  fun insert(thing: %s): %s
	fun byUser(userId: UUID): List<%s>
}

class %sRepositoryImpl @Inject constructor(
  private val db: DatabaseManager
) : %sRepository {

  override fun insert(thing: %s): %s {
    %s::class.insert()
      .values(thing)
      .generate(db.dialect)
      .execute(db.primary)
    return thing
  }
	overrride fun byUser(userId: UUID): List<%s> {
	  return %s::class.selectAll()
      .where(%s::userId eq userId)
      .generate(db.dialect)
      .query(db.primary)
	}
}
	`

	repofile := fmt.Sprintf(repo_template,
		prefix, prefix, prefix,
		model, model, model, model, model,
		model, model, model, model, model, model, model, model)

	d1 = []byte(repofile)
	ioutil.WriteFile(fmt.Sprintf("scb/repository/%sRepository.kt", model), d1, 0644)

	man_template := `package %s.core

import %s.model.%s
import %s.repository.%sRepository
import org.postgis.Point
import java.time.OffsetDateTime
import java.util.UUID

interface %sManager {
  fun insert(thing: %s): %s
	fun byUser(userId: UUID): List<%s>
}

class %sManagerImpl(
  private val %sRepository: %sRepository
): %sManager {

  override fun insert(thing: %s): %s {
    %sRepository.insert(thing)
	}
	overrride fun byUser(userId: UUID): List<%s> {
		%sRepository.get(userId)
	}
}

	`
	n := strings.Split(table, "_")[0]
	manfile := fmt.Sprintf(man_template, prefix, prefix, model,
		prefix, model, model, model,
		model, model, n, model, model,
		model, model, n, model, model, n)
	d1 = []byte(manfile)
	ioutil.WriteFile(fmt.Sprintf("scb/core/%sManager.kt", model), d1, 0644)

	res_template := `package %s.resource

import %s.core.%sManager
import %s.model.%s
import com.fasterxml.jackson.annotation.JsonProperty
import com.newrelic.api.agent.Trace
import io.dropwizard.auth.Auth
import javax.inject.Inject
import javax.validation.Valid
import javax.validation.constraints.NotNull
import javax.ws.rs.POST
import javax.ws.rs.GET
import javax.ws.rs.PUT
import javax.ws.rs.Path
import javax.ws.rs.Produces
import javax.ws.rs.core.MediaType

@Path("/%s")
@Produces(MediaType.APPLICATION_JSON)
class %sResource @Inject constructor(
  private val %sManager: %sManager
): Logging {

  @Trace(dispatcher = true)
  @POST
  fun create%s(
    @Valid @NotNull body: Create%sBody
		 ): %s = %sManager.create%s(body.to%s())
	 }

  @Trace(dispatcher = true)
  @GET
  fun get%s(
    @Auth principal: UserPrincipal,
  ): List<%s> = %sManager.get%s(principal.user_id)
	 
data class Create%sBody(
  val otherId: UUID,
  val thing: String? = null
) {
  fun to%s(): %s = %s(
    otherId = otherId
  )
}
		`
	resfile := fmt.Sprintf(res_template, prefix, prefix, model, prefix, model, table, model, n, model, model, model, model, n, model, model, model, model, model, model)
	d1 = []byte(resfile)
	ioutil.WriteFile(fmt.Sprintf("scb/resource/%sResource.kt", model), d1, 0644)

}
