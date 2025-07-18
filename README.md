# dbscript

Dbscript is a basic client for MySQL and PostgreSQL that provides JavaScript scripting capabilities for CDC events.

## Installation

```shell
go install github.com/JayJamieson/dbscript
```

## Usage

First create a basic event handling script that adds addition metadata to an event

```js
function handler() {
  // Get an event for processing
  var event = dbscript.ctx.getEvent()

  // Event has metadata key and payload key
  // {
  //   "metadata": {
  //     "id": "73a777c3-6177-478d-9075-eabba0b96922"
  //
  //     // seconds since unix epoch
  //     "timestamp": 1752569637,
  //   }
  //   "payload": {
  //   }
  // }

  // Add a new key to metadata
  event.metadata["my-key"] = "my value";

  // if there is an error case you can use the built in dbscript.error(error, event)
  // when an event is errored, it will be re-attempted by default 3 times or a configurable amount
  // each attempt will add a "retryCount" counter to event meta data
  //
  // if there is an uncaugh exception, the behavior is the same as explicitly calling dbscript.error(error, event)
  if (event.payload.someValue === "something") {
    dbscript.ctx.error(new Error("Some error handling event"), event);
  }

  // events can also be dropped, internally this sets "retryCount" to max value
  // a reason for dropping the event can also be set for logging and debuggin purposes
  if (event.payload.otherValue === "something else") {
    dbscript.ctx.drop("reason", event);
  }

  // once processing is complete, call dbscript.ok(event)
  // this moves the event forward for more processing or to it's final destination sink
  dbscript.ctx.ok(event)
}
```

For the most basic setup without a pipeline or advanced configuration can be started as follows:

```shell
dbscript start -u dbscript -H localhost -p 3306 --password password --schema dbscript --tables events --handler myhandler.js
```

This will connect to your MySQL databases and listen for all change events on table `events` and forward them to your handler. The default sink for this is stdout, to use other sinks see **Configuration** for more details.

## Development Setup

### Start MySQL Database

Start a MySQL 8.0 database with Docker Compose:

```bash
docker-compose up -d
```

This will create a MySQL database with:

- Database: `dbscript`
- User: `dbscript` / Password: `dbscript`
- Root password: `rootpassword`
- Pre-created tables: `events` and `user`

### Insert Dummy Data

To insert dummy data into the database for testing:

```bash
docker exec -i dbscript_mysql mysql -u dbscript -pdbscript dbscript < docker/mysql/scripts/insert-dummy-data.sql
```

### Connect to Database

To connect to the MySQL database directly:

```bash
docker exec -it dbscript_mysql mysql -u dbscript -pdbscript dbscript
```

### Stop Database

To stop the database:

```bash
docker-compose down
```

### Run from source

```sh
go run main.go start -u dbscript -H localhost -p 3306 --password password --schema dbscript --tables events,user --handler myhandler.js
```
