from argparse import Namespace, _SubParsersAction

from commands.common import assert_account_can_debit, require_active_account
from commands.fx.helpers import market_key
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    parts = args.parts

    if len(parts) == 3:
        bank_id = None
        from_currency = parts[0].upper()
        to_currency = parts[1].upper()
        amount = int(parts[2])
    elif len(parts) == 4:
        bank_id = parts[0]
        from_currency = parts[1].upper()
        to_currency = parts[2].upper()
        amount = int(parts[3])
    else:
        print(
            "Usage: fx-convert bank1 EUR USD 100 "
            "OR fx-convert alice bank1 EUR USD 100"
        )
        return

    key = market_key(from_currency, to_currency)
    rate = world["fx_markets"].get(key, {}).get("rate")

    if rate is None:
        print(f"FX market does not exist: {key}")
        return

    converted = round(amount * rate)

    if args.owner_id in world["humans"]:
        if bank_id is None:
            print("Human FX conversion requires bank_id.")
            return

        source_id, source = require_active_account(
            world, args.owner_id, bank_id, from_currency
        )
        target_id, target = require_active_account(
            world, args.owner_id, bank_id, to_currency
        )

        if source_id is None or source is None or target_id is None or target is None:
            print("Source and target accounts must both exist.")
            return

        error = assert_account_can_debit(world, source, amount)

        if error is not None:
            print(error)
            return

        source["booked_balance"] -= amount
        target["booked_balance"] += converted
        save_world(world)
        print(
            f"Converted {args.owner_id}: "
            f"{amount} {from_currency} -> {converted} {to_currency}"
        )
        return

    if args.owner_id in world["banks"]:
        bank = world["banks"][args.owner_id]
        bank.setdefault("fx_inventory", {})
        bank["fx_inventory"][from_currency] = (
            bank["fx_inventory"].get(from_currency, 0) - amount
        )
        bank["fx_inventory"][to_currency] = (
            bank["fx_inventory"].get(to_currency, 0) + converted
        )
        save_world(world)
        print(
            f"Converted {args.owner_id}: "
            f"{amount} {from_currency} -> {converted} {to_currency}"
        )
        return

    print(f"Owner is neither human nor bank: {args.owner_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "fx-convert", help="Convert funds for a human account or bank inventory."
    )

    parser.add_argument("owner_id")
    parser.add_argument("parts", nargs="+")

    parser.set_defaults(func=run)
