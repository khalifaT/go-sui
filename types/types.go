package types

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"
)

type Address = HexData

/**
 * Creates Address from a hex string.
 * @param addr Hex string can be with a prefix or without a prefix,
 *   e.g. '0x1aa' or '1aa'. Hex string will be left padded with 0s if too short.
 */
func NewAddressFromHex(addr string) (*Address, error) {
	if strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X") {
		addr = addr[2:]
	}
	if len(addr)%2 != 0 {
		addr = "0" + addr
	}

	bytes, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	const addressLength = 20
	if len(bytes) > addressLength {
		return nil, fmt.Errorf("Hex string is too long. Address's length is %v bytes.", addressLength)
	}

	res := [addressLength]byte{}
	copy(res[addressLength-len(bytes):], bytes[:])
	return &Address{
		data: res[:],
	}, nil
}

// Returns the address with leading zeros trimmed, e.g. 0x2
func (a Address) ShortString() string {
	return "0x" + strings.TrimLeft(hex.EncodeToString(a.data), "0")
}

type ObjectId = HexData
type Digest = Base64Data

type InputObjectKind map[string]interface{}

type TransactionBytes struct {
	// the gas object to be used
	Gas ObjectRef `json:"gas"`

	// objects to be used in this transaction
	InputObjects []InputObjectKind `json:"inputObjects"`

	// transaction data bytes
	TxBytes Base64Data `json:"txBytes"`
}

type ObjectRef struct {
	Digest   Digest   `json:"digest"`
	ObjectId ObjectId `json:"objectId"`
	Version  int      `json:"version"`
}

type SignatureScheme string

const (
	SignatureSchemeEd25519   SignatureScheme = "ED25519"
	SignatureSchemeSecp256k1 SignatureScheme = "Secp256k1"
)

type SignedTransaction struct {
	// transaction data bytes
	TxBytes *Base64Data `json:"tx_bytes"`

	// Flag of the signature scheme that is used.
	SigScheme SignatureScheme `json:"sig_scheme"`

	// transaction signature
	Signature *Base64Data `json:"signature"`

	// signer's public key
	PublicKey *Base64Data `json:"pub_key"`
	Mode      string      `json:"mode"`
}

type TransferObject struct {
	Recipient Address   `json:"recipient"`
	ObjectRef ObjectRef `json:"object_ref"`
}
type ModulePublish struct {
	Modules [][]byte `json:"modules"`
}
type MoveCall struct {
	Package  ObjectRef     `json:"package"`
	Module   string        `json:"module"`
	Function string        `json:"function"`
	TypeArgs []interface{} `json:"type_arguments"`
	Args     []interface{} `json:"arguments"`
}
type TransferSui struct {
	Recipient Address `json:"recipient"`
	Amount    uint64  `json:"amount"`
}
type ChangeEpoch struct {
	Epoch             interface{} `json:"epoch"`
	StorageCharge     uint64      `json:"storage_charge"`
	ComputationCharge uint64      `json:"computation_charge"`
}

type SingleTransactionKind struct {
	TransferObject *TransferObject `json:"TransferObject,omitempty"`
	Publish        *ModulePublish  `json:"Publish,omitempty"`
	Call           *MoveCall       `json:"Call,omitempty"`
	TransferSui    *TransferSui    `json:"TransferSui,omitempty"`
	ChangeEpoch    *ChangeEpoch    `json:"ChangeEpoch,omitempty"`
}

type SenderSignedData struct {
	Transactions []SingleTransactionKind `json:"transactions,omitempty"`

	Sender     *Address   `json:"sender"`
	GasPayment *ObjectRef `json:"gasPayment"`
	GasBudget  uint64     `json:"gasBudget"`
	// GasPrice     uint64      `json:"gasPrice"`
}

type TransactionEnvelope struct {
	TransactionDigest *Digest `json:"transactionDigest"`

	Data *SenderSignedData `json:"data"`

	TxSignature *Base64Data `json:"txSignature"`

	AuthSignInfo interface{} `json:"authSignInfo"`
}

type CertifiedTransaction = TransactionEnvelope

type OwnedObjectRef struct {
	Owner     *ObjectOwner `json:"owner"`
	Reference *ObjectRef   `json:"reference"`
}

type Event interface{}

type GasCostSummary struct {
	ComputationCost uint64 `json:"computationCost"`
	StorageCost     uint64 `json:"storageCost"`
	StorageRebate   uint64 `json:"storageRebate"`
}

const (
	TransactionStatusSuccess = "success"
	TransactionStatusFailure = "failure"
)

type TransactionStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type TransactionEffects struct {
	Status TransactionStatus `json:"status"`

	TransactionDigest *Base64Data     `json:"transactionDigest"`
	GasUsed           *GasCostSummary `json:"gasUsed"`
	GasObject         *OwnedObjectRef `json:"gasObject"`
	Events            []Event         `json:"events,omitempty"`
	Dependencies      []Digest        `json:"dependencies,omitempty"`

	// SharedObjects []ObjectRef      `json:"sharedObjects"`
	Created   []OwnedObjectRef `json:"created,omitempty"`
	Mutated   []OwnedObjectRef `json:"mutated,omitempty"`
	Unwrapped []OwnedObjectRef `json:"unwrapped,omitempty"`
	Deleted   []ObjectRef      `json:"deleted,omitempty"`
	Wrapped   []ObjectRef      `json:"wrapped,omitempty"`
}

func (te *TransactionEffects) GasFee() uint64 {
	return te.GasUsed.StorageCost - te.GasUsed.StorageRebate + te.GasUsed.ComputationCost
}

type ParsedTransactionResponse interface{}
type TransactionResponseV2 struct {
	EffectsCert TransactionResponse `json:"effectsCert"`
}

type TransactionResponse struct {
	Certificate *CertifiedTransaction     `json:"certificate"`
	Effects     *TransactionEffects       `json:"effects"`
	ParsedData  ParsedTransactionResponse `json:"parsed_data,omitempty"`
	TimestampMs uint64                    `json:"timestamp_ms,omitempty"`
}

type ObjectOwner struct {
	AddressOwner *Address `json:"AddressOwner,omitempty"`
	ObjectOwner  *Address `json:"ObjectOwner,omitempty"`
	SingleOwner  *Address `json:"SingleOwner,omitempty"`
}

type ObjectReadDetail struct {
	Data  map[string]interface{} `json:"data"`
	Owner *ObjectOwner           `json:"owner"`

	PreviousTransaction *Digest    `json:"previousTransaction"`
	StorageRebate       int        `json:"storageRebate"`
	Reference           *ObjectRef `json:"reference"`
}

type ObjectStatus string

const (
	ObjectStatusExists    ObjectStatus = "Exists"
	ObjectStatusNotExists ObjectStatus = "NotExists"
	ObjectStatusDeleted   ObjectStatus = "Deleted"
)

type ObjectRead struct {
	Details *ObjectReadDetail `json:"details"`
	Status  ObjectStatus      `json:"status"`
}

type ObjectInfo struct {
	ObjectId *ObjectId    `json:"objectId"`
	Version  int          `json:"version"`
	Digest   *Digest      `json:"digest"`
	Type     string       `json:"type"`
	Owner    *ObjectOwner `json:"owner"`

	PreviousTransaction *Digest `json:"previousTransaction"`
}

func (txn *TransactionBytes) SignWith(privateKey ed25519.PrivateKey) *SignedTransaction {
	message := txn.TxBytes.Data()
	signature := ed25519.Sign(privateKey, message)
	sign := Bytes(signature).GetBase64Data()
	publicKey := privateKey.Public().(ed25519.PublicKey)
	pub := Bytes(publicKey).GetBase64Data()

	return &SignedTransaction{
		TxBytes:   &txn.TxBytes,
		SigScheme: SignatureSchemeEd25519,
		Signature: &sign,
		PublicKey: &pub,
	}
}
