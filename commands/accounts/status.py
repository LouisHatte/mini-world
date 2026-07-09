from argparse import Namespace

from world import load_world, save_world


def set_account_status(args: Namespace, status: str) -> None:
    world = load_world()
    account_id = args.account_id

    if account_id not in world["accounts"]:
        print(f"Account does not exist: {account_id}")
        return

    account = world["accounts"][account_id]

    if status == "CLOSED" and account["booked_balance"] != 0:
        print(
            "Cannot close account with non-zero balance: "
            f"{account['booked_balance']} {account['currency']}"
        )
        return

    account["status"] = status

    save_world(world)

    print(f"{account_id} status: {status}")
