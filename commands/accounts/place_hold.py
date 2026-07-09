from argparse import Namespace, _SubParsersAction

from commands.common import (
    account_available_balance,
    assert_account_can_debit,
    next_id,
    require_active_account,
)
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    human_id = args.human_id
    bank_id = args.bank_id
    currency = args.currency.upper()
    amount = args.amount

    if human_id not in world["humans"]:
        print(f"Human does not exist: {human_id}")
        return

    if bank_id not in world["banks"]:
        print(f"Bank does not exist: {bank_id}")
        return

    account_id, account = require_active_account(world, human_id, bank_id, currency)

    if account_id is None or account is None:
        print(f"No active {currency} account for {human_id} at {bank_id}.")
        return

    error = assert_account_can_debit(world, account, amount)

    if error is not None:
        print(error)
        return

    hold_id = next_id(world["holds"], "hold")
    hold = {
        "id": hold_id,
        "account_id": account_id,
        "human_id": human_id,
        "bank_id": bank_id,
        "currency": currency,
        "amount": amount,
        "status": "HELD",
    }

    world["holds"][hold_id] = hold
    account.setdefault("holds", []).append(hold_id)

    save_world(world)

    print(f"Placed hold: {hold_id}")
    print(f"{account_id} booked_balance: {account['booked_balance']} {currency}")
    print(
        f"{account_id} available_balance: "
        f"{account_available_balance(account, world['holds'])} {currency}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "place-hold", help="Reserve part of an account balance."
    )

    parser.add_argument("human_id")
    parser.add_argument("bank_id")
    parser.add_argument("currency")
    parser.add_argument("amount", type=int)

    parser.set_defaults(func=run)
