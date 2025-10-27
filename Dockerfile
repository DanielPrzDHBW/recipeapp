# syntax=docker/dockerfile:1

FROM golang:1.25.1-alpine

# Set destination for COPY
WORKDIR /recipeapp

# Download Go modules
COPY recipeapp/go.mod recipeapp/go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY recipeapp/*.go ./
COPY recipeapp/api/ ./api/
COPY recipeapp/database/ ./database/
COPY recipeapp/client/ ./client/
COPY recipeapp/models/ ./models/
COPY recipeapp/serverError/ ./serverError/
COPY recipeapp/recipes.db ./recipes.db/
COPY recipeapp/cookie/ ./cookie/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /recipeapp

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose

EXPOSE 8080

ENTRYPOINT ["/recipeapp"]
