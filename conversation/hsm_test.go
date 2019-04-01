package conversation

import (
	"reflect"
	"testing"
	"time"
)

func TestLocalizableParameter(t *testing.T) {
	now := time.Now()

	tt := []struct {
		name   string
		got    *HSMLocalizableParameter
		expect *HSMLocalizableParameter
	}{
		{
			name: "default",
			got:  DefaultLocalizableHSMParameter("foo"),
			expect: &HSMLocalizableParameter{
				Default: "foo",
			},
		},
		{
			name: "currency",
			got:  CurrencyLocalizableHSMParameter("EUR 12.34", "EUR", int64(12340)),
			expect: &HSMLocalizableParameter{
				Default: "EUR 12.34",
				Currency: &HSMLocalizableParameterCurrency{
					Code:   "EUR",
					Amount: int64(12340),
				},
			},
		},
		{
			name: "date time",
			got:  DateTimeLocalizableHSMParameter("baz", now),
			expect: &HSMLocalizableParameter{
				Default:  "baz",
				DateTime: &now,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.got, tc.expect) {
				t.Fatalf("got %v, expected %v", tc.got, tc.expect)
			}
		})
	}
}
