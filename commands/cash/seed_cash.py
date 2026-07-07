from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    central_bank_id = args.central_bank_id
    human_id = args.human_id
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if central_bank_id not in world["central_banks"]:
        print(f"Central bank does not exist: {central_bank_id}")
        return

    if human_id not in world["humans"]:
        print(f"Human does not exist: {human_id}")
        return

    central_bank = world["central_banks"][central_bank_id]
    human = world["humans"][human_id]
    currency = central_bank["currency"]

    if central_bank["cash_vault"] < amount:
        print(
            f"Not enough cash in {central_bank_id} vault. "
            f"Available: {central_bank['cash_vault']} {currency}"
        )
        return

    central_bank["cash_vault"] -= amount
    human["cash_wallet"][currency] = human["cash_wallet"].get(currency, 0) + amount

    save_world(world)

    print(f"Seeded {amount} {currency} cash from {central_bank_id} to {human_id}.")
    print(f"{central_bank_id} cash_vault: {central_bank['cash_vault']} {currency}")
    print(
        f"{human_id} cash_wallet[{currency}]: {human['cash_wallet'][currency]} {currency}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "seed-cash",
        help="Seed already-issued physical cash from a central bank vault to a human.",
    )

    parser.add_argument("central_bank_id", help="Central bank ID, for example: ecb")

    parser.add_argument("human_id", help="Human ID, for example: alice")

    parser.add_argument("amount", type=int, help="Amount of physical cash to seed.")

    parser.set_defaults(func=run)
