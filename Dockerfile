# stage 1: build
FROM golang:1.16 as build
LABEL stage=intermediate
WORKDIR /app
COPY . .
RUN make build

# stage 2: scratch
FROM gcr.io/distroless/static as base
COPY --from=build /app/bin/honeypot /bin/honeypot
CMD ["honeypot"]
