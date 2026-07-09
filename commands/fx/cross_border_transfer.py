from argparse import Namespace, _SubParsersAction

from commands.common import assert_account_can_debit, next_id, require_active_account
from commands.fx.helpers import market_key
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    from_currency = args.from_currency.upper()
    to_currency = args.to_currency.upper()
    key = market_key(from_currency, to_currency)
    rate = world["fx_markets"].get(key, {}).get("rate")

    if rate is None:
        print(f"FX market does not exist: {key}")
        return

    source_id, source = require_active_account(
        world, args.sender_human_id, args.sender_bank_id, from_currency
    )
    target_id, target = require_active_account(
        world, args.recipient_human_id, args.recipient_bank_id, to_currency
    )

    if source_id is None or source is None or target_id is None or target is None:
        print("Source and target accounts must both exist.")
        return

    error = assert_account_can_debit(world, source, args.amount)

    if error is not None:
        print(error)
        return

    converted = round(args.amount * rate)
    source["booked_balance"] -= args.amount
    target["booked_balance"] += converted
    payment_id = next_id(world["payment_instructions"], "payment")
    world["payment_instructions"][payment_id] = {
        "id": payment_id,
        "rail": "CROSS_BORDER_FX",
        "status": "CREDITED",
        "from_currency": from_currency,
        "to_currency": to_currency,
        "rate": rate,
        "debit_amount": args.amount,
        "credit_amount": converted,
    }
    save_world(world)
    print(f"Cross-border transfer: {payment_id}")
    print(f"Credited {converted} {to_currency} to {target_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "cross-border-transfer", help="Transfer with FX conversion."
    )

    parser.add_argument("sender_human_id")
    parser.add_argument("sender_bank_id")
    parser.add_argument("recipient_human_id")
    parser.add_argument("recipient_bank_id")
    parser.add_argument("from_currency")
    parser.add_argument("to_currency")
    parser.add_argument("amount", type=int)

    parser.set_defaults(func=run)
