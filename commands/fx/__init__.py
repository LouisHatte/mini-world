from argparse import _SubParsersAction

from commands.fx.create_fx_market import register as register_create_fx_market
from commands.fx.cross_border_transfer import register as register_cross_border_transfer
from commands.fx.fx_convert import register as register_fx_convert
from commands.fx.nostro_transfer import register as register_nostro_transfer
from commands.fx.set_fx_rate import register as register_set_fx_rate


def register_fx_commands(subparsers: _SubParsersAction) -> None:
    register_create_fx_market(subparsers)
    register_set_fx_rate(subparsers)
    register_fx_convert(subparsers)
    register_cross_border_transfer(subparsers)
    register_nostro_transfer(subparsers)
