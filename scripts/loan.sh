# Reset world
python3.11 mini.py init --reset

# Basic setup
python3.11 mini.py create-central-bank ecb EUR
python3.11 mini.py issue-cash ecb 10000
python3.11 mini.py create-bank bank1
python3.11 mini.py create-human alice
python3.11 mini.py open-account alice bank1 EUR

# Give Alice a bit of extra money so she can pay interest later
python3.11 mini.py seed-cash ecb alice 100
python3.11 mini.py deposit-cash alice bank1 EUR 100

# =========================
# Loan path 1: grant -> accrue interest -> repay fully
# =========================

python3.11 mini.py grant-loan bank1 alice EUR 1000
python3.11 mini.py accrue-interest loan_000001 50

# Alice should now owe 1050 and have 1100 booked balance:
# 100 initial deposit + 1000 loan deposit

python3.11 mini.py repay-loan alice bank1 loan_000001 50
python3.11 mini.py repay-loan alice bank1 loan_000001 1000

# Expected:
# loan_000001 status: REPAID
# outstanding_principal: 0
# outstanding_interest: 0
# bank1 interest_income[EUR]: 50
# bank1 equity[EUR]: 50

# =========================
# Loan path 2: grant -> default -> cure -> repay partially
# =========================

python3.11 mini.py grant-loan bank1 alice EUR 500
python3.11 mini.py default-loan loan_000002
python3.11 mini.py cure-loan loan_000002
python3.11 mini.py repay-loan alice bank1 loan_000002 100

# Expected:
# loan_000002 status: ACTIVE
# outstanding_principal: 400

# =========================
# Loan path 3: grant -> accrue interest -> default -> write off
# =========================

python3.11 mini.py grant-loan bank1 alice EUR 700
python3.11 mini.py accrue-interest loan_000003 30
python3.11 mini.py default-loan loan_000003
python3.11 mini.py write-off-loan loan_000003

# Expected:
# loan_000003 status: WRITTEN_OFF
# outstanding_principal: 0
# outstanding_interest: 0
# written_off_principal: 700
# written_off_interest: 30
# bank1 loan_loss_expense[EUR]: 730
# bank1 equity[EUR]: -680
#
# Why -680?
# +50 interest income from loan_000001
# -730 write-off loss from loan_000003
# = -680 equity

# Final consistency check
python3.11 mini.py check-world