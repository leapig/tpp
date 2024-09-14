package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	ValidateSignatureError int = -40001
	ParseXmlError          int = -40002
	ComputeSignatureError  int = -40003
	IllegalAesKey          int = -40004
	ValidateCorpidError    int = -40005
	EncryptAESError        int = -40006
	DecryptAESError        int = -40007
	IllegalBuffer          int = -40008
	EncodeBase64Error      int = -40009
	DecodeBase64Error      int = -40010
	GenXmlError            int = -40010
	ParseJsonError         int = -40012
	GenJsonError           int = -40013
	IllegalProtocolType    int = -40014
)

type ProtocolType int

const (
	XmlType ProtocolType = 1
)

type CryptError struct {
	ErrCode int
	ErrMsg  string
}

func NewCryptError(err_code int, err_msg string) *CryptError {
	return &CryptError{ErrCode: err_code, ErrMsg: err_msg}
}

type WXBizMsg4Recv struct {
	Tousername string `xml:"ToUserName"`
	Encrypt    string `xml:"Encrypt"`
	Agentid    string `xml:"AgentID"`
}

type CDATA struct {
	Value string `xml:",cdata"`
}

type WXBizMsg4Send struct {
	XMLName   xml.Name `xml:"xml"`
	Encrypt   CDATA    `xml:"Encrypt"`
	Signature CDATA    `xml:"MsgSignature"`
	Timestamp string   `xml:"TimeStamp"`
	Nonce     CDATA    `xml:"Nonce"`
}

func NewWXBizMsg4Send(encrypt, signature, timestamp, nonce string) *WXBizMsg4Send {
	return &WXBizMsg4Send{Encrypt: CDATA{Value: encrypt}, Signature: CDATA{Value: signature}, Timestamp: timestamp, Nonce: CDATA{Value: nonce}}
}

type ProtocolProcessor interface {
	parse(src_data []byte) (*WXBizMsg4Recv, *CryptError)
	serialize(msg_send *WXBizMsg4Send) ([]byte, *CryptError)
}

type WXBizMsgCrypt struct {
	token              string
	encoding_aeskey    string
	receiver_id        string
	protocol_processor ProtocolProcessor
}

type XmlProcessor struct {
}

func (self *XmlProcessor) parse(src_data []byte) (*WXBizMsg4Recv, *CryptError) {
	var msg4_recv WXBizMsg4Recv
	err := xml.Unmarshal(src_data, &msg4_recv)
	if nil != err {
		return nil, NewCryptError(ParseXmlError, "xml to msg fail")
	}
	return &msg4_recv, nil
}

func (self *XmlProcessor) serialize(msg4_send *WXBizMsg4Send) ([]byte, *CryptError) {
	xml_msg, err := xml.Marshal(msg4_send)
	if nil != err {
		return nil, NewCryptError(GenXmlError, err.Error())
	}
	return xml_msg, nil
}

func NewWXBizMsgCrypt(receiver_id string) *WXBizMsgCrypt {
	return &WXBizMsgCrypt{token: "91BxDcWOdJ1PUBngaOkgsH9B", encoding_aeskey: "yOocubSvv527UCbHOL2TkTPygLUcbBAgSaZgtKvZlAs=", receiver_id: receiver_id, protocol_processor: new(XmlProcessor)}
}

func (self *WXBizMsgCrypt) pKCS7Padding(plaintext string, block_size int) []byte {
	padding := block_size - (len(plaintext) % block_size)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	var buffer bytes.Buffer
	buffer.WriteString(plaintext)
	buffer.Write(padtext)
	return buffer.Bytes()
}

func (self *WXBizMsgCrypt) pKCS7Unpadding(plaintext []byte, block_size int) ([]byte, *CryptError) {
	plaintext_len := len(plaintext)
	if nil == plaintext || plaintext_len == 0 {
		return nil, NewCryptError(DecryptAESError, "pKCS7Unpadding error nil or zero")
	}
	if plaintext_len%block_size != 0 {
		return nil, NewCryptError(DecryptAESError, "pKCS7Unpadding text not a multiple of the block size")
	}
	padding_len := int(plaintext[plaintext_len-1])
	return plaintext[:plaintext_len-padding_len], nil
}

func (self *WXBizMsgCrypt) cbcEncrypter(plaintext string) ([]byte, *CryptError) {
	aeskey, err := base64.StdEncoding.DecodeString(self.encoding_aeskey)
	if nil != err {
		return nil, NewCryptError(DecodeBase64Error, err.Error())
	}
	const block_size = 32
	pad_msg := self.pKCS7Padding(plaintext, block_size)

	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, NewCryptError(EncryptAESError, err.Error())
	}

	ciphertext := make([]byte, len(pad_msg))
	iv := aeskey[:aes.BlockSize]

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(ciphertext, pad_msg)
	base64_msg := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(base64_msg, ciphertext)

	return base64_msg, nil
}

func (self *WXBizMsgCrypt) cbcDecrypter(base64_encrypt_msg string) ([]byte, *CryptError) {
	aeskey, err := base64.StdEncoding.DecodeString(self.encoding_aeskey)
	if nil != err {
		return nil, NewCryptError(DecodeBase64Error, err.Error())
	}

	encrypt_msg, err := base64.StdEncoding.DecodeString(base64_encrypt_msg)
	if nil != err {
		return nil, NewCryptError(DecodeBase64Error, err.Error())
	}

	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, NewCryptError(DecryptAESError, err.Error())
	}

	if len(encrypt_msg) < aes.BlockSize {
		return nil, NewCryptError(DecryptAESError, "encrypt_msg size is not valid")
	}

	iv := aeskey[:aes.BlockSize]

	if len(encrypt_msg)%aes.BlockSize != 0 {
		return nil, NewCryptError(DecryptAESError, "encrypt_msg not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(encrypt_msg, encrypt_msg)

	return encrypt_msg, nil
}

func (self *WXBizMsgCrypt) calSignature(timestamp, nonce, data string) string {
	sort_arr := []string{self.token, timestamp, nonce, data}
	sort.Strings(sort_arr)
	var buffer bytes.Buffer
	for _, value := range sort_arr {
		buffer.WriteString(value)
	}

	sha := sha1.New()
	sha.Write(buffer.Bytes())
	signature := fmt.Sprintf("%x", sha.Sum(nil))
	return string(signature)
}

func (self *WXBizMsgCrypt) ParsePlainText(plaintext []byte) ([]byte, uint32, []byte, []byte, *CryptError) {
	const block_size = 32
	plaintext, err := self.pKCS7Unpadding(plaintext, block_size)
	if nil != err {
		return nil, 0, nil, nil, err
	}

	text_len := uint32(len(plaintext))
	if text_len < 20 {
		return nil, 0, nil, nil, NewCryptError(IllegalBuffer, "plain is to small 1")
	}
	random := plaintext[:16]
	msg_len := binary.BigEndian.Uint32(plaintext[16:20])
	if text_len < (20 + msg_len) {
		return nil, 0, nil, nil, NewCryptError(IllegalBuffer, "plain is to small 2")
	}

	msg := plaintext[20 : 20+msg_len]
	receiver_id := plaintext[20+msg_len:]

	return random, msg_len, msg, receiver_id, nil
}

func (self *WXBizMsgCrypt) VerifyURL(msg_signature, timestamp, nonce, echostr string) ([]byte, string, *CryptError) {
	signature := self.calSignature(timestamp, nonce, echostr)

	if strings.Compare(signature, msg_signature) != 0 {
		return nil, self.receiver_id, NewCryptError(ValidateSignatureError, "signature not equal")
	}

	plaintext, err := self.cbcDecrypter(echostr)
	if nil != err {
		return nil, self.receiver_id, err
	}

	_, _, msg, receiver_id, err := self.ParsePlainText(plaintext)
	if nil != err {
		return nil, self.receiver_id, err
	}

	if len(self.receiver_id) > 0 && strings.Compare(string(receiver_id), self.receiver_id) != 0 {
		fmt.Println(string(receiver_id), self.receiver_id, len(receiver_id), len(self.receiver_id))
		return nil, string(receiver_id), NewCryptError(ValidateCorpidError, "receiver_id is not equil")
	}

	return msg, string(receiver_id), nil
}

func (self *WXBizMsgCrypt) EncryptMsg(reply_msg, timestamp, nonce string) ([]byte, *CryptError) {
	rand_str := randString(16)
	var buffer bytes.Buffer
	buffer.WriteString(rand_str)

	msg_len_buf := make([]byte, 4)
	binary.BigEndian.PutUint32(msg_len_buf, uint32(len(reply_msg)))
	buffer.Write(msg_len_buf)
	buffer.WriteString(reply_msg)
	buffer.WriteString(self.receiver_id)

	tmp_ciphertext, err := self.cbcEncrypter(buffer.String())
	if nil != err {
		return nil, err
	}
	ciphertext := string(tmp_ciphertext)

	signature := self.calSignature(timestamp, nonce, ciphertext)

	msg4_send := NewWXBizMsg4Send(ciphertext, signature, timestamp, nonce)
	return self.protocol_processor.serialize(msg4_send)
}

func (self *WXBizMsgCrypt) DecryptMsg(msg_signature, timestamp, nonce string, post_data []byte) ([]byte, string, *CryptError) {
	msg4_recv, crypt_err := self.protocol_processor.parse(post_data)
	if nil != crypt_err {
		return nil, self.receiver_id, crypt_err
	}

	signature := self.calSignature(timestamp, nonce, msg4_recv.Encrypt)

	if strings.Compare(signature, msg_signature) != 0 {
		return nil, self.receiver_id, NewCryptError(ValidateSignatureError, "signature not equal")
	}

	plaintext, crypt_err := self.cbcDecrypter(msg4_recv.Encrypt)
	if nil != crypt_err {
		return nil, self.receiver_id, crypt_err
	}

	_, _, msg, receiver_id, crypt_err := self.ParsePlainText(plaintext)
	if nil != crypt_err {
		return nil, string(receiver_id), crypt_err
	}

	if len(self.receiver_id) > 0 && strings.Compare(string(receiver_id), self.receiver_id) != 0 {
		return nil, string(receiver_id), NewCryptError(ValidateCorpidError, "receiver_id is not equil")
	}

	return msg, string(receiver_id), nil
}

type LarkCrypto struct {
	Token          string
	EncodingAESKey string
	BKey           []byte
	Block          cipher.Block
}

func NewLarkCrypto() *LarkCrypto {
	token := "91BxDcWOdJ1PUBngaOkgsH9B"
	encodingAESKey := "yOocubSvv527UCbHOL2TkTPygLUcbBAgSaZgtKvZlAs"
	c := &LarkCrypto{
		Token:          token,
		EncodingAESKey: encodingAESKey,
	}
	return c
}

func (self *LarkCrypto) GetDecryptMsg(signature, timestamp, nonce string, msg []byte) (string, error) {
	// 签名验证
	if signature != "" {
		if strings.Compare(signature, self.callSignature(timestamp, nonce, string(msg))) != 0 {
			return "", errors.New("signature not equal")
		}
	}
	var msgJson map[string]interface{}
	_ = json.Unmarshal(msg, &msgJson)
	buf, err := base64.StdEncoding.DecodeString(msgJson["encrypt"].(string))
	if err != nil {
		return "", fmt.Errorf("base64StdEncode Error[%v]", err)
	}
	if len(buf) < aes.BlockSize {
		return "", errors.New("cipher  too short")
	}
	keyBs := sha256.Sum256([]byte(self.EncodingAESKey))
	block, err := aes.NewCipher(keyBs[:sha256.Size])
	if err != nil {
		return "", fmt.Errorf("AESNewCipher Error[%v]", err)
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]
	if len(buf)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	// token验证
	i := strings.Index(string(buf), self.Token)
	if i == -1 {
		return "", errors.New("token not equal")
	}
	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}
	return string(buf[n : m+1]), nil
}

func (self *LarkCrypto) callSignature(timestamp, nonce, msg string) string {
	var b strings.Builder
	b.WriteString(timestamp)
	b.WriteString(nonce)
	b.WriteString(self.EncodingAESKey)
	b.WriteString(msg)
	bs := []byte(b.String())
	h := sha256.New()
	h.Write(bs)
	bs = h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

type DingTalkCrypto struct {
	Token          string
	EncodingAESKey string
	SuiteKey       string
	BKey           []byte
	Block          cipher.Block
}

func NewDingTalkCrypto(token, encodingAESKey, suiteKey string) *DingTalkCrypto {
	if len(encodingAESKey) != 43 {
		panic("不合法的EncodingAESKey")
	}
	bkey, err := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	if err != nil {
		panic(err.Error())
	}
	block, err := aes.NewCipher(bkey)
	if err != nil {
		panic(err.Error())
	}
	c := &DingTalkCrypto{
		Token:          token,
		EncodingAESKey: encodingAESKey,
		SuiteKey:       suiteKey,
		BKey:           bkey,
		Block:          block,
	}
	return c
}

func (c *DingTalkCrypto) GetDecryptMsg(signature, timestamp, nonce, secretMsg string) (string, error) {
	if c.calSignature(c.Token, timestamp, nonce, secretMsg) != signature {
		return "", errors.New("ERROR: 签名不匹配")
	}
	decode, err := base64.StdEncoding.DecodeString(secretMsg)
	if err != nil {
		return "", err
	}
	if len(decode) < aes.BlockSize {
		return "", errors.New("ERROR: 密文太短")
	}
	blockMode := cipher.NewCBCDecrypter(c.Block, c.BKey[:c.Block.BlockSize()])
	plantText := make([]byte, len(decode))
	blockMode.CryptBlocks(plantText, decode)
	plantText = pkCS7UnPadding(plantText)
	size := binary.BigEndian.Uint32(plantText[16:20])
	plantText = plantText[20:]
	corpID := plantText[size:]
	if string(corpID) != c.SuiteKey {
		return "", errors.New("ERROR: CorpID匹配不正确")
	}
	return string(plantText[:size]), nil
}

func (c *DingTalkCrypto) GetEncryptMsg(msg, timestamp, nonce string) (string, string, error) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(msg)))
	msg = randString(16) + string(size) + msg + c.SuiteKey
	plantText := pkCS7Padding([]byte(msg), c.Block.BlockSize())
	if len(plantText)%aes.BlockSize != 0 {
		return "", "", errors.New("ERROR: 消息体size不为16的倍数")
	}
	blockMode := cipher.NewCBCEncrypter(c.Block, c.BKey[:c.Block.BlockSize()])
	chipherText := make([]byte, len(plantText))
	blockMode.CryptBlocks(chipherText, plantText)
	outMsg := base64.StdEncoding.EncodeToString(chipherText)
	signature := c.calSignature(c.Token, timestamp, nonce, outMsg)
	return outMsg, signature, nil
}

// 数据签名
func (c *DingTalkCrypto) calSignature(token, timestamp, nonce, msg string) string {
	params := make([]string, 0)
	params = append(params, token)
	params = append(params, timestamp)
	params = append(params, nonce)
	params = append(params, msg)
	sort.Strings(params)
	sha := sha1.New()
	sha.Write([]byte(strings.Join(params, "")))
	signature := fmt.Sprintf("%x", sha.Sum(nil))
	return string(signature)
}

// 解密补位
func pkCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

// 加密补位
func pkCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 随机字符串
func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GetRandStringWithCharset 获取指定字符集 下 指定长度的随机字符串
func GetRandStringWithCharset(length int, charset string) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GetRandString 获取指定长度的随机字符串
func GetRandString(length int) string {
	return GetRandStringWithCharset(length, charset)
}
