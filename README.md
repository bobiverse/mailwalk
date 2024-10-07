# 📨 **mailwalk**

Simple tool to read and loop through your mailbox, automating tasks by running custom commands on each email. 
For automating email parsing, notification systems, and more.


## **Features**  
- Read emails from your mailbox  
- Loop through each email and run custom commands  
- Extract specific fields like ```Subject``` , ```From``` , ```Cc``` , etc.  
- Supports filtering by folder (e.g., **INBOX**, **SENT**)  
- Filter emails by date range for more specific results  


## **Usage Example**

Here’s a basic example to extract email details and output the information to both a file and stdout.

```bash  
./mailwalk \
    -host $HOST \
    --port $PORT \
    -email $EMAIL \
    -pwd $PASS \
    -folder INBOX \
    -datefrom 2024-06-20 \
    -cmd 'echo "{{MessageId}}; {{DateTime}}; {{Subject}}; {{From}}" | tee -a maildump.txt'  
```  

- `-host`: Mail server host (e.g., *imap.gmail.com*)  
- `-port`: Mail server port (e.g., *993*)  
- `-email`: Your email address  
- `-pwd`: Your email password  
- `-folder`: Mailbox folder to read (default is **INBOX**)  
- `-datefrom`: Only read emails from this date onward  
- `-cmd`: Command to execute for each email, using placeholders for email fields  

## **Available Placeholders**

You can include the following placeholders in your commands to extract specific information from each email:

| Placeholder   | Description                   |
|---------------|-------------------------------|
| ```Uid```     | Unique identifier of the email |
| ```MessageId```| ID of the message             |
| ```Subject``` | Email subject                 |
| ```From```    | Sender's email address         |
| ```Cc```      | Cc recipients                 |
| ```Bcc```     | Bcc recipients                |
| ```ReplyTo``` | Reply-to address              |
| ```DateTime```| Date and time of the email     |
| ```Text```    | Plain text body of the email   |
| ```HTML```    | HTML body of the email         |


## **Roadmap**

- **Multi-threading**: Implement multi-threaded email processing to handle large mailboxes more efficiently.
- **Better Error Handling**: Improve error handling for common issues like invalid credentials, network failures, and missing email fields.
- **Support More Email Protocols**: Add support for additional protocols like POP3 to broaden compatibility.
- **Advanced Filters**: Introduce more filtering options, such as filtering by keywords in the subject, body, or specific header fields.
