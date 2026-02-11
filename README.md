# bunny-database-shell

Connect to a Bunny Database shell.

## Prerequisites

- [Go](https://go.dev/dl/)
- [Node.js](https://nodejs.org/)

## Usage

```sh
npx bunny-database-shell --url <URL> --auth-token <TOKEN>
```

### Configuration

Connection values are resolved in order:

1. CLI flags (`--url`, `--auth-token`)
2. `.env` file (`BUNNY_DB_URL`, `BUNNY_DB_TOKEN`)
3. Interactive prompt

### `.env` example

```
BUNNY_DB_URL=libsql://your-database.bunny.net
BUNNY_DB_TOKEN=your-token
```

## Install globally

```sh
npm install -g bunny-database-shell
```

Then run:

```sh
bunny-database-shell
```
