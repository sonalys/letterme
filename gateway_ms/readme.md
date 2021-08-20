# Spam Firewall
This microservice is responsible for identifying blacklisted emails, reported external virus or scam links received through unencrypted channels or sha-512 from attachments.

It will have an indexed collection of strings, which will be the email address, the external link or the sha-512 from attachments, if present, it will set a flag for spam.

This microservice will need to collect external data about these threats.

It will also need to verify for denial of services, if an ip address or domain exceeds the threshold, it will be blacklisted imediately, and all emails sent will be deleted.

The blacklisted ip will be discarded after ttl.