# go-api-example
Learning example demonstrating an api server and client.

## Project status: 

A basic server example is provided using net/http and gorilla/mux

To avoid the need for a database, this example uses as its data a package variable of type []User initialized with three User records.

The net/http ServeMux does not easily support path variables, so the example initially used a query parameter to specifiy the id of a user.  With the introduction of gorilla/mux, path variables are now supported.

A client is not yet provided, but the server can be tested from curl or a web browser:
1. http://localhost:8090 returns an html Hello, World! message, demonstrating the simplest way that static content can be served.
2. http://localhost:8090/api/v1/users returns a list of users in json format.
3. http://localhost:8090/api/v1/users?id=2 returns a single user with id=2 in json format.
4. http://localhost:8090/api/v1/users/2 returns a single user with id=2 in json format using the path variable.

Logging from the gorilla/handlers package has also been added to the server.

Since the project now has dependencies beyond the standard library, the root of the project has been initialized as a go module.  Both the client and server will exist in the same module.
