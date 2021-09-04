# Letter.Me

Letter me is a proof of concept ( POC ) of a totally private, end2end encrypted, distributed system.

It creates the possibility of business emails, that will only be stored inside the company cluster, or regular emails, stored by letter.me

## Basic concepts

Discussion of all the usecases of this microservice

### Account creation

An account must be created only by providing an available email address

It will give you a private key, which will be used to read emails from others

We will store a public key, so others can send you an email.

### Message retention

By default, messages will be stored by 30 days before deletion on the server-side

After synchronizing your device with the api, it will fetch and delete all your emails from the server.

### Receiving emails

To receive emails, others will collect your public key, encrypt the desired message with it, and send to letter.me, the message will be stored by the default period of time.

### Sending emails

To send emails, you will provide an email address and it will fetch the public key for that address, you will encrypt it using the public key, so only the owner of the private key can read.

### Group emails

When sending an email to multiple users, it will generate a sha-512 key to sign the email, and encrypt it individually with each of the recipients public keys.

### Synchronizing multiple devices

Initially, this will not be possible, could be one of those:

- Create a new email address that will receive all of your old emails stored on your device, then change the email address to the same one, synchronizing the private key from the original one encrypting through the public key of the second one.

### Email confirmation

To verify if an email was sent to outside of the letter.me domain, other email services will need to verify it within our system, I don't have any idea how that works, probably with the public key, otherwise we will have to store email history.
