package webpush

import "crypto/ecdsa"

func vapidAuthorization(endpoint string, privateKey *ecdsa.PrivateKey, publicKey []byte, subject string) (string, error) {
	
}
