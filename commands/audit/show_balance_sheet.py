from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    world = load_world()
    bank = world["banks"].get(args.bank_id)

    if bank is None:
        print(f"Bank does not exist: {args.bank_id}")
        return

    print(f"Balance sheet for {args.bank_id}")
    print(f"Assets reserves: {bank.get('reserve_balances', {})}")
    print(f"Assets cash_vault: {bank.get('cash_vault', {})}")
    print(f"Assets customer_loans: {bank.get('customer_loans', [])}")

    deposits = {}

    for account_id in bank.get("customer_accounts", []):
        account = world["accounts"][account_id]
        deposits[account["currency"]] = (
            deposits.get(account["currency"], 0) + account["booked_balance"]
        )

    print(f"Liabilities deposits: {deposits}")
    print(f"Liabilities central_bank_loans: {bank.get('loans_from_central_banks', {})}")
    print(f"Equity: {bank.get('equity', {})}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "show-balance-sheet", help="Show a simple bank balance sheet."
    )

    parser.add_argument("bank_id")

    parser.set_defaults(func=run)
