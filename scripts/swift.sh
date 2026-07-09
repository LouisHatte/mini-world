# Reset world
python3.11 mini.py init --reset

# Setup
python3.11 mini.py create-bank bank1
python3.11 mini.py create-bank bank2
python3.11 mini.py create-bank bankX
python3.11 mini.py create-bank bankY
python3.11 mini.py create-human alice
python3.11 mini.py create-human bob
python3.11 mini.py open-account alice bank1 USD
python3.11 mini.py open-account bob bank2 USD
python3.11 mini.py grant-loan bank1 alice USD 500
python3.11 mini.py create-correspondent-account bank1 bankX USD
python3.11 mini.py create-correspondent-account bank2 bankY USD

# SWIFT-style flow
python3.11 mini.py swift-transfer alice bank1 bob bank2 USD 100
python3.11 mini.py swift-message bank1 bank2 MT103
python3.11 mini.py correspondent-debit bankX bank1 USD 100
python3.11 mini.py correspondent-credit bankY bank2 USD 100
python3.11 mini.py deduct-fee bankX payment_000001 10
python3.11 mini.py swift-reconcile payment_000001
