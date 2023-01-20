FROM golang:1.19-alpine

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY src ./

RUN go build -o /goatform

EXPOSE 8080

# The collected forms are written to a local json file at ./db.json (/app/db.json)
# Goatform will listen on port *:8080
CMD [ "/goatform" ]
