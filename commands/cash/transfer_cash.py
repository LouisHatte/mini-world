from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    source_human_id = args.source_human_id
    target_human_id = args.target_human_id
    currency = args.currency.upper()
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if source_human_id == target_human_id:
        print("Source human and target human must be different.")
        return

    if source_human_id not in world["humans"]:
        print(f"Source human does not exist: {source_human_id}")
        return

    if target_human_id not in world["humans"]:
        print(f"Target human does not exist: {target_human_id}")
        return

    source_human = world["humans"][source_human_id]
    target_human = world["humans"][target_human_id]

    source_cash = source_human["cash_wallet"].get(currency, 0)

    if source_cash < amount:
        print(
            f"Not enough cash in {source_human_id}'s wallet. "
            f"Available: {source_cash} {currency}"
        )
        return

    source_human["cash_wallet"][currency] = source_cash - amount
    target_human["cash_wallet"][currency] = (
        target_human["cash_wallet"].get(currency, 0) + amount
    )

    save_world(world)

    print(
        f"Transferred {amount} {currency} cash from "
        f"{source_human_id} to {target_human_id}."
    )
    print(
        f"{source_human_id} cash_wallet[{currency}]: {source_human['cash_wallet'][currency]} {currency}"
    )
    print(
        f"{target_human_id} cash_wallet[{currency}]: {target_human['cash_wallet'][currency]} {currency}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "transfer-cash", help="Transfer physical cash from one human to another."
    )

    parser.add_argument("source_human_id", help="Human paying cash, for example: alice")

    parser.add_argument(
        "target_human_id", help="Human receiving cash, for example: bob"
    )

    parser.add_argument("currency", help="Currency, for example: EUR or USD")

    parser.add_argument("amount", type=int, help="Amount of physical cash to transfer.")

    parser.set_defaults(func=run)
