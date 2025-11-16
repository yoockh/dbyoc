# utils — Utility helpers for DBYOC

This package provides lightweight utility helpers used across the DBYOC project: retry and backoff logic, simple map/struct conversion helpers, and a small SQL query builder helper. These utilities are intentionally small, dependency-free (standard library only), and meant to be reusable across SQL and NoSQL client implementations in the repository.

- Language: Go
- Location: `./utils`

## Contents / Files
- `backoff.go` — Exponential backoff helper (Backoff struct and helpers).
- `retry.go` — Retry helper with simple configuration (RetryConfig, Retry).
- `helpers.go` — Map/struct conversion helpers and a tiny QueryBuilder.

## Key types & functions

### Backoff (utils/backoff.go)
A small exponential backoff implementation for use with retry loops.

- type Backoff
  - InitialInterval time.Duration
  - Multiplier float64
  - MaxInterval time.Duration
  - MaxElapsedTime time.Duration

- func NewBackoff() *Backoff
  - Returns a Backoff with sane defaults (100ms initial, x2 multiplier, 30s max interval, 5m max elapsed).

- func (b *Backoff) GetNextInterval(attempt int) time.Duration
  - Returns the interval for the given attempt index, capped to MaxInterval.

- func (b *Backoff) IsElapsed(start time.Time) bool
  - Returns true if MaxElapsedTime has been exceeded since start.

### Retry (utils/retry.go)
Simple retry helper to run an operation with a fixed number of attempts and fixed delay.

- var ErrMaxRetriesExceeded error
  - Error returned when all retries are exhausted.

- type RetryConfig
  - MaxRetries int
  - Delay time.Duration

- func NewRetryConfig(maxRetries int, delay time.Duration) *RetryConfig

- func Retry(operation func() error, config *RetryConfig) error
  - Runs operation up to MaxRetries times, sleeping Delay between attempts.
  - If config is nil, a default of 3 retries and 100ms delay is used.

### Helpers (utils/helpers.go)
Utility helpers for map<->struct conversion and a tiny SQL QueryBuilder.

- func MapToStruct(m map[string]interface{}, s interface{}) error
  - Marshals the map to JSON and unmarshals into the target struct (`s` must be a pointer to a struct).

- func StructToMap(s interface{}) (map[string]interface{}, error)
  - Converts a struct (pointer) to a map by reflecting over exported fields.
  - Note: This uses reflect and returns field names as keys (no tag handling).

- type QueryBuilder
  - A tiny builder to create simple SELECT queries.
  - Methods:
    - NewQueryBuilder(table string) *QueryBuilder
    - (*QueryBuilder) Select(columns ...string) *QueryBuilder
    - (*QueryBuilder) Where(condition string) *QueryBuilder
    - (*QueryBuilder) Build() string

## Examples

Backoff example:
```go
package main

import (
	"fmt"
	"time"

	"github.com/yoockh/dbyoc/utils"
)

func main() {
	b := utils.NewBackoff()
	start := time.Now()
	for i := 0; i < 6; i++ {
		if b.IsElapsed(start) {
			fmt.Println("backoff elapsed")
			break
		}
		sleep := b.GetNextInterval(i)
		fmt.Printf("attempt %d, sleeping %v\n", i+1, sleep)
		time.Sleep(sleep)
	}
}
```

Retry example:
```go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/yoockh/dbyoc/utils"
)

func main() {
	cfg := utils.NewRetryConfig(5, 200*time.Millisecond)
	attempt := 0
	err := utils.Retry(func() error {
		attempt++
		if attempt < 3 {
			return errors.New("temporary failure")
		}
		return nil
	}, cfg)
	if err != nil {
		fmt.Println("operation failed:", err)
	} else {
		fmt.Println("operation succeeded after", attempt, "attempt(s)")
	}
}
```

MapToStruct / StructToMap example:
```go
type User struct {
	ID   int
	Name string
}

m := map[string]interface{}{"ID": 1, "Name": "Alice"}
var u User
if err := utils.MapToStruct(m, &u); err != nil {
    // handle error
}
out, _ := utils.StructToMap(&u)
fmt.Println(out) // map[string]interface{}{"ID":1, "Name":"Alice"}
```

QueryBuilder example:
```go
qb := utils.NewQueryBuilder("users").Select("id", "name").Where("active = true")
sql := qb.Build() // "SELECT id, name FROM users WHERE active = true"
```

## Notes and recommendations
- MapToStruct uses JSON round-tripping. It is simple and convenient but has limitations (e.g., type coercion rules of encoding/json). For high-performance or complex mappings consider using a dedicated mapper.
- StructToMap assumes a pointer to a struct and uses reflect to read exported fields. It does not read struct tags (like `json` or `db`)—you may want to extend it if tag-aware mapping is required.
- Retry is a basic fixed-delay retry. If you need jitter, exponential backoff, or more advanced behavior combine Retry with Backoff or a third-party library.
- Backoff and Retry are not context-aware. If you need cancellation or timeouts, add context.Context support around your operation and check for cancellation inside the operation function.

## Testing
- Add unit tests that cover:
  - Backoff.GetNextInterval behavior including capping at MaxInterval.
  - Retry behavior with transient failures and final failure.
  - MapToStruct and StructToMap round-trips and error cases.
  - QueryBuilder string generation.

## Contribution
Please open an issue or PR in the repository for improvements, bug fixes or feature requests. Follow repository contribution guidelines.

## License
See the project's LICENSE file at the repository root.
