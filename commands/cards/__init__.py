from argparse import _SubParsersAction

from commands.cards.card_authorize import register as register_card_authorize
from commands.cards.card_capture import register as register_card_capture
from commands.cards.card_chargeback import register as register_card_chargeback
from commands.cards.card_reverse import register as register_card_reverse


def register_card_commands(subparsers: _SubParsersAction) -> None:
    register_card_authorize(subparsers)
    register_card_capture(subparsers)
    register_card_reverse(subparsers)
    register_card_chargeback(subparsers)
