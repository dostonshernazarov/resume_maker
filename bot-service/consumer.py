import os
import pika
import asyncio
import json
from aiogram import Bot, Dispatcher, types
import time

RABBITMQ_HOST = os.getenv('RABBITMQ_HOST', 'rabbitmq')
RABBITMQ_PORT = int(os.getenv('RABBITMQ_PORT', 5672))
RABBITMQ_USER = os.getenv('RABBITMQ_USER', 'guest')
RABBITMQ_PASS = os.getenv('RABBITMQ_PASS', 'guest')

TELEGRAM_BOT_TOKEN = os.getenv('TELEGRAM_BOT_TOKEN')
TELEGRAM_CHAT_ID = os.getenv('TELEGRAM_CHAT_ID')

bot = Bot(token=TELEGRAM_BOT_TOKEN)
dp = Dispatcher(bot)

async def send_to_telegram(message):
    await bot.send_message(chat_id=TELEGRAM_CHAT_ID, text=message, parse_mode='HTML')

def callback(ch, method, properties, body):
    try:
        bot_produce = json.loads(body)
        message = f"""
<b>🆕 NEW RESUME:</b>

<b>👱🏻‍♂️ Employee:</b> {bot_produce['name']}
<b>📧 Email:</b> {bot_produce['email']}
<b>📞 Phone:</b> {bot_produce['phone']}
<b>📚 Job:</b> {bot_produce['label']}
<b>🏠 City:</b> {bot_produce['location']}
<b>💵 Salary:</b> ${bot_produce['salary']}

<b>📄 Resume:</b>
{bot_produce['url']}

<b>🔗 Profiles:</b>
{', '.join([profile['url'] for profile in bot_produce['profiles']])}

<b>🔎 Summary:</b>
{bot_produce['summary']}
"""
        print(f" [x] Received message")
        asyncio.run(send_to_telegram(message))
    except json.JSONDecodeError:
        print(f"Failed to decode JSON: {body}")
    except KeyError as e:
        print(f"Missing key in JSON: {e}")

def main():
    connection = None
    for attempt in range(10):
        try:
            connection = pika.BlockingConnection(
                pika.ConnectionParameters(
                    host=RABBITMQ_HOST,
                    port=RABBITMQ_PORT,
                    credentials=pika.PlainCredentials(RABBITMQ_USER, RABBITMQ_PASS)
                )
            )
            break
        except pika.exceptions.AMQPConnectionError:
            print(f"Connection attempt {attempt + 1} failed. Retrying in 5 seconds...")
            time.sleep(5)

    if connection is None:
        print("Failed to connect to RabbitMQ after 10 attempts. Exiting.")
        return

    channel = connection.channel()

    channel.queue_declare(queue='cvmaker_queue')

    channel.basic_consume(queue='cvmaker_queue', on_message_callback=callback, auto_ack=True)

    print(' [*] Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()

if __name__ == '__main__':
    main()