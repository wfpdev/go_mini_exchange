FROM python:3.13-alpine

WORKDIR /app
COPY simulate_transactions.py ./simulate_transactions.py

USER nobody
ENTRYPOINT ["python", "-u", "simulate_transactions.py"]
