FROM golang as builder

RUN mkdir /app  
ADD . /app

WORKDIR /app
RUN go mod download

RUN CGO_ENABLED=0 go build -o go-products-review

CMD ["./go-products-review"]