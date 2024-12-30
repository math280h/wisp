#!/bin/bash

# Run the Prisma CLI to deploy the Prisma schema to the database
go run github.com/steebchen/prisma-client-go migrate deploy

# Start the application
/go/bin/app
