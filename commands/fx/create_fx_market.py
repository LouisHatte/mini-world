from argparse import Namespace, _SubParsersAction

from commands.fx.helpers import market_key
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    from_currency = args.from_currency.upper()
    to_currency = args.to_currency.upper()
    key = market_key(from_currency, to_currency)
    world["fx_markets"][key] = {
        "from_currency": from_currency,
        "to_currency": to_currency,
        "rate": args.rate,
    }
    save_world(world)
    print(f"FX market {key}: {args.rate}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("create-fx-market", help="Define an FX rate.")

    parser.add_argument("from_currency")
    parser.add_argument("to_currency")
    parser.add_argument("rate", type=float)

    parser.set_defaults(func=run)
