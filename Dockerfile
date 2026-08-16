FROM golang:1.25 AS build

WORKDIR /gymshark-saddiqs1
COPY ./backend/go.mod ./backend/go.sum ./
COPY ./backend/cmd/lambda ./cmd/lambda/
COPY ./backend/config ./config/
COPY ./backend/internal ./internal/
COPY ./backend/pkg ./pkg/
RUN GOARCH=arm64 GOOS=linux go build -o main ./cmd/lambda/main.go

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:1.0.0 /lambda-adapter /opt/extensions/lambda-adapter
COPY --from=build /gymshark-saddiqs1/main ./main
ENTRYPOINT [ "./main" ]
