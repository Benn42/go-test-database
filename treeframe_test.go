package gotestdatabase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeFrame(t *testing.T) {
	db := NewTestDatabase([]string{"test"})

	db.Insert("test", []*DatabaseRecord{
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
		{
			"key":  "value",
			"test": true,
		},
	})

	styler := NewUnicodeStyler()
	frame := &TreeFrame{
		styler: styler,
		lines:  make([]string, 0),
	}

	expected := `8
├─<4
│  ├─<2
│  │  ├─<1
│  │  └─>3
│  └─>6
│     ├─<5
│     └─>7
└─>12
   ├─<10
   │  ├─<9
   │  └─>11
   └─>13
      └─>14`

	renderable := frame.render(db.tables[0], false)
	assert.Equal(t, expected, renderable)
}
