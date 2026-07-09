from argparse import Namespace

from world import load_world, save_world


def adjust_correspondent(args: Namespace, field: str, verb: str) -> None:
    world = load_world()
    currency = args.currency.upper()
    account = None

    for candidate in world["correspondent_accounts"].values():
        if (
            candidate["owner_bank_id"] == args.owner_bank_id
            and candidate["correspondent_bank_id"] == args.correspondent_bank_id
            and candidate["currency"] == currency
        ):
            account = candidate
            break

    if account is None:
        print("Correspondent account does not exist.")
        return

    account[field] += args.amount
    save_world(world)
    print(f"Correspondent {verb}: {account['id']}")
    print(f"{field}: {account[field]} {currency}")
