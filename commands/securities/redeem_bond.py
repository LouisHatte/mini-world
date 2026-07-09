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
    bond["holders"][args.holder] = max(
        0, bond["holders"].get(args.holder, 0) - args.amount
    )

    if sum(bond["holders"].values()) == 0:
        bond["status"] = "REDEEMED"

    append_ledger_entry(
        world,
        args.holder,
        f"Redemption of {bond_id}",
        args.amount,
        bond["currency"],
        None,
        bond_id,
    )
    save_world(world)
    print(f"Redeemed {args.amount} {bond['currency']} of {bond_id} for {args.holder}.")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("redeem-bond", help="Repay bond principal.")

    parser.add_argument("bond_id")
    parser.add_argument("holder")
    parser.add_argument("amount", type=int)

    parser.set_defaults(func=run)
