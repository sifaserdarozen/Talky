# Secure Transport Certificates
NEVER PUT PRODUCTION CERTIFICATES TO REPOSITORY
ALL FILES EXCEPT THIS WILL BE IGNORED BY.GITIGNORE

If server certificates **certificate.key** & **sertificate.crt** are put here, 
they will be copied to container at **home/asterisk/certs/**

When there are no available ones, they will be generated inside the container
as self signed.