"""
Serviço de Plataforma/Central (P/C)
Responsável por verificar disponibilidade de slots (helipads/docks)
"""

import os
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from enum import Enum
from typing import Dict, List
import uvicorn

app = FastAPI(title="Platform/Central Service", version="1.0.0")

FREE = True
IN_USE = False

DOCKS_QTT = int(os.getenv("DOCKS_QTT", "5"))
HELIPADS_QTT = int(os.getenv("HELIPADS_QTT", "3"))


class SlotType(str, Enum):
    """Tipos de slots disponíveis"""
    DOCK = "dock"
    HELIPAD = "helipad"


class SlotState(str, Enum):
    """Estados possíveis de um slot"""
    FREE = "free"
    IN_USE = "in_use"


# Estado em memória - slots disponíveis/ocupados
# Estrutura: {"docks": [True, True, ...], "helipads": [True, True, ...]}
slots: Dict[str, List[bool]] = {"docks": [], "helipads": []}

for i in range(DOCKS_QTT):
    slots["docks"].append(FREE)

for i in range(HELIPADS_QTT):
    slots["helipads"].append(FREE)


class SlotRequest(BaseModel):
    slot_number: int
    slot_type: SlotType  #


class SlotResponse(BaseModel):
    state: SlotState 


@app.get("/")
async def root():
    """Endpoint raiz para verificação de saúde do serviço"""
    return {
        "service": "Platform/Central Service",
        "status": "operational"
    }


@app.post("/slots", response_model=SlotResponse)
async def check_slot(request: SlotRequest):
    """
    Endpoint para verificar o estado de um slot específico
    Chamado pela Torre (T) para verificar disponibilidade
    """
    slot_type = request.slot_type.value  
    
    if slot_type not in slots:
        raise HTTPException(status_code=500, detail="Internal error: slot type not initialized")
    
    if request.slot_number < 0 or request.slot_number >= len(slots[slot_type]):
        raise HTTPException(
            status_code=404,
            detail=f"Slot {request.slot_number} of type {slot_type} not found"
        )
    
    if slots[slot_type][request.slot_number] == FREE:
        state = SlotState.FREE
    else:
        state = SlotState.IN_USE
    
    return SlotResponse(state=state)


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)