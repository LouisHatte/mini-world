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

    if sender_bank == recipient_bank:
        print("Use internal-transfer for same-bank payments.")
        return

    if sender not in world["humans"] or recipient not in world["humans"]:
        print("Sender or recipient human does not exist.")
        return

    if sender_bank not in world["banks"] or recipient_bank not in world["banks"]:
        print("Sender or recipient bank does not exist.")
        return

    account_id, account = require_active_account(world, sender, sender_bank, currency)

    if account_id is None or account is None:
        print(f"No active {currency} account for {sender} at {sender_bank}.")
        return

    error = assert_account_can_debit(world, account, amount)

    if error is not None:
        print(error)
        return

    step2_id = next(iter(world["step2_systems"]), "step2")

    if step2_id not in world["step2_systems"]:
        world["step2_systems"][step2_id] = {
            "id": step2_id,
            "currency": currency,
            "messages": [],
            "payments": [],
        }

    payment_id = next_id(world["payment_instructions"], "payment")
    message_id = next_id(world["messages"], "msg")

    account["booked_balance"] -= amount
    world["payment_instructions"][payment_id] = {
        "id": payment_id,
        "rail": "SEPA",
        "status": "INITIATED",
        "sender_human_id": sender,
        "sender_bank_id": sender_bank,
        "sender_account_id": account_id,
        "recipient_human_id": recipient,
        "recipient_bank_id": recipient_bank,
        "recipient_account_id": None,
        "currency": currency,
        "amount": amount,
        "message_id": message_id,
        "step2_id": step2_id,
    }
    world["messages"][message_id] = {
        "id": message_id,
        "payment_id": payment_id,
        "rail": "SEPA",
        "sender_bank_id": sender_bank,
        "recipient_bank_id": recipient_bank,
        "status": "SENT",
    }
    world["step2_systems"][step2_id]["messages"].append(message_id)
    world["step2_systems"][step2_id]["payments"].append(payment_id)
    append_ledger_entry(
        world,
        sender_bank,
        "SEPA outgoing debit",
        -amount,
        currency,
        account_id,
        payment_id,
    )

    save_world(world)

    print(f"Created SEPA payment: {payment_id}")
    print(f"Message: {message_id}")
    print(f"{account_id} booked_balance: {account['booked_balance']} {currency}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "sepa-transfer", help="Create a full SEPA payment instruction."
    )

    parser.add_argument("sender_human_id")
    parser.add_argument("sender_bank_id")
    parser.add_argument("recipient_human_id")
    parser.add_argument("recipient_bank_id")
    parser.add_argument("currency")
    parser.add_argument("amount", type=int)

    parser.set_defaults(func=run)
