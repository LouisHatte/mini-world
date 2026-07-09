from argparse import Namespace, _SubParsersAction

from commands.common import append_ledger_entry
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    bond_id = args.bond_id

    if bond_id not in world["bonds"]:
        print(f"Bond does not exist: {bond_id}")
        return

    bond = world["bonds"][bond_id]
    append_ledger_entry(
        world,
        args.holder,
        f"Coupon on {bond_id}",
        args.amount,
        bond["currency"],
        None,
        bond_id,
    )
    save_world(world)
    print(
        f"Paid coupon on {bond_id} to {args.holder}: "
        f"{args.amount} {bond['currency']}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("pay-coupon", help="Pay bond interest.")

    parser.add_argument("bond_id")
    parser.add_argument("holder")
    parser.add_argument("amount", type=int)

    parser.set_defaults(func=run)
