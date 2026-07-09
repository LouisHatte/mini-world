from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    message = world["messages"].get(args.message_id)

    if message is None:
        print(f"Message does not exist: {args.message_id}")
        return

    message["status"] = "NETWORK_FAILED"
    save_world(world)
    print(f"Simulated network failure: {args.message_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "simulate-network-failure", help="Mark a message as not delivered."
    )

    parser.add_argument("message_id")

    parser.set_defaults(func=run)
