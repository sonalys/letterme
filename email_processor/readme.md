# Email Processor
This microservice is responsible for downloading external attachments, encryption and them storing them at email_bucket.
it's very important to preserve privacy because it avoids getting the user's ip address and read timestamp.

It is also responsible for the encryption of decrypted external messages received.