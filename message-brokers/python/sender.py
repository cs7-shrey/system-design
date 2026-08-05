import pika

conn = pika.BlockingConnection(pika.ConnectionParameters("localhost"))
ch = conn.channel()

ch.queue_declare(queue="hello", durable=True)

ch.basic_publish(
    exchange="",        # direct exchange
    routing_key="hello",        # with default exchange, the routing key is the queue name
    body="Hello, World!",
    properties=pika.BasicProperties(
        delivery_mode=2,
    ),
)

print("Message published")
conn.close()