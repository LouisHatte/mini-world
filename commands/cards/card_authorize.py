from argparse import Namespace, _SubParsersAction

from commands.common import assert_account_can_debit, next_id
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    currency = args.currency.upper()
    merchant = args.merchant

    account_id = None
    account = None
    bank_id = None

    for candidate_id in world["humans"].get(args.human_id, {}).get("bank_accounts", []):
        candidate = world["accounts"][candidate_id]

        if candidate["currency"] == currency and candidate["status"] == "ACTIVE":
            account_id = candidate_id
            account = candidate
            bank_id = candidate["bank_id"]
            break

    if account_id is None or account is None or bank_id is None:
        print(f"No active {currency} account for {args.human_id}.")
        return

    error = assert_account_can_debit(world, account, args.amount)

    if error is not None:
        print(error)
        return

    auth_id = next_id(world["card_authorizations"], "auth")
    hold_id = next_id(world["holds"], "hold")
    world["holds"][hold_id] = {
        "id": hold_id,
        "account_id": account_id,
        "human_id": args.human_id,
        "bank_id": bank_id,
        "currency": currency,
        "amount": args.amount,
        "status": "HELD",
    }
    account.setdefault("holds", []).append(hold_id)
    world["card_authorizations"][auth_id] = {
        "id": auth_id,
        "hold_id": hold_id,
        "human_id": args.human_id,
        "merchant": merchant,
        "account_id": account_id,
        "currency": currency,
        "amount": args.amount,
        "status": "AUTHORIZED",
    }

    save_world(world)

    print(f"Authorized card payment: {auth_id}")
    print(f"Hold: {hold_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "card-authorize", help="Create a card authorization hold."
    )

    parser.add_argument("human_id")
    parser.add_argument("merchant")
    parser.add_argument("currency")
    parser.add_argument("amount", type=int)

    parser.set_defaults(func=run)
