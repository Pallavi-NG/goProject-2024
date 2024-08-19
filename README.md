# goProject-2024
student project which performs CRUD operation in GO

API to authenticate a user - which takes 2 parameters i.e user_id and password
http://localhost:8080/api/v1/authenticate


POST API - which adds new student in a database
http://localhost:8080/api/v1/student
sample payload - 
{
  "first_name": "Pallavi",
  "last_name": "NG",
  "date_of_birth": "1998-04-23T00:00:00Z",
  "email": "pallavi.ng@gmail.com",
  "address": "shimoga",
  "gender": "Female",
  "age": 26,
  "created_by": "pallavi",
  "created_on": "2024-08-15T12:00:00Z"
}


GET API - which gets a student by id 
http://localhost:8080/api/v1/student/{id}

PUT API - which updates a student by ID
http://localhost:8080/api/v1/student/{id}
sample payload - 
{
    email : "updatedEmail@gmail.com"
}

DEL API - which deletes a student by ID 
http://localhost:8080/api/v1/student/{student}

To read data from .env file - 
we use os.GetEnv() to read data from .env file
