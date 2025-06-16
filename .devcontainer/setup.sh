#!/bin/bash
set -e  # Stop execution if a command fails

# echo "Downloading dependencies..."
# go mod download

echo "Installing AIR..."
go install github.com/air-verse/air@v1.61.7

echo "Installing SQLC..."
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

echo "Installing other dependencies..."
sudo apt update && sudo apt install -y curl jq

echo "Installing db migrate..."
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3

echo "Installing mock..."
go install github.com/golang/mock/mockgen@v1.6.0

echo "Installing protoc..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/ktr0731/evans@latest


# echo "Installing grpc-gateway..."
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo "Installing serve static files..."
go install github.com/rakyll/statik


echo "Installing aws cli"
pushd /tmp
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "./awscliv2.zip"
unzip ./awscliv2.zip
sudo ./aws/install
rm -r ./aws ./awscliv2.zip
popd

# if [[ -f "/path/to/file" ]]; then
#     export AWS_CONFIG_FILE=/custom/path/config
#     export AWS_SHARED_CREDENTIALS_FILE=/custom/path/credentials
# fi

echo "Install kubectl"
sudo apt-get install bash-completion
echo 'source <(kubectl completion bash)' >>~/.bashrc
echo 'alias k=kubectl' >>~/.bashrc
echo 'complete -o default -F __start_kubectl k' >>~/.bashrc

echo "Install k9s"
go install github.com/derailed/k9s@latest
# wget https://github.com/derailed/k9s/releases/latest/download/k9s_linux_amd64.deb && apt install ./k9s_linux_amd64.deb && rm k9s_linux_amd64.deb




echo "Setup complete!"


# sudo apt install plocate