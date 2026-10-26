# Todo App

![CI](https://github.com/Torrahp/Todo-App/actions/workflows/ci.yml/badge.svg)

React + Go (Gin) + PostgreSQL to-do app with JWT auth, tags, due dates, search and pagination.

## Run locally

    cp .env.example .env      # then edit the values
    make up                   # http://localhost:3000

## Test

    make test                 # unit tests (Go + Vitest)
    make test-integration     # repository tests against Postgres
