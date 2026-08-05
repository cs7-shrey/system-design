import pika

conn = pika.BlockingConnection(pika.ConnectionParameters("localhost"))
ch = conn.channel()
ch.queue_declare(queue="hello", durable=True)

def callback(ch, method, properties, body):
    print(f"Received message: {body.decode()}")

ch.basic_qos(prefetch_count=1) # prefetch count is the number of messages to prefetch
ch.basic_consume(queue="hello", on_message_callback=callback, auto_ack=True)

print("Waiting for messages. To exit press CTRL+C")
ch.start_consuming()
