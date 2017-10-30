#!/bin/bash

export PATH=$PATH:"./bin"

hostname=`hostname`
address=`ifconfig -a | grep inet | grep -v "127.0.0.1" | grep -v inet6 | awk '{print $2}' | tr -d "addr:"`


ips=($address)

address='127.0.0.1'
for ip in ${ips[@]};do
    address=${address},${ip}
done

echo "hostname: "$hostname
echo "address : "$address

echo '{"CN":"CN","key":{"algo":"rsa","size":2048}}' | cfssl gencert -initca - | cfssljson -bare ca -
echo '{"signing":{"default":{"expiry":"43800h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json

#
export ADDRESS=$address
export NAME=next
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare server

#
export ADDRESS=
export NAME=client
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME

