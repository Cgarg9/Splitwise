# Splitwise
This Readme file walks you through every step I performed to make this repo. This project is more focused on the design and experimentation

## Workflows
This contains the workflows/actions that run on pushing new code, to check for vulnerabilities, run tests etc.
Currenty this file contains 

- lint.yml : This file runs on every push to main branch and every pull request, Can be triggered manually as well(workflow_dispatch)

## DB 
This folder currently contains the file to set up db connection. This project uses postgresql 

### Migration scripts 
This folder contains the migrations scripts for db tables. Currently we have 6 tables
To run the migration script this command was used: 
```
  -path migrations \
  -database "postgres://splitwise:splitwise@localhost:5432/splitwise?sslmode=disable" \
  up
```

NOTE: migration scripts should start with 6 digit number. Usual format is <number>_<name>.up/down.sql

## Logger 
Set up logger in this project to check progress of API requests, running status of app. Will later add a local file where all logs might be saved temporarily to make it easier to debug


## Swagger documents 
Commands used 
```
# Install swag CLI tool
go install github.com/swaggo/swag/cmd/swag@latest

# Install required packages
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/files
```

After successfully installing and setting up swagger for your project add proper annotations to your route handlers and dto. Now every time you make a change to the route or add a new route run this command after updating the annotations

```
swag init -g main.go -o ./docs
```

You can visit your swagger docs here -- http://localhost:8080/swagger/index.html
It will look something like this - screenshots/Screenshot 2025-12-30 at 12.48.25 PM.png
Using this you can test APIs via postman/UI