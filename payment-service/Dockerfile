# Build Stage
# First pull Golang image
FROM golang:1.25.3-alpine as build-env

# Set envirment variable
ENV APP_NAME paymentservice
ENV CMD_PATH main.go


# Copy application data into image
COPY . $GOPATH/src/$APP_NAME
COPY ./config/dbconfig.yaml $GOPATH/src/config/dbconfig.yaml
WORKDIR $GOPATH/src/$APP_NAME

# Build application
#RUN CGO_ENABLED=0 go build -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH
RUN go build -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

# Run Stage
FROM alpine

# Set envirment variable
ENV APP_NAME paymentservice

# Copy only required data into this image
COPY --from=build-env /$APP_NAME .
COPY ./config/dbconfig.yaml ./config/dbconfig.yaml

# Expose application port
EXPOSE 3000

# Start app
CMD ./$APP_NAME