# Reset world
python3.11 mini.py init --reset

# Setup
python3.11 mini.py create-central-bank ecb EUR
python3.11 mini.py issue-cash ecb 1000
python3.11 mini.py create-bank bank1
python3.11 mini.py create-human alice
python3.11 mini.py create-human bob
python3.11 mini.py open-account alice bank1 EUR
python3.11 mini.py open-account bob bank1 EUR
python3.11 mini.py seed-cash ecb alice 500
python3.11 mini.py deposit-cash alice bank1 EUR 500

# Holds and available-balance behavior
python3.11 mini.py place-hold alice bank1 EUR 200
python3.11 mini.py internal-transfer alice bob bank1 EUR 300
python3.11 mini.py release-hold hold_000001
python3.11 mini.py place-hold alice bank1 EUR 100
python3.11 mini.py capture-hold hold_000002

# Account state controls
python3.11 mini.py freeze-account acc_bank1_alice_eur
python3.11 mini.py unfreeze-account acc_bank1_alice_eur
python3.11 mini.py internal-transfer alice bob bank1 EUR 100

# Audit
python3.11 mini.py list-accounts
python3.11 mini.py check-world
