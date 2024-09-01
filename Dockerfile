FROM golang:1.22.0 AS build-stage

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
COPY templates/ templates/
COPY static/ static/

RUN CGO_ENABLED=0 GOOS=linux go build -o /rock_paper_scissors

# Test stage
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy stage
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /rock_paper_scissors /rock_paper_scissors
COPY --from=build-stage /app/templates /templates
COPY --from=build-stage /app/static /static

EXPOSE 8085

USER nonroot:nonroot

ENTRYPOINT ["/rock_paper_scissors"]
