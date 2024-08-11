import asyncio
import websockets
import random

async def connect_and_send_message(uri, message, client_id):
    async with websockets.connect(uri) as websocket:
        print(f"Client {client_id} connected to the server.")
        
        # Send message to the server
        await websocket.send(message)
        print(f"Client {client_id} sent: {message}")
        
        # Wait for response from the server
        try:
            response = await websocket.recv()
            print(f"Client {client_id} received: {response}")
        except websockets.exceptions.ConnectionClosedError as e:
            print(f"Client {client_id} connection closed with error: {e}")
        except asyncio.TimeoutError:
            print(f"Client {client_id} timed out waiting for a response.")

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
        tasks.append(connect_and_send_message(uri, message, i + 1))
    
    # Run all client concurrently
    await asyncio.gather(*tasks)

if __name__ == "__main__":
    asyncio.run(main())
