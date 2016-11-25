# Stormrage

Backend code for Abemar's client management system.

## Configuration

Before running the project, it is necessary to have a Go workspace and the
$GOPATH environment variable. Read [How to Write Go Code](https://golang.org/doc/code.html)
to configure the project correctly.

### Downloading the project

```shell
go get -u github.com/ab22/stormrage
```

### Defining environment variables

A few environment variables must be set in order to run the application correctly.
The required env variables are:

- SECRET
- PORT - 1337 by default.
- ENV - "DEV" by default.
- PR_ADDR
- PR_PORT
- PR_USER
- PR_PASS

These variables can be copied from the heroku config variables.

### Database Migrations

It is required to have installed Postgres on the local computer. All migration
files are saved in the migrations folder. To automatically run these queries,
it is recommended to use the [migrate](https://github.com/mattes/migrate) tool.

```shell
go get -u github.com/mattes/migrate
```

Console syntax to migrate all queries:

```shell
cd github.com/ab22/abcd/
migrate -url postgres://user:pass@host:port/dbname?sslmode=disable -path ./migrations up
sudo -u postgres psql -d abemar -a -f ~/go/src/github.com/ab22/stormrage/migrations/seed.sql
```

Note: It is required to have **MinGW32/64bit** installed on **Windows**!

## Running the application

### Compiling and Running

To compile and run the project:

```shell
go build -o stormrage.o && ./stormrage.o
```
