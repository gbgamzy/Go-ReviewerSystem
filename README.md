# ReviewerService

A Go-based project for managing tasks with multi-approval workflows. This project provides APIs for user login/signup, task creation, approver assignment, and approval/comment functionality.

---

## Prerequisites

- **Go**: Ensure Go is installed ([Download Go](https://golang.org/dl/)).
- **PostgreSQL**: Install and configure PostgreSQL.

---

## Setting Up the Database

### 1. Create the Database

Navigate to the `psql` directory and execute the SQL script to create the database:

```bash
psql -U <your_postgres_user> -f psql/create_database.sql
```

This will create a database named `ReviewerDB`.

### 2. Create the Tables

Once the database is created, execute the script to create the required tables:

```bash
psql -U <your_postgres_user> -d ReviewerDB -f psql/create_tables.sql
```

### 3. Create Roles and Grant Permissions

Run the script to create roles and grant necessary permissions:

```bash
psql -U <your_postgres_user> -d ReviewerDB -f psql/create_role_and_permissions.sql
```

---

## Running the Application

### 1. Clone the Repository

Clone the project to your local machine:

```bash
git clone https://github.com/<your_username>/Go-ReviewerSystem.git
cd ReviewerService
```

### 2. Install Dependencies

Ensure all Go dependencies are installed:

```bash
go mod tidy
```

### 3. Run the Application

Start the server:

```bash
go run ./cmd/main.go
```

By default, the server will run on `http://localhost:8080`.

---

## Using the APIs

### Base URL

All APIs are accessible at `http://localhost:8080`.

### Endpoints with Examples and Expected Results

#### 1. **Login/Signup User**
   - **POST** `/users/login`
   - Request Body:
     ```json
     {
       "email": "john.doe@example.com",
       "name": "John Doe"
     }
     ```
   - Expected Response:
     ```json
     {
       "id": 1,
       "email": "john.doe@example.com",
       "name": "John Doe"
     }
     ```

#### 2. **Create Task**
   - **POST** `/tasks`
   - Request Body:
     ```json
     {
       "title": "Review Document",
       "description": "Review the project document.",
       "created_by": 1,
       "required_approvals": 3
     }
     ```
   - Expected Response:
     ```json
     {
       "id": 1,
       "title": "Review Document",
       "description": "Review the project document.",
       "status": "Pending",
       "created_by": 1,
       "required_approvals": 3,
       "current_approvals": 0
     }
     ```

#### 3. **Assign Approvers**
   - **POST** `/tasks/{task_id}/approvers`
   - Request Body:
     ```json
     {
       "approver_ids": [2, 3, 4]
     }
     ```
   - Expected Response:
     ```json
     {
       "message": "Approvers assigned successfully."
     }
     ```

#### 4. **Approve Task**
   - **POST** `/tasks/{task_id}/approve`
   - Request Body:
     ```json
     {
       "approver_id": 2
     }
     ```
   - Expected Response:
     ```json
     {
       "message": "Task approved successfully."
     }
     ```

#### 5. **Add Comment**
   - **POST** `/tasks/{task_id}/comments`
   - Request Body:
     ```json
     {
       "approver_id": 2,
       "comment": "Looks good!"
     }
     ```
   - Expected Response:
     ```json
     {
       "message": "Comment added successfully."
     }
     ```

#### 6. **Get Users**
   - **GET** `/users?page=1&limit=10&search=John`
   - Expected Response:
     ```json
     {
       "users": [
         {
           "id": 1,
           "name": "John Doe",
           "email": "john.doe@example.com"
         }
       ],
       "total": 1,
       "page": 1,
       "limit": 10
     }
     ```

---

## Notes

- Ensure the PostgreSQL server is running before starting the application.
- Update the connection string in `cmd/main.go` if your database credentials or host differ:

```go
connStr := "postgres://client:password@localhost:5432/ReviewerDB?sslmode=disable"
```

- Logs and errors are printed to the console for debugging.

---

## License

This project is licensed under the MIT License.
