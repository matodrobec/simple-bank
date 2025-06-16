FROM golang:1.24-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go


FROM alpine:3.22
WORKDIR /app

COPY --from=builder /app/db/migration ./db/migration
COPY --from=builder /app/main .
COPY ./docker/images/app/start.sh .
COPY ./wait-for.sh .
COPY app.env .

RUN chmod +x start.sh wait-for.sh

EXPOSE 8080
EXPOSE 8090
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]