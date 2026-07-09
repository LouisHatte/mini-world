from argparse import _SubParsersAction

from commands.securities.buy_security import register as register_buy_security
from commands.securities.issue_bond import register as register_issue_bond
from commands.securities.mark_to_market import register as register_mark_to_market
from commands.securities.pay_coupon import register as register_pay_coupon
from commands.securities.redeem_bond import register as register_redeem_bond
from commands.securities.sell_security import register as register_sell_security


def register_security_commands(subparsers: _SubParsersAction) -> None:
    register_issue_bond(subparsers)
    register_buy_security(subparsers)
    register_sell_security(subparsers)
    register_pay_coupon(subparsers)
    register_redeem_bond(subparsers)
    register_mark_to_market(subparsers)
