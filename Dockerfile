FROM golang:1.15-alpine AS build

WORKDIR /src/
COPY . /src/
RUN ls -alh
RUN CGO_ENABLED=0 go build -o /bin/location-service cmd/webapp/main.go

FROM scratch
COPY --from=build /bin/location-service /bin/location-service
ENTRYPOINT ["/bin/location-service"]
CMD ["-addr", ":5000", "-datastore", "inmemdb"]