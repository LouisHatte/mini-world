from argparse import Namespace, _SubParsersAction

from commands.common import append_ledger_entry, require_active_account
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    cheque_id = args.cheque_id

    if cheque_id not in world["cheques"]:
        print(f"Cheque does not exist: {cheque_id}")
        return

    cheque = world["cheques"][cheque_id]
    account_id, account = require_active_account(
        world, args.payee_human_id, args.payee_bank_id, cheque["currency"]
    )

    if account_id is None or account is None:
        print("Payee account does not exist.")
        return

    account["booked_balance"] += cheque["amount"]
    cheque["payee_bank_id"] = args.payee_bank_id
    cheque["payee_account_id"] = account_id
    cheque["status"] = "PROVISIONALLY_CREDITED"
    append_ledger_entry(
        world,
        args.payee_bank_id,
        "Provisional cheque credit",
        cheque["amount"],
        cheque["currency"],
        account_id,
        cheque_id,
    )
    save_world(world)
    print(f"Deposited cheque: {cheque_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("deposit-cheque", help="Deposit a cheque at a bank.")

    parser.add_argument("payee_human_id")
    parser.add_argument("payee_bank_id")
    parser.add_argument("cheque_id")

    parser.set_defaults(func=run)
