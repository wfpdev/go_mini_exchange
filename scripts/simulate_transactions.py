#!/usr/bin/env python3

import json
import os
import random
import time
import urllib.error
import urllib.request

API_URL = os.getenv("API_URL", "http://app:8080/create-order")
INTERVAL_SECONDS = float(os.getenv("INTERVAL_SECONDS", "1"))
START_USER_ID = int(os.getenv("START_USER_ID", "1"))

STOCKS = {
    "BBCA": 10000,
    "TLKM": 3000,
    "BBRI": 4800,
    "BMRI": 6200,
    "ASII": 5100,
    "ICBP": 11500,
    "UNVR": 2400,
    "GOTO": 80,
    "BBNI": 4700,
    "TEST": 100,
}


def send_order(stock_code, side, price, quantity, user_id):
    body = json.dumps({
        "user_id": user_id,
        "stock_code": stock_code,
        "side": side,
        "price": price,
        "quantity": quantity,
    }).encode("utf-8")

    request = urllib.request.Request(
        API_URL,
        data=body,
        headers={"Content-Type": "application/json"},
        method="POST",
    )

    try:
        with urllib.request.urlopen(request, timeout=5) as response:
            return json.loads(response.read().decode("utf-8"))
    except urllib.error.HTTPError as error:
        print(f"order rejected: HTTP {error.code} {error.read().decode('utf-8')}")
    except Exception as error:
        print(f"order request failed: {error}")
    return None


def wait_for_api():
    health_url = API_URL.rsplit("/", 1)[0]
    while True:
        try:
            urllib.request.urlopen(health_url + "/ticker?stock_code=BBCA", timeout=2).read()
            print("API is ready")
            return
        except Exception:
            print("Waiting for the API...")
            time.sleep(2)


def main():
    wait_for_api()
    sequence = 0

    while True:
        stock_code, base_price = random.choice(list(STOCKS.items()))
        variation = random.uniform(-0.002, 0.002)
        execution_price = round(base_price * (1 + variation), 2)
        quantity = random.randint(1, 100)
        buyer_id = START_USER_ID + sequence * 2
        seller_id = buyer_id + 1

        buy_response = send_order(
            stock_code, "BUY", execution_price, quantity, buyer_id
        )
        if buy_response is not None:
            print(
                f"submitted BUY {stock_code} price={execution_price:.2f} "
                f"quantity={quantity} status={buy_response.get('status')}"
            )

        sell_response = send_order(
            stock_code, "SELL", execution_price, quantity, seller_id
        )
        if sell_response is not None:
            print(
                f"submitted SELL {stock_code} price={execution_price:.2f} "
                f"quantity={quantity} status={sell_response.get('status')}"
            )

        sequence += 1
        time.sleep(INTERVAL_SECONDS)


if __name__ == "__main__":
    main()
