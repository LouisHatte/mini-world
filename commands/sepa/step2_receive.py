from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    message_id = args.message_id

    if message_id not in world["messages"]:
        print(f"Message does not exist: {message_id}")
        return

    message = world["messages"][message_id]
    message["status"] = "RECEIVED_BY_STEP2"
    world["payment_instructions"][message["payment_id"]]["status"] = "STEP2_RECEIVED"

    save_world(world)

    print(f"STEP2 received message: {message_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "step2-receive", help="Mark a SEPA message as received by STEP2."
    )

    parser.add_argument("message_id")

    parser.set_defaults(func=run)
