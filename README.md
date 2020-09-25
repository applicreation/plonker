# Plonker

> While in version `0.*` any of the config can change drastically between minor releases.

## Install

### OSX

```shell script
$ brew tap applicreation/homebrew-taps
$ brew install plonker
```

# Usage

## Connection

The `conection.yaml` is a reserved file name and is used for storing the source database connection details.

```yaml
# connection.yaml
engine: mysql
host: localhost
port: 3306
name: main
username: root
password: password
```

## Config

Each YAML file represents a table in the database.

All tables that are referenced need to be declared, including relation tables.  
However, relation tables do not need to specify a range if not needed.

Range can be declared with either `percentage` or `records` with an integer value.

```yaml
# items.yaml
name: items
primaryKey: id
range:
  percentage: 10
relations:
  - key: user_id
    table: user
    column: id
```

```yaml
# user.yaml
name: user
primaryKey: id
```

## Execution

Plonker currently only supports `dry-run` as a command and will scan for `*.yaml` files in the directory it is executed in.

```shell script
plonker dry-run
```
