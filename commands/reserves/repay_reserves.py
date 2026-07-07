from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


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

    central_bank_reserves = central_bank["reserve_accounts"].get(bank_id, 0)
    bank_reserve_mirror = bank["reserve_balances"].get(central_bank_id, 0)

    if central_bank_reserves != bank_reserve_mirror:
        print("Reserve mirror mismatch. Run check-world.")
        print(
            f"{central_bank_id}.reserve_accounts[{bank_id}] = "
            f"{central_bank_reserves}"
        )
        print(
            f"{bank_id}.reserve_balances[{central_bank_id}] = " f"{bank_reserve_mirror}"
        )
        return

    central_bank_loan = central_bank["loans_to_banks"].get(bank_id, 0)
    bank_loan_mirror = bank["loans_from_central_banks"].get(central_bank_id, 0)

    if central_bank_loan != bank_loan_mirror:
        print("Loan mirror mismatch. Run check-world.")
        print(f"{central_bank_id}.loans_to_banks[{bank_id}] = " f"{central_bank_loan}")
        print(
            f"{bank_id}.loans_from_central_banks[{central_bank_id}] = "
            f"{bank_loan_mirror}"
        )
        return

    if central_bank_reserves < amount:
        print(
            f"Not enough reserves for {bank_id} at {central_bank_id}. "
            f"Available: {central_bank_reserves} {currency}"
        )
        return

    if central_bank_loan < amount:
        print(
            f"{bank_id} does not owe that much to {central_bank_id}. "
            f"Outstanding loan: {central_bank_loan} {currency}"
        )
        return

    # Bank 1 pays with reserves.
    central_bank["reserve_accounts"][bank_id] = central_bank_reserves - amount
    bank["reserve_balances"][central_bank_id] = bank_reserve_mirror - amount

    # The central bank loan is reduced.
    central_bank["loans_to_banks"][bank_id] = central_bank_loan - amount
    bank["loans_from_central_banks"][central_bank_id] = bank_loan_mirror - amount

    save_world(world)

    print(f"{bank_id} repaid {amount} {currency} reserves to {central_bank_id}.")
    print(
        f"{central_bank_id} reserve account for {bank_id}: "
        f"{central_bank['reserve_accounts'][bank_id]} {currency}"
    )
    print(
        f"{bank_id} reserves at {central_bank_id}: "
        f"{bank['reserve_balances'][central_bank_id]} {currency}"
    )
    print(
        f"{central_bank_id} loan to {bank_id}: "
        f"{central_bank['loans_to_banks'][bank_id]} {currency}"
    )
    print(
        f"{bank_id} loan from {central_bank_id}: "
        f"{bank['loans_from_central_banks'][central_bank_id]} {currency}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "repay-reserves", help="Repay a central bank loan using reserves."
    )

    parser.add_argument("central_bank_id", help="Central bank ID, for example: ecb")

    parser.add_argument("bank_id", help="Commercial bank ID, for example: bank1")

    parser.add_argument("amount", type=int, help="Amount of reserves to repay.")

    parser.set_defaults(func=run)
