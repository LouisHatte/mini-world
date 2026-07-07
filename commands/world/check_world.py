from argparse import Namespace, _SubParsersAction

from world import load_world


def check_reserve_mirrors(world: dict) -> list[str]:
    errors: list[str] = []

    central_banks = world["central_banks"]
    banks = world["banks"]

    for central_bank_id, central_bank in central_banks.items():
        for bank_id, reserve_amount in central_bank["reserve_accounts"].items():
            if bank_id not in banks:
                errors.append(
                    f"{central_bank_id}.reserve_accounts contains unknown bank: {bank_id}"
                )
                continue

            bank = banks[bank_id]
            mirrored_amount = bank["reserve_balances"].get(central_bank_id)

            if mirrored_amount != reserve_amount:
                errors.append(
                    f"Reserve mismatch: "
                    f"{central_bank_id}.reserve_accounts[{bank_id}] = {reserve_amount}, "
                    f"but {bank_id}.reserve_balances[{central_bank_id}] = {mirrored_amount}"
                )

    for bank_id, bank in banks.items():
        for central_bank_id, mirrored_amount in bank["reserve_balances"].items():
            if central_bank_id not in central_banks:
                errors.append(
                    f"{bank_id}.reserve_balances contains unknown central bank: {central_bank_id}"
                )
                continue

            central_bank = central_banks[central_bank_id]
            reserve_amount = central_bank["reserve_accounts"].get(bank_id)

            if reserve_amount != mirrored_amount:
                errors.append(
                    f"Reserve mismatch: "
                    f"{bank_id}.reserve_balances[{central_bank_id}] = {mirrored_amount}, "
                    f"but {central_bank_id}.reserve_accounts[{bank_id}] = {reserve_amount}"
                )

    return errors


def check_loan_mirrors(world: dict) -> list[str]:
    errors: list[str] = []

    central_banks = world["central_banks"]
    banks = world["banks"]

    for central_bank_id, central_bank in central_banks.items():
        for bank_id, loan_amount in central_bank["loans_to_banks"].items():
            if bank_id not in banks:
                errors.append(
                    f"{central_bank_id}.loans_to_banks contains unknown bank: {bank_id}"
                )
                continue

            bank = banks[bank_id]
            mirrored_amount = bank["loans_from_central_banks"].get(central_bank_id)

            if mirrored_amount != loan_amount:
                errors.append(
                    f"Loan mismatch: "
                    f"{central_bank_id}.loans_to_banks[{bank_id}] = {loan_amount}, "
                    f"but {bank_id}.loans_from_central_banks[{central_bank_id}] = {mirrored_amount}"
                )

    for bank_id, bank in banks.items():
        for central_bank_id, mirrored_amount in bank[
            "loans_from_central_banks"
        ].items():
            if central_bank_id not in central_banks:
                errors.append(
                    f"{bank_id}.loans_from_central_banks contains unknown central bank: {central_bank_id}"
                )
                continue

            central_bank = central_banks[central_bank_id]
            loan_amount = central_bank["loans_to_banks"].get(bank_id)

            if loan_amount != mirrored_amount:
                errors.append(
                    f"Loan mismatch: "
                    f"{bank_id}.loans_from_central_banks[{central_bank_id}] = {mirrored_amount}, "
                    f"but {central_bank_id}.loans_to_banks[{bank_id}] = {loan_amount}"
                )

    return errors


def run(args: Namespace) -> None:
    world = load_world()

    errors: list[str] = []
    errors.extend(check_reserve_mirrors(world))
    errors.extend(check_loan_mirrors(world))

    if not errors:
        print("World check passed.")
        print("Reserve mirrors are consistent.")
        print("Loan mirrors are consistent.")
        return

    print("World check failed.")
    print()

    for error in errors:
        print(f"- {error}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "check-world",
        help="Check consistency between central banks and commercial banks.",
    )

    parser.set_defaults(func=run)
