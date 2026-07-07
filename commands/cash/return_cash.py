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

    bank_cash = bank["cash_vault"].get(currency, 0)

    if bank_cash < amount:
        print(
            f"Not enough physical cash in {bank_id}'s vault. "
            f"Available: {bank_cash} {currency}"
        )
        return

    central_bank_reserves = central_bank["reserve_accounts"].get(bank_id, 0)
    bank_reserve_mirror = bank["reserve_balances"].get(central_bank_id, 0)

    if central_bank_reserves != bank_reserve_mirror:
        print("Reserve mirror mismatch. Run check-world.")
        print(
            f"{central_bank_id}.reserve_accounts[{bank_id}] = "
            f"{central_bank_reserves}"
        )
        print(
            f"{bank_id}.reserve_balances[{central_bank_id}] = " f"{bank_reserve_mirror}"
        )
        return

    # Bank 1 gives physical cash back to the central bank.
    bank["cash_vault"][currency] = bank_cash - amount

    # The central bank receives physical cash in its vault.
    central_bank["cash_vault"] += amount

    # The central bank credits Bank 1's reserve account.
    central_bank["reserve_accounts"][bank_id] = central_bank_reserves + amount

    # Bank 1's reserve mirror increases.
    bank["reserve_balances"][central_bank_id] = bank_reserve_mirror + amount

    save_world(world)

    print(f"Returned {amount} {currency} cash from {bank_id} to {central_bank_id}.")
    print(
        f"{bank_id} cash_vault[{currency}]: {bank['cash_vault'][currency]} {currency}"
    )
    print(f"{central_bank_id} cash_vault: {central_bank['cash_vault']} {currency}")
    print(
        f"{central_bank_id} reserve account for {bank_id}: {central_bank['reserve_accounts'][bank_id]} {currency}"
    )
    print(
        f"{bank_id} reserves at {central_bank_id}: {bank['reserve_balances'][central_bank_id]} {currency}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "return-cash",
        help="Return physical cash from a commercial bank to a central bank in exchange for reserves.",
    )

    parser.add_argument("central_bank_id", help="Central bank ID, for example: ecb")

    parser.add_argument("bank_id", help="Commercial bank ID, for example: bank1")

    parser.add_argument("amount", type=int, help="Amount of physical cash to return.")

    parser.set_defaults(func=run)
