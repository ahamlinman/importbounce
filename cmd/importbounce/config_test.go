package main

import (
	"reflect"
	"testing"
)

func TestFindPackage(t *testing.T) {
	conf := &config{
		Packages: []packageConfig{
			{
				Prefix:   "go.alexhamlin.co/importbounce",
				Import:   "git https://github.com/ahamlinman/importbounce",
				Redirect: "https://github.com/ahamlinman/importbounce",
			},
		},
	}

	testCases := []struct {
		path string
		want packageConfig
	}{
		{
			path: "go.alexhamlin.co/importbounce",
			want: conf.Packages[0],
		},

		{
			path: "go.alexhamlin.co/importbounce/cmd/importbounce",
			want: conf.Packages[0],
		},

		{
			path: "go.alexhamlin.co",
			want: packageConfig{},
		},

		{
			path: "go.alexhamlin.co/importbouncer",
			want: packageConfig{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			got := conf.FindPackage(tc.path)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("FindPackage(%s) = %v; want %v", tc.path, got, tc.want)
			}
		})
	}
}
