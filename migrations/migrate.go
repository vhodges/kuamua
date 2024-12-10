package migrations

import (
    "fmt"
    "embed"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    "github.com/golang-migrate/migrate/v4/source/iofs"

)

//go:embed *.sql
var fs embed.FS

func Migrate(down bool) error {

    log.Println("Migrating database...")

    d, err := iofs.New(fs, ".")
    if err != nil {
        return err
    }

    // This is a bit of a hack.  TODO Do this better.
    purl := fmt.Sprintf("%s&x-migrations-table=kuamua_schema_migrations", os.Getenv("POSTGRESQL_URL"))
    m, err := migrate.NewWithSourceInstance("iofs", d, purl)
    if err != nil {
        return err
    }

    if down {
        err = m.Down()
    } else {
        err = m.Up()
        if err != nil && err == migrate.ErrNoChange {
            log.Printf("No database changes\n")
            return nil
        }

    }

    if err != nil {
        return err
    }

    log.Println("Done")

    return nil
}
