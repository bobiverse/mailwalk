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

# [102][2024-09-12 08:45:12]     [6 lines] peter@spiderman.example.com       `Weekly Report Summary`
# [119][2024-09-13 10:22:34]     [12 lines] diana@wonderwoman.example.com     `=?utf-8?Q?Server_Maintenance_Notice?=`
# [120][2024-09-13 14:45:09] 📎1 [5 lines] clark@superman.example.com         `Urgent: Website Downtime Alert`
# [121][2024-09-14 09:32:45]     [10 lines] bruce@batman.example.com          `=?utf-8?Q?Invitation_to_Tech_Summit_2024?=`
# [122][2024-09-14 11:15:29]     [18 lines] tony@ironman.example.com          `Security Vulnerability Notification`
# [123][2024-09-14 12:37:03] 📎2 [7 lines] natasha@blackwidow.example.com     `Monthly Financial Report`
# [124][2024-09-15 08:59:11]     [25 lines] steve@captainamerica.example.com  `=?utf-8?Q?Upcoming_Healthcare_Seminar_2024?=`

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
