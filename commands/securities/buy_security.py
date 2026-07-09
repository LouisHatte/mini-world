from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    central_bank_id = args.central_bank_id
    bank_id = args.bank_id
    bond_id = args.bond_id
    amount = args.amount

    if central_bank_id not in world["central_banks"] or bank_id not in world["banks"]:
        print("Central bank or bank does not exist.")
        return

    if bond_id not in world["bonds"]:
        print(f"Bond does not exist: {bond_id}")
        return

    central_bank = world["central_banks"][central_bank_id]
    bank = world["banks"][bank_id]
    bond = world["bonds"][bond_id]

    central_bank.setdefault("securities", {})
    central_bank["securities"][bond_id] = (
        central_bank["securities"].get(bond_id, 0) + amount
    )
    bond["holders"][central_bank_id] = bond["holders"].get(central_bank_id, 0) + amount
    central_bank["reserve_accounts"][bank_id] = (
        central_bank["reserve_accounts"].get(bank_id, 0) + amount
    )
    bank["reserve_balances"][central_bank_id] = (
        bank["reserve_balances"].get(central_bank_id, 0) + amount
    )

    save_world(world)

    print(
        f"{central_bank_id} bought {amount} {bond['currency']} "
        f"of {bond_id} from {bank_id}."
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "buy-security", help="Central bank buys a security and creates reserves."
    )

    parser.add_argument("central_bank_id")
    parser.add_argument("bank_id")
    parser.add_argument("bond_id")
    parser.add_argument("amount", type=int)

    parser.set_defaults(func=run)
