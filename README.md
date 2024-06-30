```bash

./mailwalk \
    -host $HOST \
    -port $PORT \
    -email $EMAIL \
    -pwd $PASS \
    -folder INBOX \
    -datefrom 2024-06-20 \
    -cmd 'echo "{{MessageId}}; {{DateTime}}; {{Subject}}; {{From}}" | tee -a maildump.txt' 
```
