# Test Task – Go Backend

This repository contains my solution for the backend test task.

## Structure

- `domain/` – business logic (Order entity)
- `contracts/` – repository interfaces
- `usecases/` – application layer (use cases)
- `repo/` – repository implementation
- `REVIEW.md` – code review with identified issues
- `ANSWERS.md` – answers to theoretical questions

## Notes

- The solution focuses on clean architecture principles
- Domain layer is independent from infrastructure
- Concurrency issues were fixed
- Transaction integrity is ensured for critical operations

## Build

To verify the project builds:

```bash
go build ./...
