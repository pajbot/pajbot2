# MySQL to Postgresql transition

1. get pgloader version 3.6.1 (or above?) exists in aur KKona
2. run pgloader transitions files (modify the FROM and INTO lines in the stage1 and stage2 files)
3. install the python virtual environment: `python3 -m venv venv`
4. activate the virtual environment: `source venv/bin/activate`
5. install sql clients: `pip install psycopg2 pymysql`
6. run the "stage3.py" python script: ./stage3.py
