package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/sboon-gg/svctl/svctl"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type stopOpts struct {
	*serverOpts
}

func newStopOpts() *stopOpts {
	return &stopOpts{
		serverOpts: newServerOpts(),
	}
}

func stopCmd() *cobra.Command {
	opts := newStopOpts()

	cmd := &cobra.Command{
		Use:          "stop",
		Short:        "Stops the server",
		Long:         `Stops the server`,
		SilenceUsage: true,
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *stopOpts) AddFlags(cmd *cobra.Command) {
	o.serverOpts.AddFlags(cmd)
}

func (o *stopOpts) Run(cmd *cobra.Command, args []string) error {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server at localhost:50051: %v", err)
	}
	defer conn.Close()
	c := svctl.NewServersClient(conn)

	ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
	defer cancel()

	path, err := o.Path()
	if err != nil {
		return err
	}

	r, err := c.Stop(ctx, &svctl.ServerOpts{Path: path})
	if err != nil {
		return fmt.Errorf("error calling function Stop: %v", err)
	}

	cmd.Printf("Server status: %v\n", r.GetStatus().String())
	return nil
}

// func (o *stopOpts) Run(cmd *cobra.Command, args []string) error {
// 	sv, err := o.Server()
// 	if err != nil {
// 		return err
// 	}
//
// 	cache, err := sv.Cache()
// 	if err != nil {
// 		return err
// 	}
//
// 	var proc *prbf2.Proc
// 	if cache.PID < 0 {
// 		return fmt.Errorf("server is not running")
// 	}
//
// 	proc, err = prbf2.OpenProc(cache.PID)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = proc.Kill()
// 	if err != nil {
// 		return err
// 	}
//
// 	cache.PID = -1
// 	err = sv.WriteCache(cache)
// 	if err != nil {
// 		return err
// 	}
//
// 	cmd.Println("Server stopped")
// 	return nil
// }

func init() {
	rootCmd.AddCommand(stopCmd())
}
