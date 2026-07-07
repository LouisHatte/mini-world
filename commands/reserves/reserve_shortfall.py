from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    world = load_world()

    central_bank_id = args.central_bank_id
    bank_id = args.bank_id
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if central_bank_id not in world["central_banks"]:
        print(f"Central bank does not exist: {central_bank_id}")
        return

    if bank_id not in world["banks"]:
        print(f"Bank does not exist: {bank_id}")
        return

    central_bank = world["central_banks"][central_bank_id]
    bank = world["banks"][bank_id]
    currency = central_bank["currency"]

    reserve_account = central_bank["reserve_accounts"].get(bank_id, 0)
    reserve_balance = bank["reserve_balances"].get(central_bank_id, 0)

    if reserve_account != reserve_balance:
        print("Reserve mirror mismatch. Run check-world.")
        print(f"{central_bank_id}.reserve_accounts[{bank_id}] = " f"{reserve_account}")
        print(f"{bank_id}.reserve_balances[{central_bank_id}] = " f"{reserve_balance}")
        return

    print(f"Reserve shortfall check for {bank_id} at {central_bank_id}")
    print()
    print(f"Available reserves: {reserve_balance} {currency}")
    print(f"Required reserves: {amount} {currency}")

    if reserve_balance >= amount:
        surplus = reserve_balance - amount

        print()
        print("Result: OK")
        print(f"{bank_id} can settle {amount} {currency}.")
        print(f"Remaining reserves after settlement would be: {surplus} {currency}")
        return

    shortfall = amount - reserve_balance

    print()
    print("Result: SHORTFALL")
    print(f"{bank_id} cannot settle {amount} {currency}.")
    print(f"Missing reserves: {shortfall} {currency}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "reserve-shortfall",
        help="Check whether a bank has enough reserves to settle an amount.",
    )

    parser.add_argument("central_bank_id", help="Central bank ID, for example: ecb")

    parser.add_argument("bank_id", help="Commercial bank ID, for example: bank1")

    parser.add_argument("amount", type=int, help="Amount of reserves required.")

    parser.set_defaults(func=run)
