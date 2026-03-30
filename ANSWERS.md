# Answers

## Q1: Race condition in bulk_complete

In the original implementation, multiple goroutines write to the shared variable `lastErr` without synchronization.

This is not just a "lost error" problem.

In Go, `error` is an interface type, which internally contains a type pointer and a data pointer. Concurrent unsynchronized writes to an interface value can lead to a corrupted interface state.

Possible issues:
- nondeterministic final error value
- overwriting errors from other goroutines
- corrupted interface value
- race detector failures
- unpredictable behavior in production

The correct fix is to synchronize access to shared state (e.g., using mutex) or collect errors via channels or structured aggregation.

---

## Q2: Transaction integrity in transfer_order

If `Update` succeeds but `CreateAuditLog` fails:

- The order is already updated to the new customer
- The audit log entry is missing
- The database is left in an inconsistent state

This breaks data integrity because related changes are not applied atomically.

Fix:

Both operations should be executed within a single transaction. This can be implemented by grouping mutations into a transactional unit (Plan pattern) and executing them atomically.

In practice, this can be done using a transaction boundary like:

- begin transaction
- update order
- create audit log
- commit or rollback

---

## Q3: Domain method signature violation

The method:

func (o *Order) Complete(ctx context.Context, db *sql.DB) error

is a violation of Clean Architecture.

Reasons:
- `context.Context` is an infrastructure concern
- `*sql.DB` is a database implementation detail
- the domain layer must remain pure and independent of external systems

Correct approach:

The domain method should only contain business logic:

func (o *Order) Complete() error

All persistence should be handled in the repository layer.

---

## Q4: Ignoring error from Retrieve

Code:

source, _ := uc.repo.Retrieve(ctx, req.FromAccountID)

This hides more than just a potential nil value.

Hidden problems:
- database connection errors
- timeouts or context cancellations
- "not found" vs "temporary failure" distinction is lost
- authorization or filtering issues may be ignored
- stale or inconsistent reads are not detected

Ignoring errors makes debugging extremely difficult and can lead to incorrect system behavior.

All errors must be handled explicitly and propagated properly.