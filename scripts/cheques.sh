# Reset world
python3.11 mini.py init --reset

# Setup
python3.11 mini.py create-bank bank1
python3.11 mini.py create-bank bank2
python3.11 mini.py create-human alice
python3.11 mini.py create-human bob
python3.11 mini.py open-account alice bank1 EUR
python3.11 mini.py open-account bob bank2 EUR
python3.11 mini.py grant-loan bank1 alice EUR 300

# Cheque happy path
python3.11 mini.py write-cheque alice bank1 bob EUR 100
python3.11 mini.py deposit-cheque bob bank2 cheque_000001
python3.11 mini.py clear-cheque cheque_000001
python3.11 mini.py settle-cheque cheque_000001

# Bounce path
python3.11 mini.py write-cheque alice bank1 bob EUR 1000
python3.11 mini.py deposit-cheque bob bank2 cheque_000002
python3.11 mini.py bounce-cheque cheque_000002 INSUFFICIENT_FUNDS
python3.11 mini.py reverse-cheque-credit cheque_000002
python3.11 mini.py list-accounts
