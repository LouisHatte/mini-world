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

    human_id = args.human_id
    bank_id = args.bank_id
    currency = args.currency.upper()
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if human_id not in world["humans"]:
        print(f"Human does not exist: {human_id}")
        return

    if bank_id not in world["banks"]:
        print(f"Bank does not exist: {bank_id}")
        return

    human = world["humans"][human_id]
    bank = world["banks"][bank_id]

    human_cash = human["cash_wallet"].get(currency, 0)

    if human_cash < amount:
        print(
            f"Not enough cash in {human_id}'s wallet. "
            f"Available: {human_cash} {currency}"
        )
        return

    account_id = find_account_id(world, human_id, bank_id, currency)

    if account_id is None:
        print(f"No active {currency} account for {human_id} at {bank_id}.")
        return

    account = world["accounts"][account_id]

    # Human gives physical cash to the bank.
    human["cash_wallet"][currency] = human_cash - amount

    # Bank receives physical cash in its vault.
    bank["cash_vault"][currency] = bank["cash_vault"].get(currency, 0) + amount

    # Bank credits Alice's deposit account.
    # This means the bank now owes Alice more money.
    account["booked_balance"] += amount

    save_world(world)

    print(f"{human_id} deposited {amount} {currency} cash at {bank_id}.")
    print(
        f"{human_id} cash_wallet[{currency}]: {human['cash_wallet'][currency]} {currency}"
    )
    print(
        f"{bank_id} cash_vault[{currency}]: {bank['cash_vault'][currency]} {currency}"
    )
    print(f"{account_id} booked_balance: {account['booked_balance']} {currency}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "deposit-cash",
        help="Deposit physical cash from a human wallet into a bank account.",
    )

    parser.add_argument("human_id", help="Human ID, for example: alice")

    parser.add_argument("bank_id", help="Bank ID, for example: bank1")

    parser.add_argument("currency", help="Currency, for example: EUR or USD")

    parser.add_argument("amount", type=int, help="Amount of physical cash to deposit.")

    parser.set_defaults(func=run)
