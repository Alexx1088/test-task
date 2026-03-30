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

## Issue 8: Data race on shared error variable
- File: usecases/bulk_complete/interactor.go
- Line: writes to `lastErr` inside goroutines
- Severity: CRITICAL
- Problem: Multiple goroutines write to the shared `lastErr` variable without synchronization.
- Impact: This causes a data race. In production, it can return nondeterministic errors, overwrite failures, or even produce corrupted interface values under concurrent writes.
- Fix: Protect shared state with synchronization or collect errors safely through mutex-protected aggregation, channels, or errgroup patterns.

## Other Issues

## Issue 9: Unbounded concurrency
- File: usecases/bulk_complete/interactor.go
- Line: one goroutine started per order
- Severity: WARNING
- Problem: The code starts one goroutine per order without any concurrency limit.
- Impact: A large batch can overload the database, exhaust connection pools, and increase memory usage.
- Fix: Use a semaphore or worker pool to limit concurrency.

## Issue 10: Loss of error information in bulk processing
- File: usecases/bulk_complete/interactor.go
- Line: `return lastErr`
- Severity: WARNING
- Problem: Only a single error is returned even if multiple operations fail.
- Impact: This hides the full failure picture and makes debugging harder.
- Fix: Aggregate all errors and return them in a structured way.

## Transaction Integrity Bugs

## Issue 11: Order update and audit log creation are not atomic
- File: usecases/transfer_order/interactor.go
- Line: `Update(...)` followed by `CreateAuditLog(...)`
- Severity: CRITICAL
- Problem: The order transfer and audit log creation are performed as two separate operations without a transaction.
- Impact: If the order update succeeds but audit log creation fails, the database is left in an inconsistent state: the order is transferred but the audit trail is missing.
- Fix: Execute both operations atomically in a single transaction.

## Architecture Violations

## Issue 12: Usecase directly mutates aggregate fields
- File: usecases/transfer_order/interactor.go
- Line: `order.CustomerID = newCustomerID`
- Severity: WARNING
- Problem: The usecase modifies aggregate fields directly instead of using domain behavior.
- Impact: This bypasses domain invariants and makes business rules harder to enforce consistently.
- Fix: Introduce a domain method such as `TransferTo(newCustomerID string) error`.
