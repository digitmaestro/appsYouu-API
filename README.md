# Keskin API

This repository contains the Keskin GraphQL API.

# Local Setup

### 1. Install Prisma Version 1 Cli

```bash
npm install -g prisma1
```

### 2. Install Go Version 1.13

```bash
go get golang.org/dl/go1.13
go1.13 download
go1.13 get
go1.13 generate
```

OR

You can install Go Version Manager
```bash
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)

gvm install go1.13
gvm use go1.13 --default
```

### 3. Set up database & deploy Prisma datamodel

Start the prisma server and the database using docker-compose:

```bash
docker-compose up -d
```

### 4. Generate Code With Wire and Prisma 
After you made changes to prisma or graphql schema files, just generate the files:

```bash
go generate # This actually runs below commands along with gqlgen you can look at main.go for other commands
prisma1 generate # This will generate Prisma Client
prisma1 deploy # This will migrate database to postgres using Prisma Server
```

If for  `prisam1 deploy` fails then you can change the endpoint with env `PRISMA_HOST` and `APP_NAMESPACE` to the hard coded ones in the `prisma > prisma.yml` file.

For Seeding you can run 
```bash
prisma1 seed # This will run the main function in `seed > seed.go as configured in prisma.yml`
```

### 5. Start the GraphQL server

```bash
go run ./server
```

Navigate to [http://localhost:4000](http://localhost:4000) in your browser to explore the API of your GraphQL server in a [GraphQL Playground](https://github.com/prisma/graphql-playground).

