def find_central_bank(world: dict, currency: str) -> str | None:
    for central_bank_id, central_bank in world["central_banks"].items():
        if central_bank["currency"] == currency:
            return central_bank_id

    return None
