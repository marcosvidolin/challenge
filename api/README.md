# API

The API is responsible for exposing user information.

Its is not a Swagger documentation, but I hope it serves as a helpful resource =)

Health Check:

You can verify the API’s availability with the following command:

```shell
curl http://localhost:8080/api/health
```

To get an user by ID:

To retrieve a specific user by their ID, use the following command:

```shell
curl http://localhost:8080/api/users/26
```

To search for users:

ou can search for users based on different criteria. Here's an example of how to perform a search:

```shell
http://localhost:8080/api/users?last_name=Ana&first_name=nat&fields=id,first_name,email_address,last_name,last_name
```

You can select the fields you want to see in the response using the `fields` param
