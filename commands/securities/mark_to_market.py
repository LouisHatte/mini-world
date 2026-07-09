from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    bond_id = args.bond_id

    if bond_id not in world["bonds"]:
        print(f"Bond does not exist: {bond_id}")
        return

    world["bonds"][bond_id]["market_value"] = args.market_value
    save_world(world)
    print(
        f"{bond_id} market_value: "
        f"{args.market_value} {world['bonds'][bond_id]['currency']}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "mark-to-market", help="Update market value of a security."
    )

    parser.add_argument("bond_id")
    parser.add_argument("market_value", type=int)

    parser.set_defaults(func=run)
