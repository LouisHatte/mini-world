import json
from pathlib import Path
from typing import Any

WORLD_FILE = Path("mini_world.json")


def create_empty_world() -> dict[str, Any]:
    return {
        "version": 1,
        "central_banks": {},
        "banks": {},
        "humans": {},
        "accounts": {},
        "customer_loans": {},
        "holds": {},
        "ledger_entries": [],
        "payment_instructions": {},
        "messages": {},
        "settlements": {},
        "step2_systems": {},
        "correspondent_accounts": {},
        "bonds": {},
        "cheques": {},
        "card_authorizations": {},
        "fx_markets": {},
        "snapshots": {},
        "command_history": [],
    }


def ensure_world_shape(world: dict[str, Any]) -> dict[str, Any]:
    """
    Add missing top-level keys to older world files.

    This lets us evolve the JSON schema without breaking existing worlds.
    """
    world.setdefault("version", 1)
    world.setdefault("central_banks", {})
    world.setdefault("banks", {})
    world.setdefault("humans", {})
    world.setdefault("accounts", {})
    world.setdefault("customer_loans", {})
    world.setdefault("holds", {})
    world.setdefault("ledger_entries", [])
    world.setdefault("payment_instructions", {})
    world.setdefault("messages", {})
    world.setdefault("settlements", {})
    world.setdefault("step2_systems", {})
    world.setdefault("correspondent_accounts", {})
    world.setdefault("bonds", {})
    world.setdefault("cheques", {})
    world.setdefault("card_authorizations", {})
    world.setdefault("fx_markets", {})
    world.setdefault("snapshots", {})
    world.setdefault("command_history", [])

    return world


def world_exists() -> bool:
    return WORLD_FILE.exists()


def load_world() -> dict[str, Any]:
    if not world_exists():
        raise FileNotFoundError(
            f"World does not exist: {WORLD_FILE}. Run: python3.11 mini.py init"
        )

    with WORLD_FILE.open("r", encoding="utf-8") as file:
        world = json.load(file)

    return ensure_world_shape(world)


def save_world(world: dict[str, Any]) -> None:
    world = ensure_world_shape(world)

    with WORLD_FILE.open("w", encoding="utf-8") as file:
        json.dump(world, file, indent=2, ensure_ascii=False)
        file.write("\n")
