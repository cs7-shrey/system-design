# Bloom Filters — First-Principles Revision Notes

This note is meant to rebuild the idea of a Bloom filter from something more familiar: a hash set / `unordered_set`.

The goal is not to memorize formulas first. The goal is to remember **why the data structure exists, what information it throws away, and what tradeoff that creates**.

---

## 1. Start from a normal hash set

Suppose I want a set of strings with operations like:

```text
Insert(x)
Contains(x)
Delete(x)
```

A basic chained hash set can look like:

```text
array of buckets

[0] -> nil
[1] -> "alice" -> "bob"
[2] -> nil
[3] -> "charlie"
...
```

For a key `x`:

```text
hash(x) -> integer
bucket = hash(x) % number_of_buckets
```

If multiple values land in the same bucket, we keep a linked list / chain there.

If the number of elements is `n` and number of buckets is `m`, then the expected chain length is roughly:

```text
n / m
```

This is the load factor:

```text
alpha = n / m
```

If we keep `alpha = O(1)` by resizing the bucket array as the set grows, then expected lookup remains `O(1)`.

When resizing from, say, 16 buckets to 32 buckets, we cannot copy chains verbatim because:

```text
hash(x) % 16
```

and

```text
hash(x) % 32
```

may point to different buckets.

So every element must be rehashed and reassigned.

The total resize work across capacities like:

```text
16 + 32 + 64 + 128 + ... + n
```

is still `O(n)`, so over `n` insertions the amortized insertion cost remains `O(1)`.

So a normal hash set gives us:

```text
Contains: expected O(1)
Insert:   expected amortized O(1)
Delete:   expected O(1)
Memory:   O(n)
```

The important thing is: **the set stores the actual values**.

---

## 2. What if I don't need the actual values?

Suppose I only care about this question:

```text
"Could this value have been seen before?"
```

I don't need to retrieve the original string.

I don't even need exact membership in every case.

Then storing the actual strings, nodes, pointers, bucket chains, etc. may be wasteful.

This motivates the Bloom filter.

The basic idea is:

> Remove the actual values and keep only a tiny amount of evidence that a value passed through the structure.

Instead of buckets containing strings, keep an array of bits:

```text
index:  0 1 2 3 4 5 6 7 8 9
bits:   0 0 0 0 0 0 0 0 0 0
```

---

## 3. Bloom filter with one hash function

Start with only one hash function.

Suppose:

```text
index(x) = hash(x) % 10
```

Insert:

```text
alice   -> 3
bob     -> 7
charlie -> 1
```

Then the bit array becomes:

```text
index:  0 1 2 3 4 5 6 7 8 9
bits:   0 1 0 1 0 0 0 1 0 0
```

Notice what is missing:

```text
"alice"
"bob"
"charlie"
```

None of the actual strings are stored.

Only some bits were flipped to `1`.

### Querying

Suppose:

```text
hash("david") % 10 = 5
```

and bit 5 is `0`.

Then `"david"` is **definitely not present**.

Why?

Because if it had ever been inserted, insertion would have set bit 5 to `1`.

Now suppose:

```text
hash("eve") % 10 = 3
```

and bit 3 is `1`.

Can I conclude that `"eve"` was inserted?

No.

Maybe `"alice"` set that bit.

So the Bloom filter can only say:

```text
bit = 0 -> definitely not present
bit = 1 -> maybe present
```

This asymmetry is the whole point.

A Bloom filter can have:

```text
false positives: yes
false negatives: no
```

assuming a normal insert-only Bloom filter and correct implementation.

---

## 4. False positive rate with one hash function

Let:

```text
m = number of bits
n = number of inserted elements
```

Pick one particular bit.

For one insertion, the probability that the hash does **not** choose that bit is:

```text
1 - 1/m
```

After `n` independent insertions, the probability that the bit is still `0` is:

```text
(1 - 1/m)^n
```

Therefore the probability that the bit is `1` is:

```text
q = 1 - (1 - 1/m)^n
```

Now consider a query for a key that we **know was never inserted**.

With one hash function, the query picks one random bit.

If that bit is already `1`, the Bloom filter incorrectly says "maybe present".

So for one hash function:

```text
false positive rate = q
```

or:

```text
FP = 1 - (1 - 1/m)^n
```

Important distinction:

This is not:

```text
P(Bloom filter returns yes)
```

because that would include true positives.

It is:

```text
P(Bloom filter returns yes | queried key is actually absent)
```

---

## 5. Why one hash function is not enough

With only one hash function, an absent key only needs to collide with **one already-set bit**.

Example:

```text
alice -> bit 7
```

Later:

```text
eve -> bit 7
```

Even if `eve` was never inserted, the filter says "maybe".

A natural improvement is:

> Make each key leave more than one piece of evidence.

That means using multiple hash functions.

---

## 6. Two hash functions

Now each inserted key sets two positions.

Example:

```text
alice -> h1 -> 17
alice -> h2 -> 63
```

So insertion sets both bits:

```text
bit[17] = 1
bit[63] = 1
```

For lookup, an absent key is considered "maybe present" only if **both** of its hash positions are already `1`.

Example:

```text
eve -> h1 -> 17
       h2 -> 81
```

If:

```text
bit[17] = 1
bit[81] = 0
```

then `eve` is definitely absent.

Only this case produces a false positive:

```text
bit[h1(eve)] = 1
AND
bit[h2(eve)] = 1
```

So two hashes make false positives harder.

But there is a tradeoff: every inserted key now sets **two bits instead of one**, so the array fills up faster.

---

## 7. False positive rate with two hashes

Let:

```text
m = number of bits
n = number of inserted elements
k = 2
```

There are now:

```text
2n
```

hash placements during insertion.

For a particular bit, the probability that one hash placement misses it is:

```text
1 - 1/m
```

The probability that all `2n` placements miss it is:

```text
(1 - 1/m)^(2n)
```

Therefore the probability that a particular bit is `1` is:

```text
q = 1 - (1 - 1/m)^(2n)
```

Now take an absent query.

It produces two independently distributed positions.

For a false positive, both positions need to already be set.

So approximately:

```text
FP ~= q^2
```

Therefore:

```text
FP ~= [1 - (1 - 1/m)^(2n)]^2
```

The approximation assumes the hash positions behave independently enough for the model.

---

## 8. General case: k hash functions

Now each inserted key sets `k` positions.

Across `n` inserted keys, there are:

```text
kn
```

bit-setting attempts.

The probability that one particular bit is still `0` is:

```text
(1 - 1/m)^(kn)
```

So the probability that it is `1` is:

```text
q = 1 - (1 - 1/m)^(kn)
```

An absent query checks `k` positions.

For a false positive, all `k` positions must already be `1`.

Therefore:

```text
FP ~= q^k
```

or:

```text
FP ~= [1 - (1 - 1/m)^(kn)]^k
```

This is the core Bloom filter equation.

---

## 9. Why not use infinitely many hash functions?

At first it sounds like more hashes should always be better.

But increasing `k` has two opposing effects.

### Good effect

A query needs more positions to all be `1`:

```text
1 matching bit
2 matching bits
3 matching bits
...
```

That makes false positives less likely.

### Bad effect

Every insertion also sets more bits:

```text
k = 1 -> 1 bit per item
k = 2 -> 2 bits per item
k = 7 -> 7 bits per item
k = 20 -> 20 bits per item
```

Eventually the bit array becomes mostly:

```text
111111111111111111111111111111
```

At that point almost every query returns "maybe".

So there is an optimal value of `k` for a given:

```text
m = bit array size
n = expected number of inserted items
```

A common approximation is:

```text
k_opt ~= (m / n) * ln(2)
```

I do not need to memorize this unless I am actually sizing a Bloom filter.

The important intuition is:

```text
more hashes -> stricter lookup condition
but
more hashes -> fills array faster
```

There is a sweet spot.

---

## 10. Approximation commonly used in Bloom filter math

For large `m`, this expression:

```text
(1 - 1/m)^(kn)
```

is commonly approximated by:

```text
e^(-kn/m)
```

So:

```text
q ~= 1 - e^(-kn/m)
```

and:

```text
FP ~= (1 - e^(-kn/m))^k
```

This is usually the formula shown in Bloom filter explanations.

But it is just a cleaner approximation of the probability argument above.

---

## 11. Memory intuition

This is where Bloom filters become useful.

A normal hash set stores actual values:

```text
string bytes
bucket array
pointers / nodes / metadata
allocation overhead
```

A Bloom filter stores only:

```text
m bits
```

For example, if a Bloom filter uses 10 million bits:

```text
10,000,000 bits
/ 8
= 1,250,000 bytes
~ 1.25 MB
```

That may represent approximate membership for millions of values.

The price is that it cannot distinguish:

```text
"this exact key set these bits"
```

from:

```text
"other keys happened to set all the same bits"
```

That is the source of false positives.

---

## 12. Hash set vs Bloom filter

### Hash set

```text
key
 |
 v
hash
 |
 v
bucket
 |
 v
actual key stored
```

Properties:

```text
exact membership
stores actual keys
supports normal deletion
higher memory usage
```

### Bloom filter

```text
               /-> h1 -> bit
key -> hashes ----> h2 -> bit
               \-> h3 -> bit
```

Properties:

```text
does not store actual keys
very memory efficient
false positives possible
false negatives not possible
normal deletion is not straightforward
```

---

## 13. Why normal Bloom filters cannot simply delete

Suppose:

```text
alice -> bits 3, 7, 10
bob   -> bits 2, 7, 14
```

Both `alice` and `bob` depend on bit 7.

If I delete `alice` and do:

```text
bit[7] = 0
```

then I accidentally destroy evidence needed for `bob`.

So standard Bloom filter bits cannot tell us **how many keys caused a bit to be set**.

A variant called a **counting Bloom filter** replaces bits with small counters so deletion becomes possible, at the cost of more memory.

---

## 14. Multiple hash functions in practice

Conceptually, a Bloom filter has:

```text
h1(x)
h2(x)
h3(x)
...
hk(x)
```

These should behave like independent random mappings into the bit array.

In production, implementations do not necessarily run `k` completely unrelated expensive hash algorithms.

A common technique is to compute a small number of base hashes and derive the `k` positions from them.

For understanding the data structure, it is enough to think:

```text
one key -> k independent-looking positions
```

---

## 15. The mental model I actually want to remember

A Bloom filter is basically a hash set after throwing away the expensive part: **the actual keys**.

A hash set says:

```text
"I hashed this key to a bucket, and I can inspect the real values there."
```

A Bloom filter says:

```text
"I no longer have the real values.
I only have evidence that some inserted values touched these positions."
```

With one hash:

```text
one key -> one bit
```

which creates too many accidental matches.

With multiple hashes:

```text
one key -> several bits
```

and a query returns "maybe" only if all of those bits are already set.

This lowers false positives until the extra hashes start filling the array too aggressively.

So the core knobs are:

```text
m = number of bits
n = number of inserted elements
k = number of hash functions
```

and the approximate false-positive rate is:

```text
FP ~= [1 - (1 - 1/m)^(kn)]^k
```

or, using the common exponential approximation:

```text
FP ~= (1 - e^(-kn/m))^k
```

The useful guarantee is:

```text
Bloom says NO    -> definitely absent
Bloom says MAYBE -> possibly present
```

That is enough to recognize when a Bloom filter is useful.

---

## 16. When should I think of using one?

A Bloom filter is useful when:

```text
- there are lots of keys
- exact storage of all keys is expensive
- checking the real source is expensive
- rejecting definite misses cheaply is valuable
- occasional false positives are acceptable
```

Typical pattern:

```text
request
   |
   v
Bloom filter
   |
   +-- definitely absent -> stop immediately
   |
   +-- maybe present -> check the real database / disk / network source
```

The Bloom filter is therefore often a **cheap gate in front of an expensive exact lookup**.

It does not replace the source of truth when exact membership matters.

