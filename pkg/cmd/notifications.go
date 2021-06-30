package cmd

import (
	"fmt"

	"github.com/puppetlabs/relay/pkg/client/openapi"
	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/format"
	"github.com/spf13/cobra"
)

func newNotificationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notifications",
		Short: "Manage your notifications",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(newListUserNotificationsCommand())
	cmd.AddCommand(newClearUserNotificationsCommand())

	return cmd
}

func newListUserNotificationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List notifications",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doListUserNotifications,
	}

	return cmd
}

func newClearUserNotificationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear notifications",
		Args:  cobra.MaximumNArgs(1),
	}

	cmd.AddCommand(newClearAllReadUserNotificationsCommand())

	return cmd
}

func newClearAllReadUserNotificationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Clear all read notifications",
		Args:  cobra.MaximumNArgs(1),
		RunE:  doClearAllReadUserNotifications,
	}

	return cmd
}

func doListUserNotifications(cmd *cobra.Command, args []string) error {
	Dialog.Progress("Listing notifications...")

	req := Client.Api.NotificationsApi.GetNotifications(cmd.Context())
	n, _, err := Client.Api.NotificationsApi.GetNotificationsExecute(req)
	if err != nil {
		return errors.NewClientInternalError().WithCause(err)
	}

	if n.Notifications == nil || len(*n.Notifications) == 0 {
		return nil
	}

	t := Dialog.Table()

	t.Headers([]string{"Status", "Type", "Name", "Run Number", "Link"})

	for _, un := range *n.Notifications {
		wfn := un.GetFields()["workflow_name"].(string)
		rn := int64(un.GetFields()["run_number"].(float64))
		read := un.Read

		status := ""
		if !read {
			status = "NEW"
		}

		link := format.GuiLink(Config, "/workflows/%s/runs/%d/graph", wfn, rn)
		nt := ""
		switch un.Type {
		case "workflow.failed":
			nt = "Workflow failed"
		case "workflow.succeeded":
			nt = "Workflow succeeded"
		case "step.approval":
			nt = "Approval needed"
		}

		if nt != "" {
			t.AppendRow([]string{status, nt, wfn, fmt.Sprintf("%d", rn), link})
		}
	}

	t.Flush()

	return nil
}

func doClearAllReadUserNotifications(cmd *cobra.Command, args []string) error {
	Dialog.Progress("Clearing notifications...")

	req := Client.Api.NotificationsApi.GetNotifications(cmd.Context())
	n, _, err := Client.Api.NotificationsApi.GetNotificationsExecute(req)
	if err != nil {
		return errors.NewClientInternalError().WithCause(err)
	}

	nids := make([]string, 0)
	for _, un := range *n.Notifications {
		if un.Read {
			nids = append(nids, un.GetId())
		}
	}

	if len(nids) > 0 {
		req := Client.Api.NotificationsApi.PostAllNotificationDone(cmd.Context())
		_, _, err := Client.Api.NotificationsApi.PostAllNotificationDoneExecute(
			req.NotificationIdentifiers(openapi.NotificationIdentifiers{Ids: &nids}),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
