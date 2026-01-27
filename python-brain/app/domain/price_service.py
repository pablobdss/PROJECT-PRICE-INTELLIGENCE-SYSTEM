from app.domain.price import Price

def process_price(data):
    price = Price(data.product_id, data.price, data.timestamp)
    print(f"[PRICE RECEIVED] {price.product_id} -> {price.price}")
    return price
