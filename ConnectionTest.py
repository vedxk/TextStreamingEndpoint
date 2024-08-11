import asyncio
import websockets
import random

async def connect_and_communicate(uri, message, client_id):
    async with websockets.connect(uri, ping_interval=10, ping_timeout=5) as websocket:
        print(f"Client {client_id} connected to the server.")

        try:
            while True:
                await websocket.send(message)
                print(f"Client {client_id} sent: {message}")
                
                response = await websocket.recv()
                print(f"Client {client_id} received: {response}")

                await asyncio.sleep(5) 
        except websockets.exceptions.ConnectionClosedError as e:
            print(f"Client {client_id} connection closed with error: {e}")
        except asyncio.TimeoutError:
            print(f"Client {client_id} timed out waiting for a response.")
        except Exception as e:
            print(f"Client {client_id} encountered an error: {e}")

async def main():
    uri = "ws://localhost:8080/ws" 
    clients = 100 
    messages = [
        "What is your name?",
        "How are you?",
        "Hi",
        "Unknown question"
    ]

    tasks = []
    for i in range(clients):
        message = random.choice(messages)
        # Schedule task
        tasks.append(connect_and_communicate(uri, message, i + 1))
    
    # Run all client tasks concurrently
    await asyncio.gather(*tasks)

if __name__ == "__main__":
    asyncio.run(main())

