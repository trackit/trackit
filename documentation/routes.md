# Routes

You can easily add routes to the REST API.
A route is defined it's HTTP method, the function it calls, the decorators it will use and it's path.

You can register a route in an `init` function.
You have to import the package somewhere for this function to be called. For example: 
```go
import _ "github.com/trackit/trackit/aws/routes"
```

The API exposes it's own documentation. You have to use the `Documentation` decorator in order to document the routes you create.

### Example:
This route calls the function `routeAwsNext` and use the `RequestTransaction`, `RequireAuthenticatedUser` and `Documentation` decorators. It's a GET request and it's path is `/aws/next`.
```go
func init() {
    routes.MethodMuxer{
        http.MethodGet: routes.H(routeAwsNext).With(
            db.RequestTransaction{db.Db},
            users.RequireAuthenticatedUser{users.ViewerCannot},
            routes.Documentation{
                Summary:     "get data to add next aws account",
                Description: "Gets data the user must have in order to successfully set up their account with the product.",
            },
        ),
    }.H().Register("/aws/next")
}
```