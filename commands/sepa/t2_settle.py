from argparse import Namespace, _SubParsersAction

from commands.common import next_id
from commands.sepa.helpers import find_central_bank
from world import load_world, save_world


def run(args: Namespace) -> None:
    world = load_world()
    payment_id = args.payment_id

    if payment_id not in world["payment_instructions"]:
        print(f"Payment does not exist: {payment_id}")
        return

    payment = world["payment_instructions"][payment_id]
    currency = payment["currency"]
    central_bank_id = find_central_bank(world, currency)

    if central_bank_id is None:
        print(f"No central bank for {currency}.")
        return

    central_bank = world["central_banks"][central_bank_id]
    source_bank = world["banks"][payment["sender_bank_id"]]
    target_bank = world["banks"][payment["recipient_bank_id"]]
    amount = payment["amount"]

    source_reserves = central_bank["reserve_accounts"].get(
        payment["sender_bank_id"], 0
    )
    target_reserves = central_bank["reserve_accounts"].get(
        payment["recipient_bank_id"], 0
    )

    if source_reserves < amount:
        print(
            f"Settlement failed. Missing reserves: "
            f"{amount - source_reserves} {currency}"
        )
        payment["status"] = "SETTLEMENT_FAILED"
        save_world(world)
        return

    central_bank["reserve_accounts"][payment["sender_bank_id"]] = (
        source_reserves - amount
    )
    central_bank["reserve_accounts"][payment["recipient_bank_id"]] = (
        target_reserves + amount
    )
    source_bank["reserve_balances"][central_bank_id] = source_reserves - amount
    target_bank["reserve_balances"][central_bank_id] = target_reserves + amount

    settlement_id = next_id(world["settlements"], "settlement")
    world["settlements"][settlement_id] = {
        "id": settlement_id,
        "payment_id": payment_id,
        "rail": "T2",
        "central_bank_id": central_bank_id,
        "status": "SETTLED",
        "currency": currency,
        "amount": amount,
    }
    payment["settlement_id"] = settlement_id
    payment["status"] = "SETTLED"

    save_world(world)

    print(f"T2 settled payment: {payment_id}")
    print(f"Settlement: {settlement_id}")


def register(subparsers: _SubParsersAction) -> None:
    parser = subparsers.add_parser(
        "t2-settle", help="Settle a SEPA payment in central bank money."
    )

    parser.add_argument("payment_id")

    parser.set_defaults(func=run)
