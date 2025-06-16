package util

import (
	"testing"
)

func TestLogln(t *testing.T) {
	log := NewLog(0)
	log.Println(0, "LogLn test")
	// Output: LogLn test
}

func TestLogf(t *testing.T) {
	log := NewLog(0)
	log.Printf(0, "Logf test: int %d - str %s\n", 42, "test")
	// Output: Logf test: int 42 - str test
}

func TestIf(t *testing.T) {
	type args struct {
		condition bool
		a         any
		b         any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{"TestIf", args{true, "a", "b"}, "a"},
		{"TestIf", args{false, "a", "b"}, "b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := If(tt.args.condition, tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("If() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat2string(t *testing.T) {
	type args struct {
		f    float64
		prec int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"TestFloat2string", args{1.23456789, 2}, "1.23"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float2string(tt.args.f, tt.args.prec); got != tt.want {
				t.Errorf("Float2string() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWarningErrorCheck(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"TestWarningErrorCheck", args{nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WarningErrorCheck(tt.args.err); got != tt.want {
				t.Errorf("WarningErrorCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFatalErrorCheck(t *testing.T) {
	// no puedo probar con un error porque log.Fatal() termina el programa
	FatalErrorCheck(nil)
}

/*
func TestGetHttpContent(t *testing.T) {
	type args struct {
		urlz string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"TestGetHttpContent", args{"https://www.google.com"}, false},
		{"TestGetHttpContent", args{"https://www.google.com/404"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetHttpContent(tt.args.urlz, []Header{}, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHttpContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
*/

func TestStartStopTor(t *testing.T) {
	StartTor()
	defer StopTor()
}

/*
func TestRetry(t *testing.T) {
	url := "http://checkip.dyndns.org/"

	err := Retry(10, 1*time.Second, func() error {
		body, err := GetHttpContent(url, []Header{}, false)
		if err == nil {
			// remove html tags from body
			// The policy can then be used to sanitize lots of input and it is safe to use the policy in multiple goroutines
			body := bluemonday.StripTagsPolicy().SanitizeBytes(body)

			Logln(string(body))
		}
		return err
	})
	if err != nil {
		t.Errorf("Retry() error = %v", err)
	}
}
*/

// test executableDir()
func TestExecutableDir(t *testing.T) {
	ExecutableDir()
}

// test GetHttpContent()
func TestGetHttpContent(t *testing.T) {
	_, err := GetHttpContent("https://ifconfig.me/ip")
	if err != nil {
		t.Errorf("GetHttpContent() error = %v", err)
	}
}
