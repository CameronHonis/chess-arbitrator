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
			Topic:       "auth",
			ContentType: CONTENT_TYPE_AUTH,
			Content: &AuthMessageContent{
				PublicKey:  "some-public-key",
				PrivateKey: "some-private-key",
			},
		}
		messageJson = []byte(`{"senderKey":"","privateKey":"","topic":"auth","contentType":"AUTH","content":{"publicKey":"some-public-key","privateKey":"some-private-key"}}`)
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
				messageJson = []byte(`{"contentType": "AUTH", "content":{"publicKey":"some-public-key","privateKey": "some-private-key"}}`)
			})
			It("does not return an error", func() {
				_, err := UnmarshalToMessage(messageJson)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		When("the json is missing a content type", func() {
			BeforeEach(func() {
				messageJson = []byte(`{"topic": "auth", "content":{"publicKey":"some-public-key","privateKey": "some-private-key"}}`)
			})
			It("returns an error", func() {
				_, err := UnmarshalToMessage(messageJson)
				expErr := fmt.Errorf("could not extract content type from %s while constructing Message content", string(messageJson))
				Expect(err).To(Equal(expErr))
			})
		})
		When("the json is missing content", func() {
			BeforeEach(func() {
				messageJson = []byte(`{"topic": "auth", "contentType": "AUTH"}`)
			})
			It("returns an error", func() {
				_, err := UnmarshalToMessage(messageJson)
				expErr := fmt.Errorf("could not extract content map from %s while constructing Message content", string(messageJson))
				Expect(err).To(Equal(expErr))
			})
		})
		When("when the message type is MESSAGE_TOPIC_AUTH", func() {
			BeforeEach(func() {
				message = &Message{
					Topic:       "auth",
					ContentType: CONTENT_TYPE_AUTH,
					Content: &AuthMessageContent{
						PublicKey:  "some-public-key",
						PrivateKey: "some-private-key",
					},
				}
				messageJson = []byte(`{"topic": "auth", "contentType": "AUTH", "content":{"publicKey":"some-public-key","privateKey":"some-private-key"}}`)
			})
			It("returns a message with its Content as a AuthMessageContent", func() {
				realMessage, err := UnmarshalToMessage(messageJson)
				Expect(err).ToNot(HaveOccurred())
				Expect(realMessage).To(Equal(message))
			})
		})
	})
})
