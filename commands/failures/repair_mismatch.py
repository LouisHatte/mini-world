from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    repair_count = 0

    for account in world["accounts"].values():
        account["holds"] = [
            hold_id for hold_id in account.get("holds", []) if hold_id in world["holds"]
        ]
        repair_count += 1

    save_world(world)

    print(f"Repair complete. Accounts checked: {repair_count}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "repair-mismatch", help="Apply simple repair rules for known mismatches."
    )

    parser.set_defaults(func=run)
