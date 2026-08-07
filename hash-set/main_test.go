package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNewHashSet(t *testing.T) {
	hs := NewHashSet()
	if hs == nil {
		t.Fatal("NewHashSet returned nil")
	}
	if got := hs.Size(); got != 0 {
		t.Fatalf("new set size = %d, want 0", got)
	}
	if hs.Contains("missing") {
		t.Fatal("new set contains an element that was never inserted")
	}
}

func TestInsertContainsAndDuplicates(t *testing.T) {
	hs := NewHashSet()
	values := []string{"", "alpha", "beta", "with spaces", "世界", "alpha"}

	for _, value := range values {
		hs.Insert(value)
		if !hs.Contains(value) {
			t.Fatalf("set does not contain %q immediately after insertion", value)
		}
	}

	// "alpha" was inserted twice, so it must only be counted once.
	if got, want := hs.Size(), len(values)-1; got != want {
		t.Fatalf("size after insertions = %d, want %d", got, want)
	}
}

func TestDelete(t *testing.T) {
	hs := NewHashSet()
	for _, value := range []string{"alpha", "beta", "gamma"} {
		hs.Insert(value)
	}

	hs.Delete("beta")
	if hs.Contains("beta") {
		t.Fatal("deleted value is still present")
	}
	if got := hs.Size(); got != 2 {
		t.Fatalf("size after deletion = %d, want 2", got)
	}

	// Deleting absent values and deleting the same value twice are no-ops.
	hs.Delete("beta")
	hs.Delete("never-inserted")
	if got := hs.Size(); got != 2 {
		t.Fatalf("no-op deletions changed size to %d, want 2", got)
	}

	for _, value := range []string{"alpha", "gamma"} {
		hs.Delete(value)
	}
	if got := hs.Size(); got != 0 {
		t.Fatalf("size after deleting all values = %d, want 0", got)
	}
}

func TestCollisions(t *testing.T) {
	hs := NewHashSet()
	values := stringsInSameBucket(10, len(hs.container))

	for _, value := range values {
		hs.Insert(value)
	}
	for _, value := range values {
		if !hs.Contains(value) {
			t.Errorf("colliding value %q was lost", value)
		}
	}

	// Delete the head, middle, and tail of the collision chain.
	for _, index := range []int{len(values) - 1, len(values) / 2, 0} {
		hs.Delete(values[index])
		if hs.Contains(values[index]) {
			t.Errorf("colliding value %q remains after deletion", values[index])
		}
	}
}

func TestResizePreservesElements(t *testing.T) {
	hs := NewHashSet()
	const count = 1_000

	for i := 0; i < count; i++ {
		hs.Insert(fmt.Sprintf("key-%d", i))
	}
	if got := hs.Size(); got != count {
		t.Fatalf("size after resize = %d, want %d", got, count)
	}
	for i := 0; i < count; i++ {
		value := fmt.Sprintf("key-%d", i)
		if !hs.Contains(value) {
			t.Fatalf("%q was lost during resize", value)
		}
	}
}

func TestAgainstReferenceSet(t *testing.T) {
	hs := NewHashSet()
	want := make(map[string]struct{})
	rng := rand.New(rand.NewSource(1))

	for step := 0; step < 10_000; step++ {
		value := fmt.Sprintf("value-%d", rng.Intn(2_000))
		switch rng.Intn(3) {
		case 0:
			hs.Insert(value)
			want[value] = struct{}{}
		case 1:
			hs.Delete(value)
			delete(want, value)
		case 2:
			_, expected := want[value]
			if got := hs.Contains(value); got != expected {
				t.Fatalf("step %d: Contains(%q) = %t, want %t", step, value, got, expected)
			}
		}
		if got := hs.Size(); got != len(want) {
			t.Fatalf("step %d: size = %d, want %d", step, got, len(want))
		}
	}
}

// Benchmarks are the appropriate way to evaluate complexity in Go. Run:
//
//	go test -run '^$' -bench . -benchmem
//
// Compare ns/op as the initial set size grows. Insert and Contains should stay
// roughly constant on average; a single resize may occasionally make Insert
// linear, but insertion should be O(1) amortized.
func BenchmarkHashSet(b *testing.B) {
	for _, size := range []int{16, 1_024, 65_536} {
		b.Run(fmt.Sprintf("Contains/size=%d", size), func(b *testing.B) {
			hs := populatedSet(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if !hs.Contains(fmt.Sprintf("key-%d", i%size)) {
					b.Fatal("existing key not found")
				}
			}
		})

		b.Run(fmt.Sprintf("Insert/size=%d", size), func(b *testing.B) {
			hs := populatedSet(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				hs.Insert(fmt.Sprintf("new-%d-%d", size, i))
			}
		})

		b.Run(fmt.Sprintf("DeleteAndReinsert/size=%d", size), func(b *testing.B) {
			keys := make([]string, size)
			for i := range keys {
				keys[i] = fmt.Sprintf("key-%d", i)
			}
			hs := populatedSet(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := keys[i%size]
				hs.Delete(key)
				hs.Insert(key) // Keep the set size stable for the next iteration.
			}
		})
	}
}

func populatedSet(size int) *HashSet {
	hs := NewHashSet()
	for i := 0; i < size; i++ {
		hs.Insert(fmt.Sprintf("key-%d", i))
	}
	return hs
}

func stringsInSameBucket(count, buckets int) []string {
	values := make([]string, 0, count)
	target := -1
	for i := 0; len(values) < count; i++ {
		value := fmt.Sprintf("collision-%d", i)
		bucket := int(hash(value) % uint64(buckets))
		if target == -1 {
			target = bucket
		}
		if bucket == target {
			values = append(values, value)
		}
	}
	return values
}
