# lgh
# fix eth `*.h` files missing when run `go mod tidy | vendor`
cp -r $(pwd)/libsecp256k1 $(pwd)/../vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/