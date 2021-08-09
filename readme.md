### simple das-account-indexer server

A tool service provided by DAS official, which can be used to query the information of successfully registered accounts on the chain.

More about DAS information, please visit DAS official website: <a target="_blank" href="https://da.systems/">https://da.systems/ </a>

### Quick start

> suggest use ubuntu.

* OS Support: Linux, Mac OS and Windows
* Go 1.12.x or later
* Git
* RocksDB, suggest use version 6.6.4, other version maybe cause some compile error

#### build

1. according your system to install the RocksDB, here is the document: https://github.com/facebook/rocksdb/blob/master/INSTALL.md;

2. clone the project to your `$GOPATH/src/`; 

3. `cd $GOPATH/src/das_account_indexer` then make:
    * `make rpc-win`
    * `make rpc-mac`
    * `make rpc-linux`
4. `*.h` files missing, if this error happen, see `eth-1.9.14` dir for solution;
5. other errors which your can't handle, contact me.

#### run
* start cmd:
    * `./rpc_server --config="local_server.yaml" --net_type=3`
* parameter:
    * `config`,configuration file's path;
        * demo: `conf/local_server.yaml`
    * `net_type`,server's net type. 1 means mainnet,2 means das-test2, 3 means das-test3
* support execute model(see the configuration file for details):
    1. search from chain data by GetLiveCell;
    2. local parsing block data to rocksdb and consume.
       

#### searchAccount

> search an account's info

```curl
curl --location --request POST 'http://localhost:8222' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "das_searchAccount",
    "params": [
        "linguanhong.bit"
    ]
}'
```

##### resp

* `owner_lock_chain_type`：
    * ETH
    * CKB
    * TRX
    * BTC

```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "errno": 0,
        "errmsg": "",
        "data": {
            "out_point": {
                "tx_hash": "0x62c068f066e46d53031d4aa170e7800c01b47bfd8b7f79d9d8095a2cecc23b15",
                "index": 0
            },
            "account_data": {
                "account": "linguanhong.bit",
                "account_id_hex": "0xb0e9b753b2853a464029",
                "next_account_id_hex": "0xbc4338222b62f10cec94",
                "create_at_unix": 1616489428,
                "expired_at_unix": 1679637114,
                "owner_address_chain": "ETH",
                "owner_address": "0x84ee75fd91a48c9d045840dc369fe22e045ff50a",
                "owner_lock_args_hex": "0x84ee75fd91a48c9d045840dc369fe22e045ff50a",
                "manager_address_chain": "ETH",
                "manager_address": "0x84ee75fd91a48c9d045840dc369fe22e045ff50a",
                "manager_lock_args_hex": "0x84ee75fd91a48c9d045840dc369fe22e045ff50a",
                "status": 0,
                "records": [
                    {
                        "key": "",
                        "label": "eth",
                        "value": "0x12233",
                        "ttl": "5"
                    },
                    {
                        "key": "",
                        "label": "btc",
                        "value": "mmmxxx",
                        "ttl": "5"
                    },
                    {
                        "key": "",
                        "label": "etc",
                        "value": "12222",
                        "ttl": "5"
                    }
                ]
            }
        }
    }
}
```

#### getAddressAccount

> find an address's accounts

```curl
curl --location --request POST 'http://localhost:8222' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "das_getAddressAccount",
    "params": [
        "ckt1qyqf4ehj9aaufevk5etpyt8k34pgctpgkapsdqjp6j"
    ]
}'
```

##### resp

```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "errno": 0,
        "errmsg": "",
        "data": [{object same as searchAccount}]
    }
}
```

### ErrCode

```go
const (
	Err_CallIndexer            DAS_CODE = 20000  // call ckb_indexer server error
	Err_Internal               DAS_CODE = 20001  // internal handle error
	Err_AccountExpired         DAS_CODE = 20002  
	Err_AccountFrozen          DAS_CODE = 20003
	Err_AccountAlreadyRegister DAS_CODE = 20004
	Err_AccountRecordsInvalid  DAS_CODE = 20005
	Err_AccountFormatInvalid   DAS_CODE = 20006
	Err_AccountNotExist        DAS_CODE = 20007
	Err_PubkeyHexFormatInvalid DAS_CODE = 20008
	Err_BaseParamInvalid       DAS_CODE = 20009
)
```

### Search time optimization

deploy this server in the `ckb_node` server.    
