ARG BUILDPLATFORM
ARG TARGETPLATFORM

FROM --platform=$BUILDPLATFORM golang:1.25 AS build

ARG TARGETARCH
ARG TARGETOS

WORKDIR /gymshark-saddiqs1
COPY ./backend/go.mod ./backend/go.sum ./
COPY ./backend/cmd/api ./cmd/api/
COPY ./backend/config ./config/
COPY ./backend/internal ./internal/
COPY ./backend/pkg ./pkg/
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS go build -o main ./cmd/api/main.go

FROM --platform=$TARGETPLATFORM public.ecr.aws/awsguru/aws-lambda-adapter:1.0.1 AS lambda-adapter

FROM --platform=$TARGETPLATFORM public.ecr.aws/lambda/provided:al2023
COPY --from=lambda-adapter /lambda-adapter /opt/extensions/lambda-adapter
COPY --from=build /gymshark-saddiqs1/main ./main
ENTRYPOINT [ "./main" ]
