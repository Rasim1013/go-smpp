package smpp

import (
	"fmt"
	"log"
	"time"
	"unicode"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

// Определение кодировки текста
func needsUCS2(text string) bool {
	for _, r := range text {
		if r > 127 || unicode.Is(unicode.Cyrillic, r) {
			return true
		}
	}
	return false
}

// Отправка одного SMS
func SendSMS(server, username, password, sender, receiver, message string) error {
	transceiver := &smpp.Transceiver{
		Addr:        server,
		User:        username,
		Passwd:      password,
		RespTimeout: 2 * time.Second,
	}

	// Кодировка текста
	var pdutextCodec pdutext.Codec
	if needsUCS2(message) {
		pdutextCodec = pdutext.UCS2(message)
	} else {
		pdutextCodec = pdutext.GSM7(message)
	}

	// Устанавливаем соединение
	statusChan := transceiver.Bind()
	go func() {
		for status := range statusChan {
			if status.Error() != nil {
				log.Printf("Connection error: %v", status.Error())
			} else {
				log.Printf("Connection status: %+v", status.Status())
			}
		}
	}()

        time.Sleep(100 * time.Millisecond) // 100 мс между запросами

	_, err := transceiver.SubmitLongMsg(&smpp.ShortMessage{
		Src:           sender,
		Dst:           receiver,
		Text:          pdutextCodec,
		Register:      1,
		SourceAddrTON: 5,
		SourceAddrNPI: 0,
		DestAddrTON:   1,
		DestAddrNPI:   1,
	})

	transceiver.Close()

	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}
