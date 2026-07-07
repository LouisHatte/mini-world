import argparse
import copy
import sys
from datetime import datetime, timezone
from typing import Any

from commands import register_all_commands
from world import load_world, save_world, world_exists


def world_without_command_history(
    world: dict[str, Any] | None,
) -> dict[str, Any] | None:
    if world is None:
        return None

    comparable_world = copy.deepcopy(world)
    comparable_world.pop("command_history", None)

    return comparable_world


def load_world_if_exists() -> dict[str, Any] | None:
    if not world_exists():
        return None

    return load_world()


def append_command_history_if_world_changed(
    before_world: dict[str, Any] | None, argv: list[str]
) -> None:
    if not world_exists():
        return

    after_world = load_world()

    before_comparable = world_without_command_history(before_world)
    after_comparable = world_without_command_history(after_world)

    if before_comparable == after_comparable:
        return

    command_history = after_world.setdefault("command_history", [])

    command_history.append(
        {
            "id": len(command_history) + 1,
            "timestamp_utc": datetime.now(timezone.utc).isoformat(),
            "command": argv[0] if argv else None,
            "argv": argv,
        }
    )

    save_world(after_world)


def main() -> None:
    parser = argparse.ArgumentParser(
        prog="mini", description="A small monetary and banking simulation CLI."
    )

    subparsers = parser.add_subparsers(dest="command", required=True)

    register_all_commands(subparsers)

    args = parser.parse_args()

    argv = sys.argv[1:]
    before_world = load_world_if_exists()

    args.func(args)

    append_command_history_if_world_changed(before_world, argv)


if __name__ == "__main__":
    main()
