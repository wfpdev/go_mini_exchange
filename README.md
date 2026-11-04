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

<img width="471" height="262" alt="CV Elaboration-Page-5 drawio" src="https://github.com/user-attachments/assets/a41c7504-e0be-4aa0-a3aa-5c62c31805b4" />

<h3>Flow Diagram</h3>

<img width="822" height="582" alt="CV Elaboration-Page-6 drawio" src="https://github.com/user-attachments/assets/84f0e86e-a50c-4684-bf5a-64f1cf44a73f" />

