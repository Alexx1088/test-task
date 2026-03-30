# Code Review

## Architecture Violations

## Issue 1: Domain depends on infrastructure
- File: domain/order.go
- Line: Complete method signature
- Severity: CRITICAL
- Problem: The domain method accepts `context.Context` and `*sql.DB`, which introduces infrastructure concerns into the domain layer.
- Impact: This breaks layer boundaries and makes the domain harder to test, reuse, and maintain.
- Fix: Change the method signature to a pure domain method: `func (o *Order) Complete() error`.

## Issue 2: Domain performs persistence
- File: domain/order.go
- Line: `db.ExecContext(...)`
- Severity: CRITICAL
- Problem: The entity writes directly to the database.
- Impact: This mixes business logic with repository responsibility and violates Clean Architecture.
- Fix: The domain entity should only change its own state. Saving changes must be done by the repository.

## Issue 3: Logging inside domain
- File: domain/order.go
- Line: `log.Printf("Order %s completed", o.ID)`
- Severity: WARNING
- Problem: Logging is done inside the domain entity.
- Impact: This couples business logic to infrastructure concerns and makes domain code less pure.
- Fix: Move logging to the application/service/decorator layer.

## Other Issues

## Issue 4: Missing import for errors package
- File: domain/order.go
- Line: `errors.New("cannot complete")`
- Severity: WARNING
- Problem: The code uses `errors.New(...)` but does not import the `errors` package.
- Impact: The code does not compile.
- Fix: Add the missing import or redesign error handling.

## Issue 5: float64 used for money
- File: domain/order.go
- Line: `Total float64`
- Severity: WARNING
- Problem: Monetary values are stored using `float64`, which can lead to precision errors.
- Impact: This may cause incorrect totals or rounding issues in financial calculations.
- Fix: Store money in minor units like `int64` (for example, cents) or use a decimal type.

## Issue 6: Usecase depends on concrete repository implementation
- File: usecases/complete_order/interactor.go
- Line: Interactor struct field `repo *repo.OrderRepo`
- Severity: CRITICAL
- Problem: The usecase depends on a concrete repository implementation instead of an interface.
- Impact: This breaks dependency inversion and makes testing and maintenance harder.
- Fix: Depend on a repository contract interface from the contracts layer.

## Issue 7: Usecase leaks infrastructure through DB access
- File: usecases/complete_order/interactor.go
- Line: `order.Complete(ctx, uc.repo.DB())`
- Severity: CRITICAL
- Problem: The usecase passes a database handle into the domain method.
- Impact: This leaks infrastructure concerns across architectural boundaries.
- Fix: Let the domain mutate state only, and persist changes through repository methods.

## Concurrency Bugs

## Transaction Integrity Bugs

## Other Issues