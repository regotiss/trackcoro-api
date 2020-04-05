LAMBDA_PATH="./cmd/lambda/"
SERVER_PATH="./cmd/server/"

cd $LAMBDA_PATH && set GOOS=linux && set CGO_ENABLED=0 && set GOARCH=amd64 && go build && cd ../..
cd $SERVER_PATH && go build && cd ../..
