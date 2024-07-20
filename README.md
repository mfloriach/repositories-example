# Repositories example

[![test](https://github.com/mfloriach/repositories-example/actions/workflows/test.yml/badge.svg)](https://github.com/mfloriach/repositories-example/actions/workflows/test.yml)
![go](https://img.shields.io/badge/go-1.21-blue)
[![Go Coverage](https://github.com/mfloriach/repositories-example/wiki/coverage.svg)](https://raw.githack.com/wiki/USER/REPO/coverage.html)



Small project to show how to achieve polymorphism with repositories and design flexible interfaces to filter resuts.

This project use:
- Mysql
- MongoDB

Usage:
```bash
go mod download
go test -run ^TestUser repos/repositories
```