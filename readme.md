### simple das-account-indexer server

#### build

`make rpc-win`

`make rpc-mac`

`make rpc-linux`

#### run

```
cd bin/mac
./rpc_server -config="local_server.yaml"
```

#### request

```curl
curl --location --request POST 'http://localhost:8111' \
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

#### resp

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
                "status": 0,
                "owner_lock_script": {
                    "code_hash": "0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8",
                    "hash_type": "type",
                    "args": "IK87TtHHdoqLh9L8JyQsHDpD1F8="
                },
                "manager_lock_script": {
                    "code_hash": "0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8",
                    "hash_type": "type",
                    "args": "IK87TtHHdoqLh9L8JyQsHDpD1F8="
                },
                "records": [
                    {
                        "key": "",
                        "type": "this_is_type",
                        "label": "eth",
                        "value": "0x12233",
                        "ttl": "5"
                    },
                    {
                        "key": "",
                        "type": "this_is_type",
                        "label": "btc",
                        "value": "mmmxxx",
                        "ttl": "5"
                    },
                    {
                        "key": "",
                        "type": "this_is_type",
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