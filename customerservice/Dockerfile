# Build Stage
# First pull Golang image
FROM golang:1.25.3-alpine as build-env

# Set envirment variable
ENV APP_NAME customerservice
ENV CMD_PATH main.go

# Maintainer Info
LABEL maintainer="Rajkumar Palani"

# Copy application data into image
COPY . $GOPATH/src/$APP_NAME
COPY ./configuration/dbconfig.yaml $GOPATH/src/configuration/dbconfig.yaml
WORKDIR $GOPATH/src/$APP_NAME

# Build application
#RUN CGO_ENABLED=0 go build -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH
RUN go build -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

# Run Stage
FROM alpine

# Set envirment variable
ENV APP_NAME customerservice

# Copy only required data into this image
COPY --from=build-env /$APP_NAME .
COPY ./configuration/dbconfig.yaml ./configuration/dbconfig.yaml

# Expose application port
EXPOSE 3000

# Start app
CMD ./$APP_NAME
