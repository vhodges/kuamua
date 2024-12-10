package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/vhodges/kuamua/migrations"
    "github.com/vhodges/kuamua/server"
)
	
func main() {

    serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
    enableCrud := serveCmd.Bool("enable-crud", false, "If enabled, allow access to the Patterns CRUD API")
    skipMigrations := serveCmd.Bool("skip-migrations", false, "Don't run migrations")

    migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
    migrateDown := migrateCmd.Bool("down", false, "Set to migrate down instead of up")

    // Default to running the server
    command := "serve"

    if len(os.Args) > 1 {
        command = os.Args[1]
    }

    switch command {

    case "serve":        
        if len(os.Args) > 1 {
            serveCmd.Parse(os.Args[2:])
        } else {
            serveCmd.Parse(os.Args)
        }

        // Auto migrate unless user says not too
        if !*skipMigrations {
            err := migrations.Migrate(false) // Up, not down
            if err != nil {
                log.Fatal(err)
            }    
        }

        service := server.New()
        service.Run(*enableCrud)

    case "migrate":
        if len(os.Args) > 1 {
            migrateCmd.Parse(os.Args[2:])
        } else {
            migrateCmd.Parse(os.Args)
        }
        
        // If not auto migrating, run this command to migrate the db
        err := migrations.Migrate(*migrateDown)
        if err != nil {
            log.Fatal(err)
        }
    default:
        fmt.Println("expected 'serve' or 'migrate' subcommands")
        os.Exit(1)
    }
}


