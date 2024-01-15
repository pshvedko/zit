package loop

import "testing"

type E struct {
	*E
}

func (e *E) Next() Element {
	if e.E != nil {
		return e.E
	}
	return nil
}

func TestLoop(t *testing.T) {
	type args struct {
		e Element
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				e: func() Element {
					return nil
				}(),
			},
			want: false,
		},
		{
			args: args{
				e: func() Element {
					return &E{}
				}(),
			},
			want: false,
		},
		{
			args: args{
				e: func() Element {
					e := &E{E: &E{E: &E{}}}
					return e
				}(),
			},
			want: false,
		},
		{
			args: args{
				e: func() Element {
					e := &E{}
					e.E = e
					return e
				}(),
			},
			want: true,
		},
		{
			args: args{
				e: func() Element {
					e := &E{E: &E{}}
					e.E.E = e
					return e
				}(),
			},
			want: true,
		},
		{
			args: args{
				e: func() Element {
					e := &E{E: &E{E: &E{}}}
					e.E.E.E = e
					return e
				}(),
			},
			want: true,
		},
		{
			args: args{
				e: func() Element {
					e := &E{E: &E{E: &E{}}}
					e.E.E.E = e.E
					return e
				}(),
			},
			want: true,
		},
		{
			args: args{
				e: func() Element {
					e := &E{E: &E{E: &E{}}}
					e.E.E.E = e.E.E
					return e
				}(),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Loop(tt.args.e); got != tt.want {
				t.Errorf("Loop() = %v, want %v", got, tt.want)
			}
		})
	}
}
