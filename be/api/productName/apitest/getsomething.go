package apitest

import (
	"context"
	"net/http"
	"strings"

	"github.com/template/be/ctxs"
	"github.com/template/be/lib/language"
	"github.com/template/be/lib/router"
	"github.com/template/be/locale"
)

// Get all holidays
func (m *Module) getTestingData(ctx context.Context, param *router.HandlerParam) (apiResult router.HandlerResult) {

	// Validate language
	langExists := false
	availableLanguages := []string{"id-ID", "en-US"}
	for l := range availableLanguages {
		if availableLanguages[l] == param.Lang {
			langExists = true
			break
		}
	}

	// Language didn't exists
	if !langExists {
		apiResult.SetError(locale.Translate(language.ID, "ErrorServer"), "45", http.StatusBadRequest)
		return apiResult
	}

	// Validate country
	countryID, _ := ctxs.CountryIDFromContext(ctx)
	countryID = strings.ToUpper(countryID)

	// Process result
	resultData := router.ArrObject{}

	type testing struct {
		Label string
		Value string
	}

	testRes := testing{
		Label: "testing",
		Value: "testing",
	}
	test := router.Object{
		ID:         "testing",
		Attributes: testRes,
	}

	resultData = append(resultData, test)

	apiResult.JSON = &router.Data{
		Data: resultData,
	}

	return apiResult
}
