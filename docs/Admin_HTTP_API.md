vinxi comes with an internal RESTful API for administration purposes. 
API commands can be run on any node, and vinxi will keep the configuration consistent and up-to-date in real time across the different vinxi server instances.

The RESTful Admin API server listens on port 8000 by default.

RESTful Admin API is part of [manager]() package. Read more about the manager [here](#manager).

### Content Types

vinxi RESTful API is a JSON only HTTP interface. 
Responses and payloads as designed to render and read JSON only data structures.

The MIME content type used as HTTP `Content-Type` header must be:
```
application/json
```

### Security

RESTful API admin interface can be easily configured
