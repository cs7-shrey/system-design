# RabbitMQ

## What is RabbitMQ?

Think of RabbitMQ as a programmable message-routing server between applications.

Your application does not normally send a message directly to another service or even directly to a queue. It sends the message to an exchange, along with a routing key:

```
Publisher
    |
    | publish(exchange, routingKey, message)
    v
 Exchange
    |
    | bindings decide where it goes
    v
 Queue(s)
    |
    | RabbitMQ delivers messages
    v
 Consumer(s)
```

The responsibilities are separated:

Publisher: creates a message and describes where it logically belongs.
Exchange: routes that message.
Queue: stores the message until it can be processed.
Consumer: receives and processes it.
Acknowledgement: tells RabbitMQ that processing succeeded.

RabbitMQ accepts messages from publishers, routes them through exchanges, and stores successfully routed messages in queues until consumers receive them.

RabbitMQ also has:

- a Management UI
- an HTTP management API
- metrics endpoints
- other protocols and clients
- a newer AMQP 1.0 client ecosystem
- RabbitMQ Streams



## Lifecycle of a message

Suppose an HTTP API receives:

```
POST /emails
```

Instead of sending the email synchronously, it publishes:

```json
{
  "email_id": "email-123",
  "recipient": "user@example.com"
}
```

The path is:

```
1. API serializes the message.
2. API publishes it to an exchange.
3. Exchange evaluates its bindings.
4. Message is copied into one or more matching queues.
5. RabbitMQ delivers it to a consumer.
6. Consumer sends the email.
7. Consumer acknowledges the message.
8. RabbitMQ removes it from the queue.
```

While a message is unacknowledged, it is unavailable to other consumers.
RabbitMQ remembers:

- which channel received it
- its delivery tag
- that ownership is temporarily with that consumer

Suppose Consumer 1 receives message M1:

```
Queue:
Ready:   M2, M3, M4
Unacked: M1 → Consumer 1
```

RabbitMQ then waits for one of these outcomes.

1. Consumer sends ACK

```
Consumer 1 finishes M1
Consumer 1 → ACK
RabbitMQ deletes M1
```

1. Consumer explicitly rejects it

The consumer can send:

```
NACK / REJECT with requeue=true
```

RabbitMQ puts it back into the queue, where Consumer 1 or another consumer may receive it.

Or:

```
NACK / REJECT with requeue=false
```

RabbitMQ discards it, or sends it to a dead-letter exchange if one is configured.

1. Consumer crashes or disconnects

If the process dies, its TCP connection disappears, or its channel closes before ACK, RabbitMQ automatically requeues all unacknowledged messages belonging to that channel.

```
M1 → Consumer 1
Consumer 1 crashes
M1 → queue again
M1 → Consumer 2
```

The redelivered message is marked with redelivered = true. Consumers should therefore be designed to handle the same message more than once.

1. The consumer remains alive but never ACKs

RabbitMQ cannot immediately know whether the consumer is:

- legitimately doing slow work
- stuck in an infinite loop
- deadlocked
- silently broken

So it waits for the configured consumer acknowledgement timeout. The current default is 30 minutes. When the timeout is triggered, RabbitMQ closes the affected channel and requeues its outstanding unacknowledged deliveries. The timeout is configurable through consumer_timeout.

Notice the distinction:

```
Published ≠ Routed ≠ Stored ≠ Delivered ≠ Processed
```

These are separate moments, with separate failure possibilities.

That distinction is basically the foundation of reliable RabbitMQ applications.

### Prefetch controls how many messages can become unacked

Without sensible prefetch, RabbitMQ might give one consumer many messages before it finishes processing them.

With:

```
prefetch = 1
```

RabbitMQ gives that consumer only one unacked message at a time:

```
Consumer gets M1
RabbitMQ waits for ACK
Consumer ACKs M1
RabbitMQ sends M2
```

With:

```
prefetch = 10
```

a consumer may hold up to ten unacknowledged messages simultaneously.

So the mental model is:

> Delivery temporarily leases the message to a consumer. ACK permanently completes it. Losing the consumer or timing out revokes the lease and makes the message available again.



## Key entities in RabbitMQ

- Publisher
- Exchange
- Queue
- Consumer
- Binding
- Routing key

The producer publishes a message to an exchange, usually with a routing key. The exchange uses its type and bindings to decide which queue or queues should receive that message. Consumers then read messages from those queues.

### Publisher

The publisher is the entity that sends messages to the exchange.

### Exchange

The exchange is the entity that routes messages to the queues.

### Queue

The queue is the entity that stores messages until they are delivered to a consumer.

### Consumer

The consumer is the entity that receives messages from the queue.

### Routing key

A routing key is a label the producer attaches to a message when publishing it.

### Binding

A binding is a rule connecting an exchange to a queue, telling the exchange which messages that queue should receive.

Example:

```
Producer publishes:
message = "Payment completed"
routing key = "payment.completed"
```

Bindings:

```
payments-queue   ← binding key: payment.completed
audit-queue      ← binding key: payment.*
```

The exchange compares the message’s routing key with each binding:

```
Exchange
   ├── payment.completed → payments-queue
   └── payment.*         → audit-queue
```

So both queues may receive the message.

The exact matching behavior depends on the exchange type:

- Direct exchange: routing key must exactly match the binding key.
- Topic exchange: binding keys can contain patterns such as payment.* or payment.#.
- Fanout exchange: routing keys are ignored; every bound queue receives the message.
- Headers exchange: routing is based on message headers rather than the routing key.



## Experiments

1. One queue, one sender, one receiver. Send and receive messages.
2. One sender, one queue, multiple consumers.
3. Add an exchange, multiple queues, multiple consumers
4. Add routing keys, binding keys, exchange types
