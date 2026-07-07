from argparse import _SubParsersAction

from commands.world.init_world import register as register_init_world
from commands.world.check_world import register as register_check_world


def register_world_commands(subparsers: _SubParsersAction) -> None:
    register_init_world(subparsers)
    register_check_world(subparsers)
