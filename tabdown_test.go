package tabdown

import "testing"

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
