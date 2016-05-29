vinxi comes with an internal RESTful API for administration purposes. 
API commands can be run on any node, and vinxi will keep the configuration consistent and up-to-date in real time across the different vinxi server instances.

The RESTful Admin API server listens on port 8000 by default.

RESTful Admin API is part of [manager](https://github.com/vinxi/vinxi/tree/master/manager) package. Read more about the manager [here](#manager).

## Content Types

vinxi RESTful API is a JSON only HTTP interface. 
Responses and payloads as designed to render and read JSON only data structures.

The MIME content type used as HTTP `Content-Type` header must be:
```
application/json
```

## Authentication

You can protect the admin HTTP API with a basic authentication mechanism.

Admin HTTP API is not protected by default.

### HTTP basic authentication

You can define multiple user/password credentials to authenticate users.

## Endpoints

### Node information

```
GET /
```

#### Response

```
HTTP 200 OK
```

```json
{
  "hostname": "vinxi-server",
  "version": "0.1.0",
  "runtime": "go1.6",
  "platform": "darwin",
  "cpus": 8,
  "gorutines": 13,
  "links": {
    "catalog": "/catalog",
    "instances": "/instances",
    "manager": "/manager",
    "plugins": "/plugins",
    "scopes": "/scopes"
  }
}
```

### Global plugins

#### List plugins

```
GET /plugins
```

##### Response

```json
[
  {
    "id": "UGEZGfua5S3U0nDS",
    "name": "auth",
    "description": "Authorization and authentication protection",
    "config": {
      "scheme": "Bearer",
      "token": "s3cr3t"
    }
  }
]
```

#### List plugins

```
GET /plugins
```

##### Response

```json
[
  {
    "id": "UGEZGfua5S3U0nDS",
    "name": "auth",
    "description": "Authorization and authentication protection",
    "config": {
      "scheme": "Bearer",
      "token": "s3cr3t"
    }
  }
]
```

#### Register plugin

```
GET /plugins/{id}
```

##### Request Body

```json
{
  "name":"auth", 
  "config": {
    "token": "s3cr3t"
  }
}
```

##### Response

```json
{
  "id": "UGEZGfua5S3U0nDS",
  "name": "auth",
  "description": "Authorization and authentication protection",
  "config": {
    "scheme": "Bearer",
    "token": "s3cr3t"
  }
}
```

#### Get plugin

```
GET /plugins/{id}
```

##### Response

```json
{
  "id": "UGEZGfua5S3U0nDS",
  "name": "auth",
  "description": "Authorization and authentication protection",
  "config": {
    "scheme": "Bearer",
    "token": "s3cr3t"
  }
}
```

#### Delete plugin

```
DELETE /plugins/{id}
```

#### Response

```
HTTP 204 No Content
```

### Instances

#### List instances

```
GET /instances
```

#### Response

```
HTTP 200 OK
```

```json
[
  {
    "info": {
      "id": "KamlJHrvH2owdzIG",
      "name": "default",
      "description": "This a default proxy",
      "hostname": "vinxi-server",
      "platform": "darwin",
      "server": {
        "port": 3100,
        "readTimeout": 0,
        "writeTimeout": 0,
        "address": ""
      }
    },
    "scopes": [
      {
        "id": "jd7CDZ48IsQ0RJqu",
        "name": "custom",
        "rules": [
          {
            "id": "MizprxflGgjeQe06",
            "name": "path",
            "description": "Matches HTTP request URL path againts a given path pattern",
            "config": {
              "path": "\/foo\/(.*)"
            }
          }
        ],
        "plugins": [
          {
            "id": "AsgTHxwUgUcg9oet",
            "name": "forward",
            "description": "Forward HTTP traffic to remote servers",
            "config": {
              "url": "http:\/\/httpbin.org"
            }
          }
        ]
      }
    ]
  }
]
```

#### Get instance

```
GET /instances/{id}
```

#### Response

```
HTTP 200 OK
```

```json
{
  "info": {
    "id": "KamlJHrvH2owdzIG",
    "name": "default",
    "description": "This a default proxy",
    "hostname": "test",
    "platform": "darwin",
    "server": {
      "port": 3100,
      "readTimeout": 0,
      "writeTimeout": 0,
      "address": ""
    }
  },
  "scopes": [
    {
      "id": "jd7CDZ48IsQ0RJqu",
      "name": "custom",
      "rules": [
        {
          "id": "MizprxflGgjeQe06",
          "name": "path",
          "description": "Matches HTTP request URL path againts a given path pattern",
          "config": {
            "path": "\/foo\/(.*)"
          }
        }
      ],
      "plugins": [
        {
          "id": "AsgTHxwUgUcg9oet",
          "name": "forward",
          "description": "Forward HTTP traffic to remote servers",
          "config": {
            "url": "http:\/\/httpbin.org"
          }
        }
      ]
    }
  ]
}
```

#### Delete instance

```
DELETE /instances/{id}
```

#### Response

```
HTTP 204 No Content
```
