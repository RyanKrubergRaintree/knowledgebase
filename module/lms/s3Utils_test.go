package lms

import (
	"testing"

	"github.com/raintreeinc/knowledgebase/kb"
)


func TestGetS3Client(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "Should return s3 client which should have default region",
			want: kb.DefaultRegion,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getS3Client()

			if (err != nil) != tt.wantErr {
				t.Errorf("getS3Client() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got.Config.Region != tt.want {
				t.Errorf("getS3Client() = %v, want default region ", *got.Config.Region)
			}
		})
	}
}
