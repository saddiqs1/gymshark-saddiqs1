## TODO
- [ ] readme
- [x] logic for packs
  - [ ] reivisit pack logic eventually
- [x] tests for pack logic
- [x] local api setup
- [x] setup terraform for infra
- [x] ci/cd
- [ ] lambda api entry point setup

STRETCH
- [ ] ability to manage pack sizes
- [ ] frontend
- [ ] deploy frontend with infra

# Gymshark Coding Test

TODO - fill this out

## Quick Start

TODO - update this

-   Ensure `go 1.25` is installed

## Project Structure

TODO - updated structure

```js
gymshark-saddiqs1
├── backend/
│   ├── cmd/  // entry points into the code
│   │   ├── local
│   │   └── lambda
│   ├── config/  // configures env variables
│   ├── internal/ // code that is specific to this repo
│   └── pkg/ // code that can be extracted into a shared package for other repos
│ 
├── frontend/ //...
└── infra/
```

## Running Locally

TODO - 'proper' api docs here

```bash
cd backend
go run cmd/local/main.go

curl http://localhost:8080/health
curl http://localhost:8080/packs?itemsOrdered=1500
```

## Debug

-   In VSCode, use the `Run and Debug` mode, which has 2 preset configs in order to debug either the lambda, or the local api.

## Deployments

TODO


