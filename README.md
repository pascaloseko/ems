# Employee Management System

## Detailed Requirement
1. Setup Docker Services: Provide a Docker Compose file to set up an Azure SQL Edge database sever.
2. Create an Azure SQL database: On the server from step 1, create a database with a table named 'Employee' and another called 'Department'. The Employee table should have the following fields: ID (Primary Key, Auto Increment, Not Null), First Name, Last Name, Username, Password, Email, DOB (Date of Birth), DepartmentID and Position.
3. Create an Employee GraphQL schema: The schema should represent an Employee with all the fields mentioned above and necessary resolving to the database in step 2
4.  API Endpoints:

    ◦ /login: This endpoint is for employee authentication. The endpoint should accept POST requests with 'username' and 'password' and, if authenticated successfully, return a JSON Web Token (JWT).

    ◦ /employee: This endpoint is for querying and mutating the 'Employee' table. It should only accept requests with a valid JWT in the Authorisation header. The requests should be in the GraphQL format, and the endpoint should perform the requested query/mutation on the 'Employee' table.

5. GraphQL Operations:

    ◦ Query: Get a list of all employees (filtering, sorting, pagination), get a specific employee by ID, get the current employee by context.

    ◦ Mutation: Add a new employee, update an existing employee, delete an employee

6. JWT Middleware: Implement a middleware to check the JWT from the request header. The request should only be accepted if the JWT is provided or is valid. If the JWT is valid, the request should be allowed to proceed.

## Additional Requirements Checklist

[ ] The source code pushed to a Git repository.

[ ] A README.md file with instructions on how to build and run the code, and any other relevant documentation.

[ ] Postman collection, Playground or similar for testing the endpoints.

[ ] Docker Compose file for setting up necessary services.

[ ] Any scripts required to set up the database.

## How to run the APP
 # Step 1:
- run the app

    ```
    docker-compose up --build -d
    ```

# Step 2
- test it using curl the GET employees endpoint
    ```
    curl --location --request GET 'http://localhost:8080/employees' \
    --header 'Authorization: "test"' \
    --header 'Content-Type: application/json' \
    --data '{
        "username": "test",
        "password": "test"
    }'
    ```

- POST login endpoint
    ```
    curl --location 'http://localhost:8080/login' \
    --header 'Content-Type: application/json' \
    --data '{
        "username": "test",
        "password": "test"
    }'
    ```

# Database Layer
- The database is bootstrapped from internal/pkg/db/database/mssql.go by InitDB function
- there is a connection string and I used gorm to make migrations automatically


# API Layer
- So my approach was to first initialize the graphql this will generate the relevant folders(see graph folder)
- then I created the docker compose containing the azure sql server and was able to run it(apply it on step one)
- I proceeded changing the data in schema.resolvers.go contents to match with the requirements
- the above helped me now design the API which included a middleware that checks a logged in user is authenticated.
- the login endpoint is located in the internal/handlers/handlers.go, it gets the data coming from the client and pass it down to the mutation resolver to create a user. if there are any errors they will be returned with the relevant status code and message.
- the employees handlers is also situated in the above package where it returns a list of employees from the database. The endpoint is protected in the server.go file line 43.
- if non authorized a status code of 401/403 will be thrown from the middleware in internal/auth/middleware.go

# Problems
- There is an underlying issue when testing the app using postman/curl [Update this is resolved!]

```
curl --location 'http://localhost:8080/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "test",
    "password": "test"
}'
```

the above throws a 500 with the below log
```
2024/04/16 15:23:06 ERROR failed to save employee: sql: database is closed
``` 

not sure why the database keeps getting closed

## NOTE

### Known Issues
- There are still alot to be done like testing the API and data layers to make sure they are working as expected.


- Also the issue of the datbase getting closed when doing API operations needs to be resolved. [Update this is resolved]

- Also I couldn't get the app running on docker, due to a connection problem with the mssql database [Update this is resolved]


Those are the 3 main issues that I faced but would have been able to resolve if I had more time on this and also another set of eyes to help me debug.

## nice to have
gostatic checker
ci/cd pipeline
