from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    world = load_world()
    entries = [
        entry
        for entry in world["ledger_entries"]
        if entry["entity_id"] == args.entity_id
    ]

    if not entries:
        print(f"No ledger entries for {args.entity_id}.")
        return

    for entry in entries:
        print(
            f"{entry['id']} {entry['amount']} {entry['currency']} "
            f"{entry['description']} ref={entry.get('reference_id')}"
        )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "show-ledger", help="Show ledger entries for one entity."
    )

    parser.add_argument("entity_id")

    parser.set_defaults(func=run)
