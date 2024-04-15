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