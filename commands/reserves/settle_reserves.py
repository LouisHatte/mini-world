from argparse import Namespace, _SubParsersAction

from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()

    central_bank_id = args.central_bank_id
    source_bank_id = args.source_bank_id
    target_bank_id = args.target_bank_id
    amount = args.amount

    if amount <= 0:
        print("Amount must be greater than 0.")
        return

    if source_bank_id == target_bank_id:
        print("Source bank and target bank must be different.")
        return

    if central_bank_id not in world["central_banks"]:
        print(f"Central bank does not exist: {central_bank_id}")
        return

    if source_bank_id not in world["banks"]:
        print(f"Source bank does not exist: {source_bank_id}")
        return

    if target_bank_id not in world["banks"]:
        print(f"Target bank does not exist: {target_bank_id}")
        return

    central_bank = world["central_banks"][central_bank_id]
    source_bank = world["banks"][source_bank_id]
    target_bank = world["banks"][target_bank_id]
    currency = central_bank["currency"]

    source_reserves_at_central_bank = central_bank["reserve_accounts"].get(
        source_bank_id, 0
    )
    source_reserve_mirror = source_bank["reserve_balances"].get(central_bank_id, 0)

    if source_reserves_at_central_bank != source_reserve_mirror:
        print("Reserve mirror mismatch. Run check-world.")
        print(
            f"{central_bank_id}.reserve_accounts[{source_bank_id}] = "
            f"{source_reserves_at_central_bank}"
        )
        print(
            f"{source_bank_id}.reserve_balances[{central_bank_id}] = "
            f"{source_reserve_mirror}"
        )
        return

    if source_reserves_at_central_bank < amount:
        print(
            f"Not enough reserves for {source_bank_id} at {central_bank_id}. "
            f"Available: {source_reserves_at_central_bank} {currency}"
        )
        return

    target_reserves_at_central_bank = central_bank["reserve_accounts"].get(
        target_bank_id, 0
    )
    target_reserve_mirror = target_bank["reserve_balances"].get(central_bank_id, 0)

    if target_reserves_at_central_bank != target_reserve_mirror:
        print("Reserve mirror mismatch. Run check-world.")
        print(
            f"{central_bank_id}.reserve_accounts[{target_bank_id}] = "
            f"{target_reserves_at_central_bank}"
        )
        print(
            f"{target_bank_id}.reserve_balances[{central_bank_id}] = "
            f"{target_reserve_mirror}"
        )
        return

    # Central bank ledger:
    # The central bank owes less reserves to the source bank.
    central_bank["reserve_accounts"][source_bank_id] = (
        source_reserves_at_central_bank - amount
    )

    # The central bank owes more reserves to the target bank.
    central_bank["reserve_accounts"][target_bank_id] = (
        target_reserves_at_central_bank + amount
    )

    # Commercial bank mirrors:
    # Source bank has fewer reserves at the central bank.
    source_bank["reserve_balances"][central_bank_id] = source_reserve_mirror - amount

    # Target bank has more reserves at the central bank.
    target_bank["reserve_balances"][central_bank_id] = target_reserve_mirror + amount

    save_world(world)

    print(
        f"Settled {amount} {currency} reserves from "
        f"{source_bank_id} to {target_bank_id} at {central_bank_id}."
    )
    print(
        f"{central_bank_id} reserve account for {source_bank_id}: "
        f"{central_bank['reserve_accounts'][source_bank_id]} {currency}"
    )
    print(
        f"{central_bank_id} reserve account for {target_bank_id}: "
        f"{central_bank['reserve_accounts'][target_bank_id]} {currency}"
    )
    print(
        f"{source_bank_id} reserves at {central_bank_id}: "
        f"{source_bank['reserve_balances'][central_bank_id]} {currency}"
    )
    print(
        f"{target_bank_id} reserves at {central_bank_id}: "
        f"{target_bank['reserve_balances'][central_bank_id]} {currency}"
    )


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "settle-reserves",
        help="Move central bank reserves from one commercial bank to another.",
    )

    parser.add_argument("central_bank_id", help="Central bank ID, for example: ecb")

    parser.add_argument(
        "source_bank_id", help="Bank paying reserves, for example: bank1"
    )

    parser.add_argument(
        "target_bank_id", help="Bank receiving reserves, for example: bank2"
    )

    parser.add_argument("amount", type=int, help="Amount of reserves to settle.")

    parser.set_defaults(func=run)
