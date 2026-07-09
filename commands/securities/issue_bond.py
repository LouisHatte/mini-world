from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    bond_id = args.bond_id
    currency = args.currency.upper()

    if bond_id in world["bonds"]:
        print(f"Bond already exists: {bond_id}")
        return

    world["bonds"][bond_id] = {
        "id": bond_id,
        "issuer": args.issuer,
        "currency": currency,
        "face_value": args.face_value,
        "market_value": args.face_value,
        "holders": {},
        "status": "ISSUED",
    }

    save_world(world)

    print(f"Issued bond: {bond_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser("issue-bond", help="Issue a bond.")

    parser.add_argument("issuer")
    parser.add_argument("bond_id")
    parser.add_argument("currency")
    parser.add_argument("face_value", type=int)

    parser.set_defaults(func=run)
