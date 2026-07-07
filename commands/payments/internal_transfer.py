from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def find_account_id(
    world: dict, human_id: str, bank_id: str, currency: str
) -> str | None:
    for account_id in world["humans"][human_id]["bank_accounts"]:
        account = world["accounts"][account_id]

        if (
            account["owner_human_id"] == human_id
            and account["bank_id"] == bank_id
            and account["currency"] == currency
            and account["status"] == "ACTIVE"
        ):
            return account_id

    return None


def run(args: Namespace) -> None:
    world = load_world()

    sender_human_id = args.sender_human_id
    recipient_human_id = args.recipient_human_id
    bank_id = args.bank_id
    currency = args.currency.upper()
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if sender_human_id == recipient_human_id:
        print("Sender and recipient must be different humans.")
        return

    if sender_human_id not in world["humans"]:
        print(f"Sender human does not exist: {sender_human_id}")
        return

    if recipient_human_id not in world["humans"]:
        print(f"Recipient human does not exist: {recipient_human_id}")
        return

    if bank_id not in world["banks"]:
        print(f"Bank does not exist: {bank_id}")
        return

    sender_account_id = find_account_id(world, sender_human_id, bank_id, currency)

    if sender_account_id is None:
        print(f"No active {currency} account for {sender_human_id} at {bank_id}.")
        return

    recipient_account_id = find_account_id(world, recipient_human_id, bank_id, currency)

    if recipient_account_id is None:
        print(f"No active {currency} account for {recipient_human_id} at {bank_id}.")
        return

    sender_account = world["accounts"][sender_account_id]
    recipient_account = world["accounts"][recipient_account_id]

    if sender_account["booked_balance"] < amount:
        print(
            f"Not enough money in {sender_human_id}'s account. "
            f"Available: {sender_account['booked_balance']} {currency}"
        )
        return

    # Bank 1 reduces its debt to Alice.
    sender_account["booked_balance"] -= amount

    # Bank 1 increases its debt to Bob.
    recipient_account["booked_balance"] += amount

    save_world(world)

    print(
        f"Transferred {amount} {currency} from "
        f"{sender_human_id} to {recipient_human_id} inside {bank_id}."
    )
    print(
        f"{sender_account_id} booked_balance: {sender_account['booked_balance']} {currency}"
    )
    print(
        f"{recipient_account_id} booked_balance: {recipient_account['booked_balance']} {currency}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "internal-transfer",
        help="Transfer deposits between two humans inside the same commercial bank.",
    )

    parser.add_argument("sender_human_id", help="Sender human ID, for example: alice")

    parser.add_argument(
        "recipient_human_id", help="Recipient human ID, for example: bob"
    )

    parser.add_argument("bank_id", help="Commercial bank ID, for example: bank1")

    parser.add_argument("currency", help="Currency, for example: EUR or USD")

    parser.add_argument("amount", type=int, help="Amount to transfer.")

    parser.set_defaults(func=run)
