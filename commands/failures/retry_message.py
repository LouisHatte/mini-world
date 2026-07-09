from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    message = world["messages"].get(args.message_id)

    if message is None:
        print(f"Message does not exist: {args.message_id}")
        return

    message["retry_count"] = message.get("retry_count", 0) + 1
    message["status"] = "SENT"
    save_world(world)
    print(f"Retried message: {args.message_id}")
    print(f"retry_count: {message['retry_count']}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "retry-message", help="Retry a failed outbox message."
    )

    parser.add_argument("message_id")

    parser.set_defaults(func=run)
