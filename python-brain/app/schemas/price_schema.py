from pydantic import BaseModel

class PriceInput(BaseModel):
    product_id: str
    price: float
    timestamp: str
