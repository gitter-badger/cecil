
## Generating goa code

```
./goagen.sh
```

## Swagger Endpoint

To view the Swagger spec in JSON format, go to:

```
curl http://host:port/swagger.json
```

Replacing `host:port` with the host and port where you are running cecil

## Regenerate gobindata email templates

```bash
$ cd core/email-templates
$ $GOPATH/bin/go-bindata -o=../email-templates.go -pkg="emailtemplates" .
$ mv ../email-templates.go ../../emailtemplates/templates.go
```

## Listing of code directories/files and their purposes

See the [Code Inventory](docs/CodeInventory.md)
