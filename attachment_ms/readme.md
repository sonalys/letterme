# Email bucket
This microservice will be responsible for storing extremely high number of emails, everything needs to be encrypted already or by using email_processor.

Reading or sending emails is done through here.

### Email receiving
Messages received without encryption must be encrypted imediately through email_processor, and a flag must be set, stating that it was received without encryption.

It must have a timestamp, to live through the default time span of the email, before getting deleted.

All attachments and external links detected will need to be encrypted as well.

Messages sent through SMTP will be encrypted.

### Attachment
This service will also handle the attachments, emails received from outside letter.me will have all of it's links and attachments encrypted and stored within this service.

An email can only be sent after all it's attachments are properly sent, to do so, it will need to upload the encrypted files first, then changing the links inside the email body to the files, then encrypting the email and sending it.

All attachments are deleted after ttl of the email.

### Deletion
Emails are deleted by default after ttl, but users can confirm they read emails or attachments by sending a deletion request to the resourceId using a provided JWT token by the account_manager