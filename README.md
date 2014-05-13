Example redis app
=================

![example](http://cl.ly/image/35140k2A2K1Z/example-redis-app.gif)

Installation
------------

```
go get -u github.com/drnic/redis-counter-example
```

Configuration
-------------

Create a configuration file with connection details to a Redis server:

```json
{
  "host": "localhost",
  "port": 6379
}
```

Usage
-----

Run the server with `-config PATH` to specify the Redis connection details:

```
redis-counter-example -config path/to/redis.json
```

Development
-----------

Within the project is a default `config.json` which assumes a local Redis server (contains the example configuration above).

```
go run server.go
```
