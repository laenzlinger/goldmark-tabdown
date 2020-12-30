package tabdown

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

func Test_isChordLine(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty line should not be a chord line",
			args: args{line: []byte("")},
			want: false,
		},
		{
			name: "should detect single chord",
			args: args{line: []byte(" [Am] ")},
			want: true,
		},
		{
			name: "should reject other stuff chord",
			args: args{line: []byte(" [Am] other stuff")},
			want: false,
		},
		{
			name: "should detect multiple chords",
			args: args{line: []byte(" [Am] [Cmaj] [] ")},
			want: true,
		},
		{
			name: "should not detect empty chords",
			args: args{line: []byte(" [] ")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isChordLine(tt.args.line); got != tt.want {
				t.Errorf("isChordLine(%v) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestMeta(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			Tabdown,
		),
	)
	source := `
# Let it be

       C              G                 Am     Am7  Fmaj7    F6
When I find myself in times of trouble, Mother Mary comes to me

`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	assert.Equal(t, `<h1>Let it be</h1>
`, buf.String())
}
