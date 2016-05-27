vinxi comes with an internal RESTful API for administration purposes. 
API commands can be run on any node, and vinxi will keep the configuration consistent and up-to-date in real time across the different vinxi server instances.

The RESTful Admin API server listens on port 8000 by default.

RESTful Admin API is part of [manager](https://github.com/vinxi/vinxi/tree/master/manager) package. Read more about the manager [here](#manager).

### Content Types

vinxi RESTful API is a JSON only HTTP interface. 
Responses and payloads as designed to render and read JSON only data structures.

The MIME content type used as HTTP `Content-Type` header must be:
```
application/json
```

### Authentication

You can protect the admin HTTP API with a basic authentication mechanism.

Admin HTTP API is not protected by default.

#### HTTP basic authentication

You can define multiple user/password credentials to authenticate users.

### Endpoints

#### Node information

##### Endpoint

| **GET**        |  /          |

##### Response

```
HTTP 200 OK
```

```json
{
  "hostname": "tomas-laptop",
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





