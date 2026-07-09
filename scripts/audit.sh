# Reset world
python3.11 mini.py init --reset

# Setup
python3.11 mini.py create-central-bank ecb EUR
python3.11 mini.py create-bank bank1
python3.11 mini.py create-human alice
python3.11 mini.py open-account alice bank1 EUR
python3.11 mini.py grant-loan bank1 alice EUR 100

# Debugging and snapshots
python3.11 mini.py list-banks
python3.11 mini.py list-humans
python3.11 mini.py list-accounts
python3.11 mini.py show-account acc_bank1_alice_eur
python3.11 mini.py show-balance-sheet bank1
python3.11 mini.py snapshot before_payment
python3.11 mini.py show-json
python3.11 mini.py load-snapshot before_payment
