from argparse import Namespace, _SubParsersAction

from commands.common import (
    append_ledger_entry,
    assert_account_can_debit,
    next_id,
    require_active_account,
)
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    sender = args.sender_human_id
    sender_bank = args.sender_bank_id
    recipient = args.recipient_human_id
    recipient_bank = args.recipient_bank_id
    currency = args.currency.upper()
    amount = args.amount
    account_id, account = require_active_account(world, sender, sender_bank, currency)

    if account_id is None or account is None:
        print(f"No active {currency} account for {sender} at {sender_bank}.")
        return

    error = assert_account_can_debit(world, account, amount)

    if error is not None:
        print(error)
        return

    payment_id = next_id(world["payment_instructions"], "payment")
    message_id = next_id(world["messages"], "msg")
    account["booked_balance"] -= amount
    world["payment_instructions"][payment_id] = {
        "id": payment_id,
        "rail": "SWIFT",
        "status": "INSTRUCTED",
        "sender_human_id": sender,
        "sender_bank_id": sender_bank,
        "sender_account_id": account_id,
        "recipient_human_id": recipient,
        "recipient_bank_id": recipient_bank,
        "currency": currency,
        "amount": amount,
        "fees": 0,
        "message_ids": [message_id],
    }
    world["messages"][message_id] = {
        "id": message_id,
        "payment_id": payment_id,
        "rail": "SWIFT",
        "type": "MT103",
        "sender_bank_id": sender_bank,
        "recipient_bank_id": recipient_bank,
        "status": "SENT",
    }
    append_ledger_entry(
        world,
        sender_bank,
        "SWIFT outgoing debit",
        -amount,
        currency,
        account_id,
        payment_id,
    )

    save_world(world)

    print(f"Created SWIFT payment: {payment_id}")
    print(f"Message: {message_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "swift-transfer", help="Create an international SWIFT-style transfer."
    )

    parser.add_argument("sender_human_id")
    parser.add_argument("sender_bank_id")
    parser.add_argument("recipient_human_id")
    parser.add_argument("recipient_bank_id")
    parser.add_argument("currency")
    parser.add_argument("amount", type=int)

    parser.set_defaults(func=run)
