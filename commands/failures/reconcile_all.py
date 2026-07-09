from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    world = load_world()
    problems: list[str] = []

    for account_id, account in world["accounts"].items():
        for hold_id in account.get("holds", []):
            if hold_id not in world["holds"]:
                problems.append(f"{account_id} references missing hold {hold_id}")

    for payment_id, payment in world["payment_instructions"].items():
        message_id = payment.get("message_id")

        if message_id is not None and message_id not in world["messages"]:
            problems.append(f"{payment_id} references missing message {message_id}")

        settlement_id = payment.get("settlement_id")

        if settlement_id is not None and settlement_id not in world["settlements"]:
            problems.append(
                f"{payment_id} references missing settlement {settlement_id}"
            )

    if not problems:
        print("Reconciliation passed.")
        return

    print("Reconciliation found problems:")

    for problem in problems:
        print(f"- {problem}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "reconcile-all", help="Search for inconsistent states."
    )

    parser.set_defaults(func=run)
