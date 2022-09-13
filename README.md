# Bank API

## REST Endpoints

### Create account
URI: POST http://localhost:8080/v1/account/

Body request example:

     {
		"name": "bill smith",
		"amount": 10.0 
	 }

### Add money
URI: PATCH http://localhost:8080/v1/account/[accountID]/money

Example: PATCH http://localhost:8080/v1/account/5b7a411e-051c-4010-b9f1-f102c09768a0/money

Body request example:

     {
		"amount": 10.0 
	 }


### Transfer money
URI: POST http://localhost:8080/v1/transfer/

Body request example:

    {
		"from":"9a937bcf-7351-4f2b-8087-ab7dc076621c",
		"to":"e7569452-8f05-4d59-891c-36b7a5156f16".
		"amount": 10.0 
	}

Where ***from*** is account source and ***to*** is account destiny.

### Get account
URI: GET http://localhost:8080/v1/account/[accountID]/

Example: http://localhost:8080/v1/account/5b7a411e-051c-4010-b9f1-f102c09768a0

### Get all accounts
URI: GET http://localhost:8080/v1/account/

## Running the application

Execute in a terminal

    docker-compose up

It will load api (using 8080 port) and mysql docker images.

## Running tests

`docker-compose up` should be running in order to execute `main_test.go`

## Caveats

Concurrency was managed by used *mutex*. This works for as matter of this quick example but not scalable at all. Better using channels but implementation might take longer as being addressed as event driven app.
