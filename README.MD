# sql_idempotency_doctor

A CLI tool to check the idempotency of SQL scripts in PostgreSQL.

sql_idempotency_doctor is a user-friendly command-line interface that verifies if the PLPGSQL scripts in the specified `deploy` and `revert` directories (usually created & maintained by [`sqitch`](https://github.com/sqitchers/sqitch)) are idempotent. 

This tool is built using the [Cobra](https://github.com/spf13/cobra) framework and focuses on PostgreSQL databases.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Development](#development)

## Installation

1. Download the latest binary for your platform from the [releases](https://github.com/your_username/sql_idempotency_doctor/releases) page.
2. Extract the archive and place the `sql_idempotency_doctor` executable in a directory included in your system's `PATH`.

Alternatively, you can build the project from source:

1. Ensure you have [Go](https://golang.org/doc/install) installed on your system.
2. Clone the repository and navigate to the project root folder:

```sh
git clone https://github.com/your_username/sql_idempotency_doctor.git
cd sql_idempotency_doctor
```

3. Build the executable:
```sh
go build -o sql_idempotency_doctor
```
4. Move the sql_idempotency_doctor executable to a directory included in your system's PATH.

## Usage
To check the idempotency of SQL scripts, run the following command:

```sh
sql_idempotency_doctor check -p /path/to/directory -n namespace
```
Replace /path/to/directory with the actual path to the directory containing the deploy and revert folders, and replace namespace with the namespace you want to check.

## Examples
Consider the following directory structure:

```
/path/to/directory
├── deploy
│   └── namespace.sql
└── revert
    └── namespace.sql
```

To check the idempotency of the SQL scripts in the deploy and revert directories, run:

```sh
sql_idempotency_doctor check -p /path/to/directory -n namespace
```
where scripts in `<namespace.sql>` that require assessment must be decorated with:
```sql
--@ddl:start
<SQL here>
--@ddl:end
```

## Development
To set up the development environment, follow these steps:

1. Clone the repository:
```
git clone https://github.com/your_username/sql_idempotency_doctor.git
```
2. Navigate to the project root folder:
```
cd sql_idempotency_doctor
```
3. Install dependencies:
```
go mod download
```
4. Make your changes and build the executable:
```
go build -o sql_idempotency_doctor
```
5. Test your changes:
```
./sql_idempotency_doctor check -p /path/to/directory -n namespace
```

