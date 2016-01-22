package loggerplugin

import (
	"crypto/rand"
    "crypto/sha256"
    // TODO: log to our logger instead of system?
    "log"
)

const ID_LEN = 32

func id() string {
    b := make([]byte, ID_LEN)
    n, err := rand.Read(b)
    if err != nil {
        log.Fatal(err)
    }
    s := sha256.Sum256(b)
    return string(s[:n])
}
