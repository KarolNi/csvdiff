package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRunContext(t *testing.T) {
	t.Run("should find diff in happy path", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		{
			baseContent := []byte(`id,name,age,desc
0,tom,2,developer
2,ryan,20,qa
4,emin,40,pm

`)
			err := afero.WriteFile(fs, "/base.csv", baseContent, os.ModePerm)
			assert.NoError(t, err)
		}
		{
			deltaContent := []byte(`id,name,age,desc
0,tom,2,developer
1,caprio,3,developer
2,ryan,23,qa
`)
			err := afero.WriteFile(fs, "/delta.csv", deltaContent, os.ModePerm)
			assert.NoError(t, err)
		}

		ctx, err := NewContext(
			fs,
			digest.Positions{0},
			nil,
			digest.Positions{1, 2},
			nil,
			nil,
			digest.Positions{0, 1, 2},
			nil,
			"json",
			"/base.csv",
			"/delta.csv",
			',',
			false,
		)
		assert.NoError(t, err)

		outStream := &bytes.Buffer{}
		errStream := &bytes.Buffer{}

		err = runContext(ctx, outStream, errStream)
		expected := `{
  "Additions": [
    "1,caprio,3"
  ],
  "Modifications": [
    {
      "Original": "2,ryan,20",
      "Current": "2,ryan,23"
    }
  ],
  "Deletions": [
    "4,emin,40"
  ]
}`

		assert.NoError(t, err)
		assert.Equal(t, expected, outStream.String())

	})
}

func TestRunContext2(t *testing.T) {
	t.Run("should find diff in happy path", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		{
			baseContent := []byte(`id,name,age,desc
0,tom,2,developer
2,ryan,20,qa
4,emin,40,pm
`)
			err := afero.WriteFile(fs, "/base.csv", baseContent, os.ModePerm)
			assert.NoError(t, err)
		}
		{
			deltaContent := []byte(`rownr,id,name,desc,age
1,0,tom,developer,2
2,1,caprio,developer,3
3,2,ryan,qa,23
`)
			err := afero.WriteFile(fs, "/delta.csv", deltaContent, os.ModePerm)
			assert.NoError(t, err)
		}

		ctx, err := NewContext(
			fs,
			digest.Positions{0},
			digest.Positions{1},
			digest.Positions{2, 3},
			digest.Positions{4, 3},
			nil,
			digest.Positions{0, 1, 2},
			digest.Positions{1, 2, 4},
			"json",
			"/base.csv",
			"/delta.csv",
			',',
			false,
		)
		assert.NoError(t, err)

		outStream := &bytes.Buffer{}
		errStream := &bytes.Buffer{}

		err = runContext(ctx, outStream, errStream)
		expected := `{
  "Additions": [
    "1,caprio,3"
  ],
  "Modifications": [
    {
      "Original": "2,ryan,20",
      "Current": "2,ryan,23"
    }
  ],
  "Deletions": [
    "4,emin,40"
  ]
}`

		assert.NoError(t, err)
		assert.Equal(t, expected, outStream.String())

	})
}
