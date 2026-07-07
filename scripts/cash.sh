# Reset world
python3.11 mini.py init --reset

# Create actors
python3.11 mini.py create-central-bank ecb EUR
python3.11 mini.py create-bank bank1
python3.11 mini.py create-bank bank2
python3.11 mini.py create-human alice
python3.11 mini.py create-human bob

# Create bank accounts
python3.11 mini.py open-account alice bank1 EUR
python3.11 mini.py open-account bob bank1 EUR

# Create physical cash at the central bank
python3.11 mini.py issue-cash ecb 10000

# Give banks reserves so they can later request cash from the ECB
python3.11 mini.py lend-reserves ecb bank1 1000
python3.11 mini.py lend-reserves ecb bank2 500

# 1. Initial cash distribution to a human
python3.11 mini.py seed-cash ecb alice 1000

# 2. Human physically gives cash to another human
python3.11 mini.py transfer-cash alice bob EUR 100

# 3. Human deposits cash into a commercial bank
python3.11 mini.py deposit-cash alice bank1 EUR 500

# 4. Human withdraws cash from a bank deposit
python3.11 mini.py withdraw-cash alice bank1 EUR 200

# 5. Commercial bank converts reserves into physical cash from ECB
python3.11 mini.py supply-cash ecb bank1 700

# 6. Physical cash transport between commercial banks
python3.11 mini.py move-cash bank1 bank2 EUR 250

# 7. Commercial bank returns physical cash to the central bank
python3.11 mini.py return-cash ecb bank2 100

# 8. Central bank destroys physical cash from its own vault
python3.11 mini.py destroy-cash ecb 50

# Final consistency check
python3.11 mini.py check-world