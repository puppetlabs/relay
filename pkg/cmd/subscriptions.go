package cmd

import (
	"github.com/puppetlabs/relay/pkg/client/openapi"
	"github.com/spf13/cobra"
)

func newSubscriptionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "Manage your Relay subscriptions",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newListUserWorkflowSubscriptions())
	cmd.AddCommand(newSubscribeUserWorkflow())
	cmd.AddCommand(newUnsubscribeUserWorkflow())

	return cmd
}

func newListUserWorkflowSubscriptions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workflow subscriptions",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doListUserWorkflowSubscriptions,
	}

	return cmd
}

func newSubscribeUserWorkflow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribe [workflow name]",
		Short: "Subscribe to workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doSubscribeUserWorkflow,
	}

	return cmd
}

func newUnsubscribeUserWorkflow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unsubscribe [workflow name]",
		Short: "Unsubscribe to workflow",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doUnsubscribeUserWorkflow,
	}

	return cmd
}

func doListUserWorkflowSubscriptions(cmd *cobra.Command, args []string) error {
	Dialog.Progress("Listing workflow subscriptions...")

	req := Client.Api.SubscriptionsApi.GetWorkflowsSubscriptions(cmd.Context())
	uws, _, err := Client.Api.SubscriptionsApi.GetWorkflowsSubscriptionsExecute(req)
	if err != nil {
		return err
	}

	for _, wf := range *uws.Workflows {
		if wf.Subscriptions != nil &&
			*wf.Subscriptions.Subscribe {
			Dialog.Infof(wf.Name)
		}
	}

	return nil
}

func doSubscribeUserWorkflow(cmd *cobra.Command, args []string) error {
	name, err := getWorkflowName(args)

	if err != nil {
		return err
	}

	Dialog.Progress("Subscribing...")

	var subscribe = true
	req := Client.Api.SubscriptionsApi.PutWorkflowSubscriptions(cmd.Context(), name)
	_, _, cerr := Client.Api.SubscriptionsApi.PutWorkflowSubscriptionsExecute(
		req.UserWorkflowSubscriptions(
			openapi.UserWorkflowSubscriptions{Subscribe: &subscribe}))
	if cerr != nil {
		return cerr
	}

	return nil
}

func doUnsubscribeUserWorkflow(cmd *cobra.Command, args []string) error {
	name, err := getWorkflowName(args)

	if err != nil {
		return err
	}

	Dialog.Progress("Unsubscribing...")

	var subscribe = false
	req := Client.Api.SubscriptionsApi.PutWorkflowSubscriptions(cmd.Context(), name)
	_, _, cerr := Client.Api.SubscriptionsApi.PutWorkflowSubscriptionsExecute(
		req.UserWorkflowSubscriptions(
			openapi.UserWorkflowSubscriptions{Subscribe: &subscribe}))
	if cerr != nil {
		return cerr
	}

	return nil
}
