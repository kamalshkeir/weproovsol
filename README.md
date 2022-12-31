# weproovsol

## Setup
After cloning the repos and cd into the project
```sh
go mod tidy
```

Then update in .env `DB_PASS`

### you can run
```sh
go run main.go
```

### available endpoints:
- GET /
- GET /user  (html response, list all users in a table)
- POST /user (called from /user form to create a new user)
- DELETE /user/:id (called from /user to delete a user using id)

### NOT COMPLETE
- PATCH /user/:id (due to lack of time because of table called `user` which is a reserved word)
