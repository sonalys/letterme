# Email Micro-Service
This microservice is responsible for managing all tasks related to emails.

From receiving / sending emails, to / from outside letter.me or inner emails, it will handle all the cryptographic
work needed, as well, will trigger all the required tasks for other micro-services to do.

## Basic Concepts
Discussion of all the usecases of this micro-service.

### Incoming emails
Incoming emails are treated differently depending on the source.

- External emails: those received from outside letter.me will reach the server as decrypted messages, they will need to be analyzed, purged from any tracking techniques, all of it's body media will need to be re-uploaded into our buckets.

- Internal emails: internal emails will come encrypted, with the exception of the recipient, because we need this value to properly address this email to the right person.

#### Multiple recipients
Emails can have multiple recipients, so we need to create a copy of the received email for each user inside our system, and properly encrypt it using the recipient's public key.

### Outgoing emails
Emails that are going to the outside will not have end to end encryption ( e2e ), because the SMTP that will receive that email will not be able to interpret and handle it correctly, due to the lack of protocols of e2e encryption on emails.
We will not store any email going outside, so it will be deleted imediately after sent.

Attachments of outgoing emails to outside letter.me will also need special treatment, it's not defined yet how we can implement this, but we could try allowing attachments to be downloaded a limited amount of times, or set a TTL for the files to exist.