package conversation

import (
	"encoding/json"
	"github.com/messagebird/go-rest-api/v7/internal/mbtest"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnmarshalContact(t *testing.T) {
	data := mbtest.Testdata(t, "contact.json")

	c := &Contact{}
	err := json.Unmarshal(data, c)

	assert.NoError(t, err)
	assert.Equal(t, "9354647c5b144a2b4c99f2n42497249", c.ID)
	assert.Equal(t, "https://rest.messagebird.com/1/contacts/9354647c5b144a2b4c99f2n42497249", c.Href)
	assert.Equal(t, "316123456789", c.MSISDN)
	assert.Equal(t, "Jen", c.FirstName)
	assert.Equal(t, "Smith", c.LastName)
	assert.Equal(t, "2022-06-03T20:06:03Z", c.CreatedDatetime.Format(time.RFC3339))
	assert.Nil(t, c.UpdatedDatetime)
	assert.Equal(
		t,
		map[string]interface{}{
			"custom1": nil,
			"custom2": nil,
			"custom3": nil,
			"custom4": nil,
		},
		c.CustomDetails,
	)
}
