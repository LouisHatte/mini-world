from __future__ import annotations

from typing import Any


def next_id(collection: dict[str, Any], prefix: str) -> str:
    return f"{prefix}_{len(collection) + 1:06d}"


def account_available_balance(account: dict[str, Any], holds: dict[str, Any]) -> int:
    held = 0

    for hold_id in account.get("holds", []):
        hold = holds.get(hold_id)

        if hold is not None and hold.get("status") == "HELD":
            held += hold["amount"]

    return account["booked_balance"] - held


def find_account_id(
    world: dict[str, Any], human_id: str, bank_id: str, currency: str
) -> str | None:
    if human_id not in world["humans"]:
        return None

    for account_id in world["humans"][human_id].get("bank_accounts", []):
        account = world["accounts"].get(account_id)

        if (
            account is not None
            and account["owner_human_id"] == human_id
            and account["bank_id"] == bank_id
            and account["currency"] == currency
            and account["status"] == "ACTIVE"
        ):
            return account_id

    return None


def require_active_account(
    world: dict[str, Any], human_id: str, bank_id: str, currency: str
) -> tuple[str | None, dict[str, Any] | None]:
    account_id = find_account_id(world, human_id, bank_id, currency)

    if account_id is None:
        return None, None

    return account_id, world["accounts"][account_id]


def assert_account_can_debit(
    world: dict[str, Any], account: dict[str, Any], amount: int
) -> str | None:
    if amount <= 0:
        return "Amount must be greater than 0."

    if account["status"] != "ACTIVE":
        return f"Account is not active: {account['id']} ({account['status']})"

    available = account_available_balance(account, world["holds"])

    if available < amount:
        return (
            f"Not enough available balance in {account['id']}. "
            f"Available: {available} {account['currency']}"
        )

    return None


def append_ledger_entry(
    world: dict[str, Any],
    entity_id: str,
    description: str,
    amount: int,
    currency: str,
    account_id: str | None = None,
    reference_id: str | None = None,
) -> None:
    world["ledger_entries"].append(
        {
            "id": next_id({str(i): i for i, _ in enumerate(world["ledger_entries"])}, "ledger"),
            "entity_id": entity_id,
            "account_id": account_id,
            "reference_id": reference_id,
            "description": description,
            "amount": amount,
            "currency": currency,
        }
    )
