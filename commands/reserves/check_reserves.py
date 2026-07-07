from argparse import Namespace, _SubParsersAction

from world import load_world


def run(args: Namespace) -> None:
    world = load_world()

    central_bank_id = args.central_bank_id
    bank_id = args.bank_id

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

    loan_to_bank = central_bank["loans_to_banks"].get(bank_id, 0)
    loan_from_central_bank = bank["loans_from_central_banks"].get(central_bank_id, 0)

    print(f"Reserve position for {bank_id} at {central_bank_id}")
    print()

    print("Reserves:")
    print(
        f"  {central_bank_id}.reserve_accounts[{bank_id}]: "
        f"{reserve_account} {currency}"
    )
    print(
        f"  {bank_id}.reserve_balances[{central_bank_id}]: "
        f"{reserve_balance} {currency}"
    )

    if reserve_account == reserve_balance:
        print("  Reserve mirror: OK")
    else:
        print("  Reserve mirror: MISMATCH")

    print()

    print("Central bank loan:")
    print(
        f"  {central_bank_id}.loans_to_banks[{bank_id}]: " f"{loan_to_bank} {currency}"
    )
    print(
        f"  {bank_id}.loans_from_central_banks[{central_bank_id}]: "
        f"{loan_from_central_bank} {currency}"
    )

    if loan_to_bank == loan_from_central_bank:
        print("  Loan mirror: OK")
    else:
        print("  Loan mirror: MISMATCH")

    print()

    print("Interpretation:")
    print(
        f"  {bank_id} has {reserve_balance} {currency} in reserves at {central_bank_id}."
    )
    print(
        f"  {bank_id} owes {loan_from_central_bank} {currency} to {central_bank_id} as a loan."
    )

    net_position = reserve_balance - loan_from_central_bank

    print()
    print("Simplified net position:")
    print(
        f"  reserves - central bank loan = "
        f"{reserve_balance} - {loan_from_central_bank} = {net_position} {currency}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "check-reserves", help="Show a bank's reserve and central bank loan position."
    )

    parser.add_argument("central_bank_id", help="Central bank ID, for example: ecb")

    parser.add_argument("bank_id", help="Commercial bank ID, for example: bank1")

    parser.set_defaults(func=run)
