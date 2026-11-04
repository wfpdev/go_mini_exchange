<h1 align="center">Golang Mini Exchange System</h1>

<p align="center">
  A simple in-memory exchange system built with Go.
</p>

<h3>How To Run</h3>

Simply run using `docker compose up` command. Then `docker compose down` to turn it off.

Containers included:
- Golang app (http://localhost:8080)
- Redpanda Console GUI (http://localhost:8081/topics)
- Redpanda Message Broker
- Simple trading simulator (Python bot). Made from `scripts` folder.

<h3>Architecture</h3>

<img width="471" height="262" src="https://github.com/user-attachments/assets/a41c7504-e0be-4aa0-a3aa-5c62c31805b4" />

<h3>Flow Diagram</h3>

<img width="822" height="582" src="https://github.com/user-attachments/assets/84f0e86e-a50c-4684-bf5a-64f1cf44a73f" />

<h3>System Assumption / Setup</h3>

- No authorization. UserID is being inputted within API request.
- No database is being used. Business data is stored within in-memory variables.
- List of stocks and their initial prices are defined in the [cmd/main.go](https://github.com/wfpdev/go_mini_exchange/blob/227e99682d9c86ed00db7324e72beb331205307d/cmd/main/main.go#L25-L34) file.

<h3>Race Condition</h3>

Because of using memory variables, the system is prone to race condition. For example, upon multiple create order requests at the same time.
The system has put the necessary mutex to mitigate it.

<h3>WS Broadcast</h3>

Each of client WebSocket connection is [handled in separate goroutine](https://github.com/wfpdev/go_mini_exchange/blob/227e99682d9c86ed00db7324e72beb331205307d/internal/handler/market.go#L105-L113) so that all connected clients can interact with the server concurrently without blocking each other.

<h3>Bottlenecks</h3>

- Too many orders that may end up flooding Kafka topic. One possible solution is a rate limiter.
- Idempotency keys stored using in-memory cache may end up degrading system performance. One possible solution is to move to Redis.
- Ideally business entities should be stored in database with separate Writer and Reader instances so that read-heavy traffic can be controlled.
