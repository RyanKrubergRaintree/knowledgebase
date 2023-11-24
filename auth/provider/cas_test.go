package provider

import (
	"os"
	"testing"
)

func TestGetCompanyIdRedirect(t *testing.T) {
	type args struct {
		company   string
		companyid string
	}
	tests := []struct {
		name  string
		args  args
		want string
	}{
		{
			name: "Should return original companyid",
			args: args{
				company:   "OÜ Viive sünnipäev",
				companyid: "123",
			},
			want: "123",
		},
		{
			name: "Should return redirected companyid",
			args: args{
				company:   "OÜ Juss Anna Ampsu",
				companyid: "123",
			},
			want: "1234",
		},
	}


	os.Setenv("OÜ_JUSS_ANNA_AMPSU", "1234")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotCompanyId := getCompanyIdRedirect(tt.args.company, tt.args.companyid)

			if gotCompanyId != tt.want {
				t.Errorf("getCompanyIdRedirect() got = %v, want %v", gotCompanyId, tt.want)
			}
			if gotCompanyId != tt.want {
				t.Errorf("getCompanyIdRedirect() got1 = %v, want %v", gotCompanyId, tt.want)
			}
		})
	}

	os.Unsetenv("OÜ_JUSS_ANNA_AMPSU")
}
