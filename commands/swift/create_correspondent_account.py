from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    bank_id = args.bank_id
    correspondent_bank_id = args.correspondent_bank_id
    currency = args.currency.upper()
    account_id = f"corr_{bank_id}_{correspondent_bank_id}_{currency.lower()}"

    if bank_id not in world["banks"] or correspondent_bank_id not in world["banks"]:
        print("Both banks must exist.")
        return

    if account_id in world["correspondent_accounts"]:
        print(f"Correspondent account already exists: {account_id}")
        return

    world["correspondent_accounts"][account_id] = {
        "id": account_id,
        "owner_bank_id": bank_id,
        "correspondent_bank_id": correspondent_bank_id,
        "currency": currency,
        "nostro_balance": 0,
        "vostro_balance": 0,
        "status": "ACTIVE",
    }

    save_world(world)

    print(f"Created correspondent account: {account_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "create-correspondent-account", help="Open a nostro/vostro relationship."
    )

    parser.add_argument("bank_id")
    parser.add_argument("correspondent_bank_id")
    parser.add_argument("currency")

    parser.set_defaults(func=run)
