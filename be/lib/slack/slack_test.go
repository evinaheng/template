package slack_test

import (
	"errors"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	. "github.com/template/be/lib/slack"
)

func TestSend(t *testing.T) {

	defer gock.Off()
	gock.New("https://www.dummy.com").Post("").MatchType("application/x-www-form-urlencoded").Reply(200).BodyString(`ok`)

	assert.Nil(t, Send("Unit testing for slack. Just ignore this message, my fellow friends", "https://www.dummy.com", true))
}
func TestSendError(t *testing.T) {

	defer gock.Off()
	gock.New("https://www.dummy.com").Post("").MatchType("application/x-www-form-urlencoded").ReplyError(errors.New("foo"))

	assert.EqualError(t, Send("Unit testing for slack. Just ignore this message, my fellow friends", "https://www.dummy.com", true), "Post https://www.dummy.com: foo")
}

func TestSendAttachment(t *testing.T) {

	defer gock.Off()
	gock.New("https://www.dummy.com").Post("").MatchType("application/x-www-form-urlencoded").Reply(200).BodyString(`ok`)

	a := Attachment{}

	assert.Nil(t, SendAttachment(a, "https://www.dummy.com", true))
}

func TestSendAttachmentError(t *testing.T) {

	defer gock.Off()
	gock.New("https://www.dummy.com").Post("").MatchType("application/x-www-form-urlencoded").ReplyError(errors.New("foo"))

	a := Attachment{}

	assert.EqualError(t, SendAttachment(a, "https://www.dummy.com", true), "Post https://www.dummy.com: foo")
}
