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