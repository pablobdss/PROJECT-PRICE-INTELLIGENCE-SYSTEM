from fastapi import FastAPI
from pathlib import Path
import json

app = FastAPI()

# pega a raiz do projeto (projeto/)
BASE_DIR = Path(__file__).resolve().parents[3]  
# infrastructure -> app -> python-brain -> projeto

@app.get("/contracts")
def get_price_contract():
    file_path = BASE_DIR / "contracts" / "price_events.json"
    return json.loads(file_path.read_text(encoding="utf-8"))
