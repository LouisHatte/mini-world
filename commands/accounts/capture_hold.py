from argparse import Namespace, _SubParsersAction

from commands.common import append_ledger_entry
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    hold_id = args.hold_id

    if hold_id not in world["holds"]:
        print(f"Hold does not exist: {hold_id}")
        return

    hold = world["holds"][hold_id]

    if hold["status"] != "HELD":
        print(f"Hold is not capturable: {hold_id} ({hold['status']})")
        return

    account = world["accounts"].get(hold["account_id"])

    if account is None:
        print(f"Account does not exist: {hold['account_id']}")
        return

    if account["status"] != "ACTIVE":
        print(f"Account is not active: {account['id']} ({account['status']})")
        return

    account["booked_balance"] -= hold["amount"]
    hold["status"] = "CAPTURED"

    append_ledger_entry(
        world,
        hold["bank_id"],
        "Captured account hold",
        -hold["amount"],
        hold["currency"],
        account["id"],
        hold_id,
    )

    save_world(world)

    print(f"Captured hold: {hold_id}")
    print(
        f"{account['id']} booked_balance: "
        f"{account['booked_balance']} {hold['currency']}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "capture-hold", help="Book a debit for a held account reservation."
    )

    parser.add_argument("hold_id")

    parser.set_defaults(func=run)
