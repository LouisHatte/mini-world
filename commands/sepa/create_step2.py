from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    step2_id = args.step2_id
    currency = args.currency.upper()

    if step2_id in world["step2_systems"]:
        print(f"STEP2 system already exists: {step2_id}")
        return

    world["step2_systems"][step2_id] = {
        "id": step2_id,
        "currency": currency,
        "messages": [],
        "payments": [],
    }

    save_world(world)

    print(f"Created STEP2 system: {step2_id} ({currency})")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "create-step2", help="Create SEPA STEP2 CSM infrastructure."
    )

    parser.add_argument("step2_id")
    parser.add_argument("currency")

    parser.set_defaults(func=run)
