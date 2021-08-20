# Account Micro-Service
This microservice is responsible for:
- Creating new accounts
- Fetching information about any account
- Retrieval of public key

## Basic concepts
Discussion of all the usecases of this microservice

### Account creation
An account must be created only by providing an available email address

It will give you a private key, which will be used to read emails from others

We will store a public key, so others can send you an email.

### Fetching information
When sending an email to someone, it might be possible to fetch basic information about that email address, such as profile picture, name, etc...

The user will have to opt-in or opt-out about giving this informations to others, but he can always fetch it for himself.

To do so, he will need to encrypt these informations using his private key.

### Public key
We will need to have a blazing fast cache converting email addresses to public keys, one of the main throttles of the project.

### Authencity check
This service will be responsible for creating short lived JWT tokens for users.
It encrypts it using the public key, and it can be used to confirm receivement of attachments and reading of emails.