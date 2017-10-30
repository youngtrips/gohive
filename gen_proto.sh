#!/bin/bash

OS=`uname -s`
GOPATH=`pwd`

if [ "$OS" == "Linux" ]; then
    export PATH=$GOPATH/tools/linux/protoc/bin:$PATH
else
    export PATH=$GOPATH/tools/darwin/protoc/bin:$PATH
fi

echo "which protoc: "
which protoc

#git@lab.qeetap.com:geekjoys/gohive_proto.git
#if test ! -e pb/msg; then
#    if test ! -e .build/pb; then
#        git clone https://lab.qeetap.com/geekjoys/gohive_proto.git .build/pb
#    fi
#    pushd pb
#        rm -rf msg
#        rm -rf def
#        rm -rf cfg
#        ln -s ../.build/pb/msg msg
#        ln -s ../.build/pb/cfg cfg
#        ln -s ../.build/pb/def def
#    popd
#fi

rm ./internal/pb/cfg/*.pb.go
rm ./internal/pb/msg/*.pb.go
rm ./internal/pb/def/*.pb.go

protoc --go_out=import_prefix=gohive/internal/pb/,plugins=grpc:internal/pb pb/cfg/*.proto --proto_path=pb
protoc --go_out=import_prefix=gohive/internal/pb/,plugins=grpc:internal/pb pb/def/*.proto --proto_path=pb
protoc --go_out=import_prefix=gohive/internal/pb/,plugins=grpc:internal/pb pb/msg/*.proto --proto_path=pb
protoc --go_out=import_prefix=gohive/internal/pb/,plugins=grpc:internal/pb pb/service/*.proto --proto_path=pb

if [ "$OS" == "Linux" ]; then
    find ./internal/pb -name "*.pb.go" | xargs sed -i "s/gohive\/internal\/pb\/github.com/github.com/g"
    find ./internal/pb -name "*.pb.go" | xargs sed -i "s/gohive\/internal\/pb\/golang.org/golang.org/g"
    find ./internal/pb -name "*.pb.go" | xargs sed -i "s/gohive\/internal\/pb\/google.golang.org/google.golang.org/g"
    find ./internal/pb -name "*.pb.go" | xargs sed -i "s/,omitempty//g"
else
    find ./internal/pb -name "*.pb.go" | xargs sed -i '' "s/gohive\/internal\/pb\/github.com/github.com/g"
    find ./internal/pb -name "*.pb.go" | xargs sed -i '' "s/gohive\/internal\/pb\/golang.org/golang.org/g"
    find ./internal/pb -name "*.pb.go" | xargs sed -i '' "s/gohive\/internal\/pb\/google.golang.org/google.golang.org/g"
    find ./internal/pb -name "*.pb.go" | xargs sed -i '' "s/,omitempty//g"
fi


types=`cat internal/pb/def/common.pb.go| grep "type" | cut -d ' ' -f2`
for type in ${types[@]};do
    if [ "$OS" == "Linux" ]; then
        sed -i "s/^\t${type}_/\t/g" internal/pb/def/common.pb.go
    else
        #sed -e 's/\/[CTR+V][CTR+I]/g' filename
        sed -i '' "s/^	${type}_/	/g" internal/pb/def/common.pb.go
    fi
done

./gen_handler.py ./pb ./server/game ./server/game/handler

