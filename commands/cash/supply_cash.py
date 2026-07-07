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
    currency = central_bank["currency"]

    central_bank_cash_vault = central_bank["cash_vault"]
    bank_reserves_at_central_bank = bank["reserve_balances"].get(central_bank_id, 0)

    if central_bank_cash_vault < amount:
        print(
            f"Not enough physical cash in {central_bank_id} vault. "
            f"Available: {central_bank_cash_vault} {currency}"
        )
        return

    if bank_reserves_at_central_bank < amount:
        print(
            f"Not enough reserves for {bank_id} at {central_bank_id}. "
            f"Available: {bank_reserves_at_central_bank} {currency}"
        )
        return

    # Central bank side:
    # Physical cash leaves the central bank vault.
    central_bank["cash_vault"] -= amount

    # Bank 1 pays for the cash by reducing its reserve balance at the central bank.
    central_bank["reserve_accounts"][bank_id] -= amount

    # Commercial bank side:
    # Bank 1 receives physical cash in its vault.
    bank["cash_vault"][currency] = bank["cash_vault"].get(currency, 0) + amount

    # Bank 1's reserve mirror decreases.
    bank["reserve_balances"][central_bank_id] -= amount

    save_world(world)

    print(f"Supplied {amount} {currency} cash from {central_bank_id} to {bank_id}.")
    print(f"{central_bank_id} cash_vault: {central_bank['cash_vault']} {currency}")
    print(
        f"{central_bank_id} reserve account for {bank_id}: {central_bank['reserve_accounts'][bank_id]} {currency}"
    )
    print(
        f"{bank_id} cash_vault[{currency}]: {bank['cash_vault'][currency]} {currency}"
    )
    print(
        f"{bank_id} reserves at {central_bank_id}: {bank['reserve_balances'][central_bank_id]} {currency}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "supply-cash",
        help="Move physical cash from a central bank vault to a commercial bank vault.",
    )

    parser.add_argument("central_bank_id", help="Central bank ID, for example: ecb")

    parser.add_argument("bank_id", help="Commercial bank ID, for example: bank1")

    parser.add_argument("amount", type=int, help="Amount of physical cash to supply.")

    parser.set_defaults(func=run)
