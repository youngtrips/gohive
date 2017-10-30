#!/bin/bash

openssl genrsa -out rsa_private_key.pem 2048 
openssl rsa -in rsa_private_key.pem -out rsa_public_key.pem -pubout
openssl pkcs8 -topk8 -in rsa_private_key.pem -out pkcs8_rsa_private_key.pem -nocrypt

