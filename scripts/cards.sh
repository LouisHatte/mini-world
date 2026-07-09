# Reset world
python3.11 mini.py init --reset

# Setup
python3.11 mini.py create-bank bank1
python3.11 mini.py create-human alice
python3.11 mini.py open-account alice bank1 EUR
python3.11 mini.py grant-loan bank1 alice EUR 300

# Capture path
python3.11 mini.py card-authorize alice merchant EUR 50
python3.11 mini.py card-capture auth_000001

# Reverse and chargeback paths
python3.11 mini.py card-authorize alice merchant EUR 40
python3.11 mini.py card-reverse auth_000002
python3.11 mini.py card-chargeback payment_000001
python3.11 mini.py show-account acc_bank1_alice_eur
