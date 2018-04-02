// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"fmt"
	"strconv"

	"github.com/meyskens/justexpenseit/expensify"
	"github.com/spf13/cobra"
)

type totalCmdOptions struct {
	From  string
	To    string
	Limit float64
}

// NewTotalCmd generates the `total` command
func NewTotalCmd() *cobra.Command {
	options := totalCmdOptions{}
	c := &cobra.Command{
		Use:   "total",
		Short: "get the total amount you've expensed during a certain period",
		Long:  `get the total amount you've expensed during a certain period`,
		Example: `
Get the total expenses made between 2018-03-01 and 2018-04-01
# justexpenseit total --from 2018-03-01 --to 2018-04-01
		`,
		RunE: options.RunE,
	}

	c.Flags().StringVar(&options.From, "from", "", "The start date for looking up expenses (eg. 2018-03-01)")
	c.Flags().StringVar(&options.To, "to", "", "The end date for looking up expenses (eg. 2018-04-01)")
	c.Flags().Float64Var(&options.Limit, "limit", 0, "This sets a limit you can expense so we can tell you how much you have left")

	c.MarkFlagRequired("from")
	c.MarkFlagRequired("to")

	return c
}

func (t *totalCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	api := expensify.New(partnerUserID, partnerUserSecret)

	expenses, err := api.GetExpenses("OPEN,SUBMITTED,APPROVED,REIMBURSED,ARCHIVED", t.From, t.To, 100)
	if err != nil {
		return err
	}

	var total float64
	for _, expense := range expenses {
		total += float64(expense.Amount) / 100.0
	}

	fmt.Printf("You have expensed a total of %s\n", strconv.FormatFloat(total, 'f', 2, 64))
	if t.Limit != 0 {
		fmt.Printf("You have %s left\n", strconv.FormatFloat(t.Limit-total, 'f', 2, 64))
	}

	return nil
}
