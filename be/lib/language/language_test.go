package language

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDefault(t *testing.T) {
	assert.Equal(t, ID, GetDefault("XX"))
	assert.Equal(t, EN, GetDefault("en-us"))
}

func TestGetLangID(t *testing.T) {
	assert.Equal(t, -1, GetLangID(""))
	assert.Equal(t, 0, GetLangID(ID))
	assert.Equal(t, 1, GetLangID(EN))
	assert.Equal(t, 2, GetLangID(MY))
	assert.Equal(t, 3, GetLangID(CN))
}

func TestGetLangTypeByID(t *testing.T) {
	assert.Equal(t, "", GetLangTypeByID(-1))
	assert.Equal(t, ID, GetLangTypeByID(0))
	assert.Equal(t, EN, GetLangTypeByID(1))
	assert.Equal(t, MY, GetLangTypeByID(2))
	assert.Equal(t, CN, GetLangTypeByID(3))
}

func TestDate(t *testing.T) {

	d := time.Date(1990, 10, 11, 12, 12, 12, 12, time.UTC)
	assert.Equal(t, "Thu 11 Thursday Oct October", Date(EN, "Mon 02 Monday Jan January", d))
	assert.Equal(t, "Kam 11 Kamis Okt Oktober", Date(ID, "Mon 02 Monday Jan January", d))
	assert.Equal(t, "Kha 11 Khamis Okt Oktober", Date(MY, "Mon 02 Monday Jan January", d))
	assert.Equal(t, "星期四 11 星期四 十月 十月", Date(CN, "Mon 02 Monday Jan January", d))
}

func TestGenerateDuration(t *testing.T) {
	assert.Equal(t, "1j 40m", Duration(ID, 100, true))
	assert.Equal(t, "17h 30m", Duration(EN, 1050, true))
	assert.Equal(t, "1h 9j 20m", Duration(ID, 2000, true))
	assert.Equal(t, "30 menit", Duration(ID, 30, false))
	assert.Equal(t, "1 jam 1 menit", Duration(ID, 61, false))
	assert.Equal(t, "5 jam 45 menit", Duration(ID, 345, false))
	assert.Equal(t, "5 days 8 hours 2 minutes", Duration(EN, 7682, false))
	assert.Equal(t, "5 days 8 hours", Duration(EN, 7680, false))
	assert.Equal(t, "1 hari 12 jam", Duration(ID, 2160, false))
}
