# Reset world
python3.11 mini.py init --reset

# Setup and bond lifecycle
python3.11 mini.py create-central-bank ecb EUR
python3.11 mini.py create-bank bank1
python3.11 mini.py issue-bond treasury bond1 EUR 1000
python3.11 mini.py buy-security ecb bank1 bond1 1000
python3.11 mini.py mark-to-market bond1 950
python3.11 mini.py pay-coupon bond1 ecb 50
python3.11 mini.py sell-security ecb bank1 bond1 500
python3.11 mini.py redeem-bond bond1 ecb 500
python3.11 mini.py show-balance-sheet bank1
python3.11 mini.py check-world
