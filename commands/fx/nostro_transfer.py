from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    currency = args.currency.upper()
    account = None

    for candidate in world["correspondent_accounts"].values():
        if (
            candidate["owner_bank_id"] == args.bank_id
            and candidate["correspondent_bank_id"] == args.correspondent_bank_id
            and candidate["currency"] == currency
        ):
            account = candidate
            break

    if account is None:
        print("Correspondent account does not exist.")
        return

    account["nostro_balance"] += args.amount
    save_world(world)
    print(f"Nostro transfer: {account['id']} {account['nostro_balance']} {currency}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "nostro-transfer", help="Move funds through a nostro account."
    )

    parser.add_argument("bank_id")
    parser.add_argument("correspondent_bank_id")
    parser.add_argument("currency")
    parser.add_argument("amount", type=int)

    parser.set_defaults(func=run)
