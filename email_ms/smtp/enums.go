package smtp

import "bytes"

type messageType uint

const (
	mGreet messageType = iota
	mHELO
	mEHLO
	mTLS
	mSize
	mPipeline
	mEnhance
	mHelp
	mOK
	mReady
	mErrTransaction
	mErrTooBig
	mErrInvalidEmail
	mErrNoRecipients
	mErrTooManyRecipients
)

const (
	_         = iota             // ignore first value by assigning to blank identifier
	KB uint32 = 1 << (10 * iota) // 1 << (10*1)
	MB                           // 1 << (10*2)
)

type command []byte

var (
	cmdHELO     command = []byte("HELO")
	cmdEHLO     command = []byte("EHLO")
	cmdHELP     command = []byte("HELP")
	cmdMAIL     command = []byte("MAIL FROM:")
	cmdRCPT     command = []byte("RCPT TO:")
	cmdRSET     command = []byte("RSET")
	cmdVRFY     command = []byte("VRFY")
	cmdNOOP     command = []byte("NOOP")
	cmdQUIT     command = []byte("QUIT")
	cmdDATA     command = []byte("DATA")
	cmdSTARTTLS command = []byte("STARTTLS")
)

func (c command) match(in []byte) bool {
	lenIn := len(in)
	lenC := len(c)
	if lenIn < lenC {
		return false
	}
	return bytes.Equal(in[:lenC], c)
}
