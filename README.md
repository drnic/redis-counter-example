Example redis app
=================

![example](http://cl.ly/image/35140k2A2K1Z/example-redis-app.gif)

Installation
------------

```
go get -u github.com/drnic/redis-counter-example
```

Usage
-----

### With configuration file

Create a configuration file with connection details to a Redis server:

```json
{
  "host": "localhost",
  "port": 6379
}
```

Run the server with `-config PATH` to specify the Redis connection details:

```
redis-counter-example -config path/to/redis.json
```

### With environment variables

Set the environment variables:

```
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=password
```

Run the server and Redis connection details loaded from `$HOST` and `$PORT`:

```
redis-counter-example
```

Development
-----------

Within the project there is `config.json` which assumes a local Redis server (contains the example configuration above).

```
go run server.go -config config.json
```
