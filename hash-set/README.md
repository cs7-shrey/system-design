# Hash Set

A hash set implementation in Go for storing unique strings. It uses FNV-1a hashing, separate chaining for collisions, and automatically doubles its bucket count when needed.

Supported operations:

- `Insert`
- `Contains`
- `Delete`
- `Size`

Operations have expected **O(1)** time complexity and **O(n)** worst-case complexity when many keys collide. Insertion is O(1) amortized because resizing occasionally takes O(n).

## Testing

Run the correctness tests:

```bash
go test ./...
```

Run the performance benchmarks:

```bash
go test -run '^$' -bench . -benchmem
```

The benchmarks compare operation times at different set sizes. Times that remain roughly constant as the set grows are consistent with expected O(1) behavior.
