# kuamua

A thin wrapper around Quamina for languages that are not Go. Quamina is a fast 
pattern matching libary written in Go. See https://github.com/timbray/quamina for more information

Kuamua is Swahili for "Decide" 

## Prerequisites

The things you need before installing the software.

* You'll need a postgres database.  This can be one dedicated for kuamua or it can be the database used by the calling application/service. 
* POSTGRESQL_URL pointing at said database.
* the SERVERPORT env variable is used to set the port number.  Defaults to 3000

```eg 'postgres://<user>:<password>@<host>:5432/kuamua?sslmode=disable&search_path=public' ```

Note: The migration code will tack on a 'x-migrations-table=kuamua_schema_migrations' to the end of that url. It's a bit gross but works for now.

### Local Development

A devenv.sh configuration is available for local development. But if not then you'll need

* Go (1.23)
* sqlc (to regenerate the queries/models)
* k6 for load testing (or your favorite)
* Bruno for API testing (or your favorite)
* And a Postgres database of course.

## Installation

Binaries are available as well as an OCI image.  Building from source is easy too.

## Usage

### Command Line

```
$ ./kuamua serve [-enable-crud] [-skip-migrations] # serve is the default command
$ ./kuamua migrate [-down]
```
The server auto-migrates the database (unless started with -skip-migrations) by default. 
You can manage the timing of migrations yourself by using the migrate subcommand.

Since it's insecure at this point, the patterns crud API is disabled by default.  Additionally if are colocating your app and kuamua together in the same database, the calling application could be used to manage the patterns.

### Patterns/API

Some information on patterns can be found at [Patterns](https://github.com/timbray/quamina/blob/main/PATTERNS.md) 

There is a bruno collection in [testing](testing/) that should point the way at using the API. And a k6 loadtest.js that should help with figuring out how to call the /document/patterns route

TODO Improve this section

## TODO

  - Improve Readme/Docs
  - Authentication for the CRUD API (jwt?)
  - OTEL support
  - Settings for cache size?

## Contributing

### Bug Reports & Feature Requests

Please use the [issue tracker](https://github.com/vhodges/kuamua/issues) to report any bugs or file feature requests.

### Pull Requests are welcome!

## License

MIT Please see [LICENSE](LICENSE) for more details.








