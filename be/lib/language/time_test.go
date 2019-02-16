package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLongDay(t *testing.T) {
	assert.Equal(t, "Minggu", getLongDay(ID, "Sunday"))
	assert.Equal(t, "Ahad", getLongDay(MY, "Sunday"))
	assert.Equal(t, "星期日", getLongDay(CN, "Sunday"))
	assert.Equal(t, "ABC", getLongDay("", "ABC"))
}

func TestGetShortDay(t *testing.T) {
	assert.Equal(t, "Min", getShortDay(ID, "Sun"))
	assert.Equal(t, "Ahd", getShortDay(MY, "Sun"))
	assert.Equal(t, "星期日", getShortDay(CN, "Sun"))
	assert.Equal(t, "ABC", getShortDay("", "ABC"))
}

func TestGetLongMonth(t *testing.T) {

	assert.Equal(t, "Januari", getLongMonth(ID, "January"))
	assert.Equal(t, "Januari", getLongMonth(MY, "January"))
	assert.Equal(t, "一月", getLongMonth(CN, "January"))
	assert.Equal(t, "ABC", getLongMonth("", "ABC"))
}

func TestGetShortMonth(t *testing.T) {

	assert.Equal(t, "Jan", getShortMonth(ID, "Jan"))
	assert.Equal(t, "Jan", getShortMonth(MY, "Jan"))
	assert.Equal(t, "一月", getShortMonth(CN, "Jan"))
	assert.Equal(t, "ABC", getShortMonth("", "ABC"))
}

func TestTimeTranslationEmpty(t *testing.T) {
	assert.Equal(t, "", timeTranslate("foo", "bar"))
}
