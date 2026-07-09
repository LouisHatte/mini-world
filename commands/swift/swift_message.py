from argparse import Namespace, _SubParsersAction

from commands.common import next_id
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    message_id = next_id(world["messages"], "msg")
    world["messages"][message_id] = {
        "id": message_id,
        "rail": "SWIFT",
        "sender_bank_id": args.sender_bank_id,
        "recipient_bank_id": args.recipient_bank_id,
        "type": args.message_type,
        "status": "SENT",
    }
    save_world(world)
    print(f"Created SWIFT message: {message_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "swift-message", help="Create a SWIFT payment instruction message."
    )

    parser.add_argument("sender_bank_id")
    parser.add_argument("recipient_bank_id")
    parser.add_argument("message_type")

    parser.set_defaults(func=run)
