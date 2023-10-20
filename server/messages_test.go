package server_test

import (
	"fmt"
	. "github.com/CameronHonis/chess-arbitrator/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Messages", func() {
	var message *Message
	var messageJson []byte

	BeforeEach(func() {
		message = &Message{
			Topic: MESSAGE_TOPIC_AUTH,
			Content: &AuthMessageContent{
				PublicKey:  "some-public-key",
				PrivateKey: "some-private-key",
			},
		}
		messageJson = []byte(`{"topic":1,"content":{"publicKey":"some-public-key","privateKey":"some-private-key"}}`)
	})
	Describe("marshalling a message to JSON", func() {
		It("converts the Message into JSON", func() {
			realJsonBytes, err := message.Marshal()
			Expect(err).ToNot(HaveOccurred())
			Expect(realJsonBytes).To(Equal(messageJson))
		})
	})
	Describe("unmarshalling JSON to a message", func() {
		When("the json is missing a topic", func() {
			BeforeEach(func() {
				messageJson = []byte(`{"content":{"publicKey":"some-public-key","privateKey": "some-private-key"}}`)
			})
			It("returns an error", func() {
				_, err := UnmarshalToMessage(messageJson)
				Expect(err).To(Equal(fmt.Errorf("could not determine required field 'Topic' from %s while constructing Message", string(messageJson))))
			})
		})
		When("the json is missing content", func() {
			BeforeEach(func() {
				messageJson = []byte(`{"topic": 1}`)
			})
			It("returns an error", func() {
				_, err := UnmarshalToMessage(messageJson)
				Expect(err).To(Equal(fmt.Errorf("could not extract content map from %s while constructing Message content", string(messageJson))))
			})
		})
		When("when the message type is MESSAGE_TOPIC_AUTH", func() {
			BeforeEach(func() {
				message = &Message{
					Topic: MESSAGE_TOPIC_AUTH,
					Content: &AuthMessageContent{
						PublicKey:  "some-public-key",
						PrivateKey: "some-private-key",
					},
				}
				messageJson = []byte(`{"topic":1,"content":{"publicKey":"some-public-key","privateKey":"some-private-key"}}`)
			})
			It("returns a message with its Content as a AuthMessageContent", func() {
				realMessage, err := UnmarshalToMessage(messageJson)
				Expect(err).ToNot(HaveOccurred())
				Expect(realMessage).To(Equal(message))
			})
		})
	})
})
