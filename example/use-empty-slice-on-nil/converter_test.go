package emptyonnil_test

import (
	"encoding/json"
	"testing"

	emptyonnil "github.com/jmattheis/goverter/example/use-empty-slice-on-nil"
	"github.com/jmattheis/goverter/example/use-empty-slice-on-nil/generated"
	"github.com/stretchr/testify/require"
)

func TestEmptyOnNil(t *testing.T) {
	var c emptyonnil.Converter = &generated.ConverterImpl{}

	items := c.Convert(nil)
	require.NotNil(t, items, "nil slice must be converted to an empty (non-nil) slice")
	require.Empty(t, items)

	raw, err := json.Marshal(items)
	require.NoError(t, err)
	require.JSONEq(t, `[]`, string(raw))

	m := c.ConvertMap(nil)
	require.NotNil(t, m, "nil map must be converted to an empty (non-nil) map")
	require.Empty(t, m)

	raw, err = json.Marshal(m)
	require.NoError(t, err)
	require.JSONEq(t, `{}`, string(raw))

	// non-nil input still converts normally
	actual := c.Convert([]emptyonnil.Input{{Name: "jmattheis"}})
	require.Equal(t, []emptyonnil.Output{{Name: "jmattheis"}}, actual)
}

func TestStrictConverterKeepsNil(t *testing.T) {
	var c emptyonnil.StrictConverter = &generated.StrictConverterImpl{}

	require.Nil(t, c.Convert(nil), "without the setting nil must stay nil (JSON null)")

	raw, err := json.Marshal(c.Convert(nil))
	require.NoError(t, err)
	require.Equal(t, `null`, string(raw))

	require.Nil(t, c.ConvertMap(nil))

	raw, err = json.Marshal(c.ConvertMap(nil))
	require.NoError(t, err)
	require.Equal(t, `null`, string(raw))
}
