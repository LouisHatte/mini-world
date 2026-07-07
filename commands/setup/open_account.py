from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def build_account_id(bank_id: str, human_id: str, currency: str) -> str:
    return f"acc_{bank_id}_{human_id}_{currency.lower()}"


def run(args: Namespace) -> None:
    world = load_world()

    human_id = args.human_id
    bank_id = args.bank_id
    currency = args.currency.upper()

    if human_id not in world["humans"]:
        print(f"Human does not exist: {human_id}")
        return

    if bank_id not in world["banks"]:
        print(f"Bank does not exist: {bank_id}")
        return

    account_id = build_account_id(bank_id, human_id, currency)

    if account_id in world["accounts"]:
        print(f"Account already exists: {account_id}")
        return

    account = {
        "id": account_id,
        "owner_human_id": human_id,
        "bank_id": bank_id,
        "currency": currency,
        "booked_balance": 0,
        "holds": [],
        "status": "ACTIVE",
    }

    world["accounts"][account_id] = account

    world["humans"][human_id]["bank_accounts"].append(account_id)
    world["banks"][bank_id]["customer_accounts"].append(account_id)

    save_world(world)

    print(f"Opened account: {account_id}")
    print(f"Owner: {human_id}")
    print(f"Bank: {bank_id}")
    print(f"Currency: {currency}")
    print("Booked balance: 0")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "open-account", help="Open a bank account for a human at a commercial bank."
    )

    parser.add_argument("human_id", help="Human ID, for example: alice")

    parser.add_argument("bank_id", help="Bank ID, for example: bank1")

    parser.add_argument("currency", help="Account currency, for example: EUR or USD")

    parser.set_defaults(func=run)
