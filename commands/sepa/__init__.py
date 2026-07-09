from argparse import _SubParsersAction

from commands.sepa.bank_credit_incoming import register as register_bank_credit_incoming
from commands.sepa.create_step2 import register as register_create_step2
from commands.sepa.sepa_reconcile import register as register_sepa_reconcile
from commands.sepa.sepa_return import register as register_sepa_return
from commands.sepa.sepa_transfer import register as register_sepa_transfer
from commands.sepa.step2_clear import register as register_step2_clear
from commands.sepa.step2_deliver import register as register_step2_deliver
from commands.sepa.step2_receive import register as register_step2_receive
from commands.sepa.t2_settle import register as register_t2_settle


def register_sepa_commands(subparsers: _SubParsersAction) -> None:
    register_create_step2(subparsers)
    register_sepa_transfer(subparsers)
    register_step2_receive(subparsers)
    register_step2_clear(subparsers)
    register_t2_settle(subparsers)
    register_step2_deliver(subparsers)
    register_bank_credit_incoming(subparsers)
    register_sepa_return(subparsers)
    register_sepa_reconcile(subparsers)
