package uuidx

import "uuid"

// UUIDv5 generates a deterministic UUIDv5 from a namespace and a name.
func UUIDv5(namespace uuid.UUID, name []byte) uuid.UUID {
	return NewSHA1(namespace, name)
}

// UUIDv5OID generates a deterministic UUIDv5 using the OID namespace.
func UUIDv5OID(name []byte) uuid.UUID {
	return UUIDv5(NameSpaceOID, name)
}
