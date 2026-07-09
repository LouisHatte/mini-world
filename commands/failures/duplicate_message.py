from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    message = world["messages"].get(args.message_id)

    if message is None:
        print(f"Message does not exist: {args.message_id}")
        return

    if message.get("duplicate_seen"):
        print(f"Duplicate ignored idempotently: {args.message_id}")
        return

    message["duplicate_seen"] = True
    save_world(world)
    print(f"Marked duplicate message as seen: {args.message_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "duplicate-message", help="Test message idempotency."
    )

    parser.add_argument("message_id")

    parser.set_defaults(func=run)
