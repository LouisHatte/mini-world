# Reset world
python3.11 mini.py init --reset

# Setup
python3.11 mini.py create-bank bank1
python3.11 mini.py create-bank bank2
python3.11 mini.py create-bank bankX
python3.11 mini.py create-human alice
python3.11 mini.py create-human bob
python3.11 mini.py open-account alice bank1 EUR
python3.11 mini.py open-account alice bank1 USD
python3.11 mini.py open-account bob bank2 USD
python3.11 mini.py grant-loan bank1 alice EUR 1000
python3.11 mini.py create-correspondent-account bank1 bankX USD
python3.11 mini.py create-fx-market EUR USD 1.10

# FX and cross-border flows
python3.11 mini.py fx-convert alice bank1 EUR USD 100
python3.11 mini.py fx-convert bank1 EUR USD 250
python3.11 mini.py set-fx-rate EUR USD 1.08
python3.11 mini.py cross-border-transfer alice bank1 bob bank2 EUR USD 100
python3.11 mini.py nostro-transfer bank1 bankX USD 100
python3.11 mini.py list-accounts
