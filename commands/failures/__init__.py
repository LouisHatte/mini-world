from argparse import _SubParsersAction

from commands.failures.duplicate_message import register as register_duplicate_message
from commands.failures.fail_settlement import register as register_fail_settlement
from commands.failures.reconcile_all import register as register_reconcile_all
from commands.failures.reject_recipient import register as register_reject_recipient
from commands.failures.repair_mismatch import register as register_repair_mismatch
from commands.failures.retry_message import register as register_retry_message
from commands.failures.simulate_network_failure import (
    register as register_simulate_network_failure,
)


def register_failure_commands(subparsers: _SubParsersAction) -> None:
    register_simulate_network_failure(subparsers)
    register_retry_message(subparsers)
    register_duplicate_message(subparsers)
    register_fail_settlement(subparsers)
    register_reject_recipient(subparsers)
    register_reconcile_all(subparsers)
    register_repair_mismatch(subparsers)
