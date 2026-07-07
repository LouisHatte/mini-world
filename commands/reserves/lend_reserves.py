from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    central_bank_id = args.central_bank_id
    bank_id = args.bank_id
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if central_bank_id not in world["central_banks"]:
        print(f"Central bank does not exist: {central_bank_id}")
        return

    if bank_id not in world["banks"]:
        print(f"Bank does not exist: {bank_id}")
        return

    central_bank = world["central_banks"][central_bank_id]
    bank = world["banks"][bank_id]

    central_bank["reserve_accounts"][bank_id] = (
        central_bank["reserve_accounts"].get(bank_id, 0) + amount
    )

    central_bank["loans_to_banks"][bank_id] = (
        central_bank["loans_to_banks"].get(bank_id, 0) + amount
    )

    bank["reserve_balances"][central_bank_id] = (
        bank["reserve_balances"].get(central_bank_id, 0) + amount
    )

    bank["loans_from_central_banks"][central_bank_id] = (
        bank["loans_from_central_banks"].get(central_bank_id, 0) + amount
    )

    save_world(world)

    currency = central_bank["currency"]

    print(f"Lent {amount} {currency} reserves from {central_bank_id} to {bank_id}.")
    print(
        f"{central_bank_id} reserve account for {bank_id}: {central_bank['reserve_accounts'][bank_id]} {currency}"
    )
    print(
        f"{central_bank_id} loan to {bank_id}: {central_bank['loans_to_banks'][bank_id]} {currency}"
    )
    print(
        f"{bank_id} reserves at {central_bank_id}: {bank['reserve_balances'][central_bank_id]} {currency}"
    )
    print(
        f"{bank_id} loan from {central_bank_id}: {bank['loans_from_central_banks'][central_bank_id]} {currency}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "lend-reserves",
        help="Create central bank reserves by lending them to a commercial bank.",
    )

    parser.add_argument("central_bank_id", help="Central bank ID, for example: ecb")

    parser.add_argument("bank_id", help="Commercial bank ID, for example: bank1")

    parser.add_argument("amount", type=int, help="Amount of reserves to lend.")

    parser.set_defaults(func=run)
