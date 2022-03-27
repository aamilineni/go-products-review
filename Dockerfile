FROM golang as builder

RUN mkdir /app  
ADD . /app

WORKDIR /app
RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 go build -o go-products-review

# Set the environment variables for basic auth token
ENV AUTH_USERNAME=IDT
ENV AUTH_PASSWORD=password

CMD ["./go-products-review"]