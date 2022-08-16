# `x/nftstaking`

## Abstract

`nftstaking` module is provide rewards to nft stakers on different networks.

The owner registers nft stakers on the module.

Part of inflation rewards are allocated to registered nft stakers and broadcasted on everyblock.

## State

### NftStaking

`nftstaking` module keeps the information of `NftStaking` that shows nft information and reward address with weights.

```protobuf
enum NftType {
  NFT_TYPE_DEFAULT = 0 [ (gogoproto.enumvalue_customname) = "NftTypeDefault" ];
}
message NftStaking {
  NftType nft_type = 1;
  string nft_identifier = 2;
  string nft_metadata = 3;
  string reward_address = 4;
  uint64 reward_weight = 5;
}
```

### NftTypePerms

`NftTypePerms` defines allocated permissions per nft type. Module admin can set this and manage permissions.

```protobuf
enum Permission {
  SET_SERVER_ACCESS = 0 [ (gogoproto.enumvalue_customname) = "PermSetServerAccess" ];
}

message NftTypePerms {
  NftType nft_type = 1;
  repeated Permission perms = 2;
}
```

### Access

`ServerAccess` defines server name and accessible channels on the server.
`Access` defines servers access info for an address on-chain.

```protobuf
message ServerAccess {
  string server = 1;
  repeated string channels = 2;
}
message Access {
  string address = 1;
  repeated ServerAccess servers = 2 [ (gogoproto.nullable) = false ];
}
```

## Messages

### MsgRegisterNftStaking

`MsgRegisterNftStaking` describes the message to register staked nfts by admin.

```go
type MsgRegisterNftStaking struct {
	Sender     string
	NftStaking NftStaking
}
```

### MsgSetAccessInfo

`MsgSetAccessInfo` describes the message to set access info by permissioned nft holder or by module admin.

```protobuf
message MsgSetAccessInfo {
  string sender = 1;
  teritori.nftstaking.v1beta1.Access access_info = 2 [ (gogoproto.nullable) = false ];
}
```

### MsgSetNftTypePerms

`MsgSetAccessInfo` to set required permissions to a nft type by module admin.

```protobuf
message MsgSetNftTypePerms {
  string sender = 1;
  teritori.nftstaking.v1beta1.NftTypePerms nft_type_perms = 2 [ (gogoproto.nullable) = false ];
}
```
