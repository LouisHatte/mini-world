# Reset world
python3.11 mini.py init --reset

# Setup
python3.11 mini.py create-central-bank ecb EUR
python3.11 mini.py create-bank bank1
python3.11 mini.py create-bank bank2
python3.11 mini.py create-human alice
python3.11 mini.py create-human bob
python3.11 mini.py open-account alice bank1 EUR
python3.11 mini.py open-account bob bank2 EUR
python3.11 mini.py lend-reserves ecb bank1 1000
python3.11 mini.py lend-reserves ecb bank2 100
python3.11 mini.py grant-loan bank1 alice EUR 500
python3.11 mini.py create-step2 step2 EUR

# SEPA flow
python3.11 mini.py sepa-transfer alice bank1 bob bank2 EUR 100
python3.11 mini.py step2-receive msg_000001
python3.11 mini.py step2-clear payment_000001
python3.11 mini.py t2-settle payment_000001
python3.11 mini.py step2-deliver payment_000001
python3.11 mini.py bank-credit-incoming bank2 payment_000001
python3.11 mini.py sepa-reconcile payment_000001
python3.11 mini.py check-world
