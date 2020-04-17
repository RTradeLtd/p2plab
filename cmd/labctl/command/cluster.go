// Copyright 2019 Netflix, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"errors"

	"github.com/Netflix/p2plab"
	"github.com/Netflix/p2plab/pkg/cliutil"
	"github.com/Netflix/p2plab/printer"
	"github.com/Netflix/p2plab/query"
	"github.com/rs/zerolog"
	"github.com/urfave/cli"
)

var clusterCommand = cli.Command{
	Name:    "cluster",
	Aliases: []string{"c"},
	Usage:   "Manage clusters.",
	Subcommands: []cli.Command{
		{
			Name:      "create",
			Aliases:   []string{"c"},
			Usage:     "Creates a new cluster.",
			ArgsUsage: "<name>",
			Action:    createClusterAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "definition,d",
					Usage: "Create cluster from a cluster definition.",
				},
				&cli.IntFlag{
					Name:  "size,s",
					Usage: "Size of cluster.",
					Value: 3,
				},
				&cli.StringFlag{
					Name:  "instance-type,t",
					Usage: "EC2 Instance type of cluster.",
					Value: "t2.micro",
				},
				&cli.StringFlag{
					Name:  "region,r",
					Usage: "AWS Region to deploy to.",
					Value: "us-west-2",
				},
			},
		},
		{
			Name:      "inspect",
			Aliases:   []string{"inspect"},
			Usage:     "Displays detailed information on a cluster.",
			ArgsUsage: "<name>",
			Action:    inspectClusterAction,
		},
		{
			Name:      "label",
			Aliases:   []string{"l"},
			Usage:     "Add or remove labels from clusters.",
			ArgsUsage: " ",
			Action:    labelClustersAction,
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:  "add",
					Usage: "Adds a label.",
				},
				&cli.StringSliceFlag{
					Name:  "remove,rm",
					Usage: "Removes a label.",
				},
			},
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "List clusters.",
			ArgsUsage: " ",
			Action:    listClusterAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "query,q",
					Usage: "Runs a query to filter the listed clusters.",
				},
			},
		},
		{
			Name:      "remove",
			ArgsUsage: "[<name> ...]",
			Aliases:   []string{"rm"},
			Usage:     "Remove clusters.",
			Action:    removeClustersAction,
		},
	},
}

func createClusterAction(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("cluster name must be provided")
	}

	p, err := CommandPrinter(c, printer.OutputID)
	if err != nil {
		return err
	}

	control, err := ResolveControl(c)
	if err != nil {
		return err
	}
	ctx := cliutil.CommandContext(c)

	var options []p2plab.CreateClusterOption
	if c.IsSet("definition") {
		options = append(options,
			p2plab.WithClusterDefinition(c.String("definition")),
		)
	} else {
		options = append(options,
			p2plab.WithClusterSize(c.Int("size")),
			p2plab.WithClusterInstanceType(c.String("instance-type")),
			p2plab.WithClusterRegion(c.String("region")),
		)
	}

	name := c.Args().First()
	id, err := control.Cluster().Create(ctx, name, options...)
	if err != nil {
		return err
	}

	cluster, err := control.Cluster().Get(ctx, id)
	if err != nil {
		return err
	}

	zerolog.Ctx(ctx).Info().Msgf("Created cluster %q", cluster.Metadata().ID)
	return p.Print(cluster.Metadata())
}

func inspectClusterAction(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.New("cluster name must be provided")
	}

	p, err := CommandPrinter(c, printer.OutputJSON)
	if err != nil {
		return err
	}

	control, err := ResolveControl(c)
	if err != nil {
		return err
	}

	ctx := cliutil.CommandContext(c)
	name := c.Args().First()
	cluster, err := control.Cluster().Get(ctx, name)
	if err != nil {
		return err
	}

	return p.Print(cluster.Metadata())
}

func labelClustersAction(c *cli.Context) error {
	var names []string
	for i := 0; i < c.NArg(); i++ {
		names = append(names, c.Args().Get(i))
	}

	p, err := CommandPrinter(c, printer.OutputTable)
	if err != nil {
		return err
	}

	control, err := ResolveControl(c)
	if err != nil {
		return err
	}

	ctx := cliutil.CommandContext(c)
	cs, err := control.Cluster().Label(ctx, names, c.StringSlice("add"), c.StringSlice("remove"))
	if err != nil {
		return err
	}

	l := make([]interface{}, len(cs))
	for i, c := range cs {
		l[i] = c.Metadata()
	}

	return p.Print(l)
}

func listClusterAction(c *cli.Context) error {
	control, err := ResolveControl(c)
	if err != nil {
		return err
	}

	p, err := CommandPrinter(c, printer.OutputTable)
	if err != nil {
		return err
	}

	var opts []p2plab.ListOption
	ctx := cliutil.CommandContext(c)
	if c.IsSet("query") {
		q, err := query.Parse(ctx, c.String("query"))
		if err != nil {
			return err
		}

		opts = append(opts, p2plab.WithQuery(q.String()))
	}

	cs, err := control.Cluster().List(ctx, opts...)
	if err != nil {
		return err
	}

	l := make([]interface{}, len(cs))
	for i, c := range cs {
		l[i] = c.Metadata()
	}

	return p.Print(l)
}

func removeClustersAction(c *cli.Context) error {
	var names []string
	for i := 0; i < c.NArg(); i++ {
		names = append(names, c.Args().Get(i))
	}

	control, err := ResolveControl(c)
	if err != nil {
		return err
	}

	ctx := cliutil.CommandContext(c)
	err = control.Cluster().Remove(ctx, c.Bool("force"), names...)
	if err != nil {
		return err
	}

	return nil
}
