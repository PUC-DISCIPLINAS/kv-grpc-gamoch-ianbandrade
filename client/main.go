package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"key_value/client/pb"
)

const (
	address = "localhost:50051"
)

func getClient(ctx context.Context) pb.KeyValueServiceClient {
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Did not connect - %s\n", err)
		os.Exit(1)
	}

	return pb.NewKeyValueServiceClient(conn)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cmdPut = &cobra.Command{
		Use:   "put [key] [value]",
		Short: "Put a entry",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = getClient(ctx).Put(ctx, &pb.KeyValue{Key: args[0], Value: args[1]})
		},
	}

	var cmdGet = &cobra.Command{
		Use:   "get [key]",
		Short: "Get value of key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			getResponse, _ := getClient(ctx).Get(ctx, &pb.Key{Key: args[0]})

			if getResponse.Defined {
				fmt.Printf("Get: Key='%s' Value='%s'\n", args[0], getResponse.GetValue())
			} else {
				fmt.Printf("Get: Key not exists\n")
			}
		},
	}

	var cmdGetAllKeys = &cobra.Command{
		Use:   "getAllKeys",
		Short: "Get all keys",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			getAllKeysResponse, _ := getClient(ctx).GetAllKeys(ctx, &pb.Empty{})
			fmt.Printf("GetAllKeys: %s\n", strings.Join(getAllKeysResponse.GetKeys(), ", "))
		},
	}

	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(cmdPut, cmdGet, cmdGetAllKeys)
	_ = rootCmd.Execute()
}
