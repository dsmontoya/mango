package mango

import (
	"context"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
)

func Test_portString(t *testing.T) {
	Convey("Given a valid port number", t, func() {
		var port uint = 3000
		Convey("When the port is parsed to string", func() {
			s := portString(port)
			Convey("The string should be valid", func() {
				So(s, ShouldEqual, ":3000")
			})
		})
	})

	Convey("Given a zero port number", t, func() {
		var port uint = 0
		Convey("When the port is parsed to string", func() {
			s := portString(port)
			Convey("The string should be empty", func() {
				So(s, ShouldBeEmpty)
			})
		})
	})
}

func Test_hostString(t *testing.T) {
	Convey("Given a address and port", t, func() {
		var port uint = 3000
		host := "example.com"

		Convey("When the host string is generated", func() {
			s := hostString(host, port)
			Convey("The string should be valid", func() {
				So(s, ShouldEqual, "example.com:3000")
			})
		})
	})
}

func TestConnect(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name         string
		args         args
		wantDatabase string
		wantErr      bool
	}{
		{
			"connect",
			args{
				&Config{
					ClientOptions: *new(mongooptions.ClientOptions).ApplyURI(os.Getenv("MONGO_URI")),
					Database:      "test",
					Context:       context.Background(),
				},
			},
			"test",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, err := Connect(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = connection.client.Ping(context.Background(), nil)
			if err != nil {
				t.Errorf("Unable to ping")
			}
		})
	}
}
