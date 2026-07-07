from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    source_bank_id = args.source_bank_id
    target_bank_id = args.target_bank_id
    currency = args.currency.upper()
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if source_bank_id == target_bank_id:
        print("Source bank and target bank must be different.")
        return

    if source_bank_id not in world["banks"]:
        print(f"Source bank does not exist: {source_bank_id}")
        return

    if target_bank_id not in world["banks"]:
        print(f"Target bank does not exist: {target_bank_id}")
        return

    source_bank = world["banks"][source_bank_id]
    target_bank = world["banks"][target_bank_id]

    source_cash = source_bank["cash_vault"].get(currency, 0)

    if source_cash < amount:
        print(
            f"Not enough physical cash in {source_bank_id}'s vault. "
            f"Available: {source_cash} {currency}"
        )
        return

    source_bank["cash_vault"][currency] = source_cash - amount
    target_bank["cash_vault"][currency] = (
        target_bank["cash_vault"].get(currency, 0) + amount
    )

    save_world(world)

    print(
        f"Moved {amount} {currency} cash from " f"{source_bank_id} to {target_bank_id}."
    )
    print(
        f"{source_bank_id} cash_vault[{currency}]: {source_bank['cash_vault'][currency]} {currency}"
    )
    print(
        f"{target_bank_id} cash_vault[{currency}]: {target_bank['cash_vault'][currency]} {currency}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "move-cash",
        help="Move physical cash from one commercial bank vault to another.",
    )

    parser.add_argument(
        "source_bank_id", help="Bank sending physical cash, for example: bank1"
    )

    parser.add_argument(
        "target_bank_id", help="Bank receiving physical cash, for example: bank2"
    )

    parser.add_argument("currency", help="Currency, for example: EUR or USD")

    parser.add_argument("amount", type=int, help="Amount of physical cash to move.")

    parser.set_defaults(func=run)
