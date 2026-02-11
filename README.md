# bunny-database-shell

Connect to a Bunny Database shell.

## Prerequisites

- [Node.js](https://nodejs.org/)

## Install

```sh
npm install -g bunny-database-shell
```

Or run directly with `npx`:

```sh
npx bunny-database-shell
```

## Usage

Connection values are resolved in order: CLI flags, `.env` file, then interactive prompt.

### Open an interactive shell

Pass your connection details as flags:

```sh
bunny-database-shell --url libsql://your-database.lite.bunnydb.net --auth-token your-token
```

Or use a `.env` file and the shell will connect automatically:

```sh
bunny-database-shell
```

If no flags or `.env` values are found, you'll be prompted to enter them.

### Execute a single statement

Run a SQL statement and exit without entering the interactive shell:

```sh
bunny-database-shell -e "SELECT * FROM users"
bunny-database-shell "SELECT * FROM users"
```

### `.env` file

Create a `.env` file in your working directory:

```
BUNNY_DB_URL=libsql://your-database.lite.bunnydb.net
BUNNY_DB_TOKEN=your-token
```

This lets you run `bunny-database-shell` without passing flags every time.
