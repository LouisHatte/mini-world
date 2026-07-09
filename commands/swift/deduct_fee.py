from argparse import Namespace, _SubParsersAction

from commands.common import append_ledger_entry
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    bank_id = args.bank_id
    payment_id = args.payment_id
    amount = args.amount

    if payment_id not in world["payment_instructions"]:
        print(f"Payment does not exist: {payment_id}")
        return

    payment = world["payment_instructions"][payment_id]
    payment["fees"] = payment.get("fees", 0) + amount
    payment["amount"] -= amount
    append_ledger_entry(
        world,
        bank_id,
        "Intermediary fee",
        amount,
        payment["currency"],
        None,
        payment_id,
    )

    save_world(world)

    print(f"Deducted fee from {payment_id}: {amount} {payment['currency']}")
    print(f"Remaining transfer amount: {payment['amount']} {payment['currency']}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "deduct-fee", help="Deduct an intermediary fee from a payment."
    )

    parser.add_argument("bank_id")
    parser.add_argument("payment_id")
    parser.add_argument("amount", type=int)

    parser.set_defaults(func=run)
