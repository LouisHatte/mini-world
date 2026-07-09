from argparse import Namespace, _SubParsersAction

from commands.common import append_ledger_entry, require_active_account
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    bank_id = args.bank_id
    payment_id = args.payment_id

    if bank_id not in world["banks"]:
        print(f"Bank does not exist: {bank_id}")
        return

    if payment_id not in world["payment_instructions"]:
        print(f"Payment does not exist: {payment_id}")
        return

    payment = world["payment_instructions"][payment_id]

    if payment["recipient_bank_id"] != bank_id:
        print(f"Payment is not addressed to {bank_id}.")
        return

    account_id, account = require_active_account(
        world, payment["recipient_human_id"], bank_id, payment["currency"]
    )

    if account_id is None or account is None:
        print(f"No active recipient account at {bank_id}.")
        return

    account["booked_balance"] += payment["amount"]
    payment["recipient_account_id"] = account_id
    payment["status"] = "CREDITED"
    append_ledger_entry(
        world,
        bank_id,
        "SEPA incoming credit",
        payment["amount"],
        payment["currency"],
        account_id,
        payment_id,
    )

    save_world(world)

    print(f"Credited incoming SEPA payment: {payment_id}")
    print(
        f"{account_id} booked_balance: "
        f"{account['booked_balance']} {payment['currency']}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "bank-credit-incoming",
        help="Credit a recipient account for an incoming payment.",
    )

    parser.add_argument("bank_id")
    parser.add_argument("payment_id")

    parser.set_defaults(func=run)
