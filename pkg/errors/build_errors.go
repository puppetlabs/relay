// Package errors contains errors for the domain "neb".
//
// This file is automatically generated by errawr-gen. Do not modify it.
package errors

import (
	errawrgo "github.com/puppetlabs/errawr-go"
	impl "github.com/puppetlabs/errawr-go/impl"
)

// Error is the type of all errors generated by this package.
type Error interface {
	errawrgo.Error
}

// External contains methods that can be used externally to help consume errors from this package.
type External struct{}

// API is a singleton instance of the External type.
var API External

// Domain is the general domain in which all errors in this package belong.
var Domain = &impl.ErrorDomain{
	Key:   "neb",
	Title: "Nebula",
}

// GcpSection defines a section of errors with the following scope:
// GCP related errors
var GcpSection = &impl.ErrorSection{
	Key:   "gcp",
	Title: "GCP related errors",
}

// GcpClientCreateErrorCode is the code for an instance of "client_create_error".
const GcpClientCreateErrorCode = "neb_gcp_client_create_error"

// IsGcpClientCreateError tests whether a given error is an instance of "client_create_error".
func IsGcpClientCreateError(err errawrgo.Error) bool {
	return err != nil && err.Is(GcpClientCreateErrorCode)
}

// IsGcpClientCreateError tests whether a given error is an instance of "client_create_error".
func (External) IsGcpClientCreateError(err errawrgo.Error) bool {
	return IsGcpClientCreateError(err)
}

// GcpClientCreateErrorBuilder is a builder for "client_create_error" errors.
type GcpClientCreateErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "client_create_error" from this builder.
func (b *GcpClientCreateErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "an error occurred while creating the GCP client",
		Technical: "an error occurred while creating the GCP client",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "client_create_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     GcpSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "Client create error",
		Version:          1,
	}
}

// NewGcpClientCreateErrorBuilder creates a new error builder for the code "client_create_error".
func NewGcpClientCreateErrorBuilder() *GcpClientCreateErrorBuilder {
	return &GcpClientCreateErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewGcpClientCreateError creates a new error with the code "client_create_error".
func NewGcpClientCreateError() Error {
	return NewGcpClientCreateErrorBuilder().Build()
}

// GcpClusterDoesNotExistCode is the code for an instance of "cluster_does_not_exist".
const GcpClusterDoesNotExistCode = "neb_gcp_cluster_does_not_exist"

// IsGcpClusterDoesNotExist tests whether a given error is an instance of "cluster_does_not_exist".
func IsGcpClusterDoesNotExist(err errawrgo.Error) bool {
	return err != nil && err.Is(GcpClusterDoesNotExistCode)
}

// IsGcpClusterDoesNotExist tests whether a given error is an instance of "cluster_does_not_exist".
func (External) IsGcpClusterDoesNotExist(err errawrgo.Error) bool {
	return IsGcpClusterDoesNotExist(err)
}

// GcpClusterDoesNotExistBuilder is a builder for "cluster_does_not_exist" errors.
type GcpClusterDoesNotExistBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "cluster_does_not_exist" from this builder.
func (b *GcpClusterDoesNotExistBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "the remote cluster does not exist",
		Technical: "the remote cluster does not exist",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "cluster_does_not_exist",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     GcpSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "Cluster does not exist",
		Version:          1,
	}
}

// NewGcpClusterDoesNotExistBuilder creates a new error builder for the code "cluster_does_not_exist".
func NewGcpClusterDoesNotExistBuilder() *GcpClusterDoesNotExistBuilder {
	return &GcpClusterDoesNotExistBuilder{arguments: impl.ErrorArguments{}}
}

// NewGcpClusterDoesNotExist creates a new error with the code "cluster_does_not_exist".
func NewGcpClusterDoesNotExist() Error {
	return NewGcpClusterDoesNotExistBuilder().Build()
}

// GcpClusterReadErrorCode is the code for an instance of "cluster_read_error".
const GcpClusterReadErrorCode = "neb_gcp_cluster_read_error"

// IsGcpClusterReadError tests whether a given error is an instance of "cluster_read_error".
func IsGcpClusterReadError(err errawrgo.Error) bool {
	return err != nil && err.Is(GcpClusterReadErrorCode)
}

// IsGcpClusterReadError tests whether a given error is an instance of "cluster_read_error".
func (External) IsGcpClusterReadError(err errawrgo.Error) bool {
	return IsGcpClusterReadError(err)
}

// GcpClusterReadErrorBuilder is a builder for "cluster_read_error" errors.
type GcpClusterReadErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "cluster_read_error" from this builder.
func (b *GcpClusterReadErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "an error occurred while fetching the remote cluster",
		Technical: "an error occurred while fetching the remote cluster",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "cluster_read_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     GcpSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "Cluster read error",
		Version:          1,
	}
}

// NewGcpClusterReadErrorBuilder creates a new error builder for the code "cluster_read_error".
func NewGcpClusterReadErrorBuilder() *GcpClusterReadErrorBuilder {
	return &GcpClusterReadErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewGcpClusterReadError creates a new error with the code "cluster_read_error".
func NewGcpClusterReadError() Error {
	return NewGcpClusterReadErrorBuilder().Build()
}

// GcpClusterSyncErrorCode is the code for an instance of "cluster_sync_error".
const GcpClusterSyncErrorCode = "neb_gcp_cluster_sync_error"

// IsGcpClusterSyncError tests whether a given error is an instance of "cluster_sync_error".
func IsGcpClusterSyncError(err errawrgo.Error) bool {
	return err != nil && err.Is(GcpClusterSyncErrorCode)
}

// IsGcpClusterSyncError tests whether a given error is an instance of "cluster_sync_error".
func (External) IsGcpClusterSyncError(err errawrgo.Error) bool {
	return IsGcpClusterSyncError(err)
}

// GcpClusterSyncErrorBuilder is a builder for "cluster_sync_error" errors.
type GcpClusterSyncErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "cluster_sync_error" from this builder.
func (b *GcpClusterSyncErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "an error occurred while syncing the cluster",
		Technical: "an error occurred while syncing the cluster",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "cluster_sync_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     GcpSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "Cluster sync error",
		Version:          1,
	}
}

// NewGcpClusterSyncErrorBuilder creates a new error builder for the code "cluster_sync_error".
func NewGcpClusterSyncErrorBuilder() *GcpClusterSyncErrorBuilder {
	return &GcpClusterSyncErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewGcpClusterSyncError creates a new error with the code "cluster_sync_error".
func NewGcpClusterSyncError() Error {
	return NewGcpClusterSyncErrorBuilder().Build()
}

// WorkflowSection defines a section of errors with the following scope:
// Workflow errors
var WorkflowSection = &impl.ErrorSection{
	Key:   "workflow",
	Title: "Workflow errors",
}

// WorkflowActionDecodeErrorCode is the code for an instance of "action_decode_error".
const WorkflowActionDecodeErrorCode = "neb_workflow_action_decode_error"

// IsWorkflowActionDecodeError tests whether a given error is an instance of "action_decode_error".
func IsWorkflowActionDecodeError(err errawrgo.Error) bool {
	return err != nil && err.Is(WorkflowActionDecodeErrorCode)
}

// IsWorkflowActionDecodeError tests whether a given error is an instance of "action_decode_error".
func (External) IsWorkflowActionDecodeError(err errawrgo.Error) bool {
	return IsWorkflowActionDecodeError(err)
}

// WorkflowActionDecodeErrorBuilder is a builder for "action_decode_error" errors.
type WorkflowActionDecodeErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "action_decode_error" from this builder.
func (b *WorkflowActionDecodeErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "could not decode action configuration: {{reason}}",
		Technical: "could not decode action configuration: {{reason}}",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "action_decode_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     WorkflowSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "Action decode error",
		Version:          1,
	}
}

// NewWorkflowActionDecodeErrorBuilder creates a new error builder for the code "action_decode_error".
func NewWorkflowActionDecodeErrorBuilder(reason string) *WorkflowActionDecodeErrorBuilder {
	return &WorkflowActionDecodeErrorBuilder{arguments: impl.ErrorArguments{"reason": impl.NewErrorArgument(reason, "the reason there was a decoding error")}}
}

// NewWorkflowActionDecodeError creates a new error with the code "action_decode_error".
func NewWorkflowActionDecodeError(reason string) Error {
	return NewWorkflowActionDecodeErrorBuilder(reason).Build()
}

// WorkflowFileNotFoundCode is the code for an instance of "file_not_found".
const WorkflowFileNotFoundCode = "neb_workflow_file_not_found"

// IsWorkflowFileNotFound tests whether a given error is an instance of "file_not_found".
func IsWorkflowFileNotFound(err errawrgo.Error) bool {
	return err != nil && err.Is(WorkflowFileNotFoundCode)
}

// IsWorkflowFileNotFound tests whether a given error is an instance of "file_not_found".
func (External) IsWorkflowFileNotFound(err errawrgo.Error) bool {
	return IsWorkflowFileNotFound(err)
}

// WorkflowFileNotFoundBuilder is a builder for "file_not_found" errors.
type WorkflowFileNotFoundBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "file_not_found" from this builder.
func (b *WorkflowFileNotFoundBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "filepath `{{path}}` does not exist",
		Technical: "filepath `{{path}}` does not exist",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "file_not_found",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     WorkflowSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "File not found",
		Version:          1,
	}
}

// NewWorkflowFileNotFoundBuilder creates a new error builder for the code "file_not_found".
func NewWorkflowFileNotFoundBuilder(path string) *WorkflowFileNotFoundBuilder {
	return &WorkflowFileNotFoundBuilder{arguments: impl.ErrorArguments{"path": impl.NewErrorArgument(path, "the path that doesn't exist")}}
}

// NewWorkflowFileNotFound creates a new error with the code "file_not_found".
func NewWorkflowFileNotFound(path string) Error {
	return NewWorkflowFileNotFoundBuilder(path).Build()
}

// WorkflowLoaderErrorCode is the code for an instance of "loader_error".
const WorkflowLoaderErrorCode = "neb_workflow_loader_error"

// IsWorkflowLoaderError tests whether a given error is an instance of "loader_error".
func IsWorkflowLoaderError(err errawrgo.Error) bool {
	return err != nil && err.Is(WorkflowLoaderErrorCode)
}

// IsWorkflowLoaderError tests whether a given error is an instance of "loader_error".
func (External) IsWorkflowLoaderError(err errawrgo.Error) bool {
	return IsWorkflowLoaderError(err)
}

// WorkflowLoaderErrorBuilder is a builder for "loader_error" errors.
type WorkflowLoaderErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "loader_error" from this builder.
func (b *WorkflowLoaderErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "an error occurred while loading the workflow",
		Technical: "an error occurred while loading the workflow",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "loader_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     WorkflowSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "Loader error",
		Version:          1,
	}
}

// NewWorkflowLoaderErrorBuilder creates a new error builder for the code "loader_error".
func NewWorkflowLoaderErrorBuilder() *WorkflowLoaderErrorBuilder {
	return &WorkflowLoaderErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewWorkflowLoaderError creates a new error with the code "loader_error".
func NewWorkflowLoaderError() Error {
	return NewWorkflowLoaderErrorBuilder().Build()
}

// WorkflowNoCommandToExecuteErrorCode is the code for an instance of "no_command_to_execute_error".
const WorkflowNoCommandToExecuteErrorCode = "neb_workflow_no_command_to_execute_error"

// IsWorkflowNoCommandToExecuteError tests whether a given error is an instance of "no_command_to_execute_error".
func IsWorkflowNoCommandToExecuteError(err errawrgo.Error) bool {
	return err != nil && err.Is(WorkflowNoCommandToExecuteErrorCode)
}

// IsWorkflowNoCommandToExecuteError tests whether a given error is an instance of "no_command_to_execute_error".
func (External) IsWorkflowNoCommandToExecuteError(err errawrgo.Error) bool {
	return IsWorkflowNoCommandToExecuteError(err)
}

// WorkflowNoCommandToExecuteErrorBuilder is a builder for "no_command_to_execute_error" errors.
type WorkflowNoCommandToExecuteErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "no_command_to_execute_error" from this builder.
func (b *WorkflowNoCommandToExecuteErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "A valid command was not found to execute",
		Technical: "A valid command was not found to execute",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "no_command_to_execute_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     WorkflowSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "No command available to execute",
		Version:          1,
	}
}

// NewWorkflowNoCommandToExecuteErrorBuilder creates a new error builder for the code "no_command_to_execute_error".
func NewWorkflowNoCommandToExecuteErrorBuilder() *WorkflowNoCommandToExecuteErrorBuilder {
	return &WorkflowNoCommandToExecuteErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewWorkflowNoCommandToExecuteError creates a new error with the code "no_command_to_execute_error".
func NewWorkflowNoCommandToExecuteError() Error {
	return NewWorkflowNoCommandToExecuteErrorBuilder().Build()
}

// WorkflowNonExistentActionErrorCode is the code for an instance of "non_existent_action_error".
const WorkflowNonExistentActionErrorCode = "neb_workflow_non_existent_action_error"

// IsWorkflowNonExistentActionError tests whether a given error is an instance of "non_existent_action_error".
func IsWorkflowNonExistentActionError(err errawrgo.Error) bool {
	return err != nil && err.Is(WorkflowNonExistentActionErrorCode)
}

// IsWorkflowNonExistentActionError tests whether a given error is an instance of "non_existent_action_error".
func (External) IsWorkflowNonExistentActionError(err errawrgo.Error) bool {
	return IsWorkflowNonExistentActionError(err)
}

// WorkflowNonExistentActionErrorBuilder is a builder for "non_existent_action_error" errors.
type WorkflowNonExistentActionErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "non_existent_action_error" from this builder.
func (b *WorkflowNonExistentActionErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "the request action does not exist: {{action}}",
		Technical: "the request action does not exist: {{action}}",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "non_existent_action_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     WorkflowSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "The requested action does not exist",
		Version:          1,
	}
}

// NewWorkflowNonExistentActionErrorBuilder creates a new error builder for the code "non_existent_action_error".
func NewWorkflowNonExistentActionErrorBuilder(action string) *WorkflowNonExistentActionErrorBuilder {
	return &WorkflowNonExistentActionErrorBuilder{arguments: impl.ErrorArguments{"action": impl.NewErrorArgument(action, "the action that was missing")}}
}

// NewWorkflowNonExistentActionError creates a new error with the code "non_existent_action_error".
func NewWorkflowNonExistentActionError(action string) Error {
	return NewWorkflowNonExistentActionErrorBuilder(action).Build()
}

// WorkflowRunnerDecodeErrorCode is the code for an instance of "runner_decode_error".
const WorkflowRunnerDecodeErrorCode = "neb_workflow_runner_decode_error"

// IsWorkflowRunnerDecodeError tests whether a given error is an instance of "runner_decode_error".
func IsWorkflowRunnerDecodeError(err errawrgo.Error) bool {
	return err != nil && err.Is(WorkflowRunnerDecodeErrorCode)
}

// IsWorkflowRunnerDecodeError tests whether a given error is an instance of "runner_decode_error".
func (External) IsWorkflowRunnerDecodeError(err errawrgo.Error) bool {
	return IsWorkflowRunnerDecodeError(err)
}

// WorkflowRunnerDecodeErrorBuilder is a builder for "runner_decode_error" errors.
type WorkflowRunnerDecodeErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "runner_decode_error" from this builder.
func (b *WorkflowRunnerDecodeErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "there was an error decoding the action runner",
		Technical: "there was an error decoding the action runner",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "runner_decode_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     WorkflowSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "Runner decode error",
		Version:          1,
	}
}

// NewWorkflowRunnerDecodeErrorBuilder creates a new error builder for the code "runner_decode_error".
func NewWorkflowRunnerDecodeErrorBuilder() *WorkflowRunnerDecodeErrorBuilder {
	return &WorkflowRunnerDecodeErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewWorkflowRunnerDecodeError creates a new error with the code "runner_decode_error".
func NewWorkflowRunnerDecodeError() Error {
	return NewWorkflowRunnerDecodeErrorBuilder().Build()
}

// WorkflowRunnerNotFoundCode is the code for an instance of "runner_not_found".
const WorkflowRunnerNotFoundCode = "neb_workflow_runner_not_found"

// IsWorkflowRunnerNotFound tests whether a given error is an instance of "runner_not_found".
func IsWorkflowRunnerNotFound(err errawrgo.Error) bool {
	return err != nil && err.Is(WorkflowRunnerNotFoundCode)
}

// IsWorkflowRunnerNotFound tests whether a given error is an instance of "runner_not_found".
func (External) IsWorkflowRunnerNotFound(err errawrgo.Error) bool {
	return IsWorkflowRunnerNotFound(err)
}

// WorkflowRunnerNotFoundBuilder is a builder for "runner_not_found" errors.
type WorkflowRunnerNotFoundBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "runner_not_found" from this builder.
func (b *WorkflowRunnerNotFoundBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "the runner `{{kind}}` was not found",
		Technical: "the runner `{{kind}}` was not found",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "runner_not_found",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     WorkflowSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "Runner not found",
		Version:          1,
	}
}

// NewWorkflowRunnerNotFoundBuilder creates a new error builder for the code "runner_not_found".
func NewWorkflowRunnerNotFoundBuilder(kind string) *WorkflowRunnerNotFoundBuilder {
	return &WorkflowRunnerNotFoundBuilder{arguments: impl.ErrorArguments{"kind": impl.NewErrorArgument(kind, "the kind of runner that was not found")}}
}

// NewWorkflowRunnerNotFound creates a new error with the code "runner_not_found".
func NewWorkflowRunnerNotFound(kind string) Error {
	return NewWorkflowRunnerNotFoundBuilder(kind).Build()
}

// WorkflowStageErrorCode is the code for an instance of "stage_error".
const WorkflowStageErrorCode = "neb_workflow_stage_error"

// IsWorkflowStageError tests whether a given error is an instance of "stage_error".
func IsWorkflowStageError(err errawrgo.Error) bool {
	return err != nil && err.Is(WorkflowStageErrorCode)
}

// IsWorkflowStageError tests whether a given error is an instance of "stage_error".
func (External) IsWorkflowStageError(err errawrgo.Error) bool {
	return IsWorkflowStageError(err)
}

// WorkflowStageErrorBuilder is a builder for "stage_error" errors.
type WorkflowStageErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "stage_error" from this builder.
func (b *WorkflowStageErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "an error occurred while setting up stage of workflow",
		Technical: "an error occurred while setting up stage of workflow",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "stage_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     WorkflowSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "Stage error",
		Version:          1,
	}
}

// NewWorkflowStageErrorBuilder creates a new error builder for the code "stage_error".
func NewWorkflowStageErrorBuilder() *WorkflowStageErrorBuilder {
	return &WorkflowStageErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewWorkflowStageError creates a new error with the code "stage_error".
func NewWorkflowStageError() Error {
	return NewWorkflowStageErrorBuilder().Build()
}

// WorkflowUnknownCommandExecutionErrorCode is the code for an instance of "unknown_command_execution_error".
const WorkflowUnknownCommandExecutionErrorCode = "neb_workflow_unknown_command_execution_error"

// IsWorkflowUnknownCommandExecutionError tests whether a given error is an instance of "unknown_command_execution_error".
func IsWorkflowUnknownCommandExecutionError(err errawrgo.Error) bool {
	return err != nil && err.Is(WorkflowUnknownCommandExecutionErrorCode)
}

// IsWorkflowUnknownCommandExecutionError tests whether a given error is an instance of "unknown_command_execution_error".
func (External) IsWorkflowUnknownCommandExecutionError(err errawrgo.Error) bool {
	return IsWorkflowUnknownCommandExecutionError(err)
}

// WorkflowUnknownCommandExecutionErrorBuilder is a builder for "unknown_command_execution_error" errors.
type WorkflowUnknownCommandExecutionErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "unknown_command_execution_error" from this builder.
func (b *WorkflowUnknownCommandExecutionErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "something unexpected happened with command exection",
		Technical: "something unexpected happened with command exection",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "unknown_command_execution_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     WorkflowSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "Unknown command execution error",
		Version:          1,
	}
}

// NewWorkflowUnknownCommandExecutionErrorBuilder creates a new error builder for the code "unknown_command_execution_error".
func NewWorkflowUnknownCommandExecutionErrorBuilder(errorString string, commandString string) *WorkflowUnknownCommandExecutionErrorBuilder {
	return &WorkflowUnknownCommandExecutionErrorBuilder{arguments: impl.ErrorArguments{
		"command_string": impl.NewErrorArgument(commandString, "Command being executed"),
		"error_string":   impl.NewErrorArgument(errorString, "Error being thrown"),
	}}
}

// NewWorkflowUnknownCommandExecutionError creates a new error with the code "unknown_command_execution_error".
func NewWorkflowUnknownCommandExecutionError(errorString string, commandString string) Error {
	return NewWorkflowUnknownCommandExecutionErrorBuilder(errorString, commandString).Build()
}

// WorkflowUnknownRuntimeErrorCode is the code for an instance of "unknown_runtime_error".
const WorkflowUnknownRuntimeErrorCode = "neb_workflow_unknown_runtime_error"

// IsWorkflowUnknownRuntimeError tests whether a given error is an instance of "unknown_runtime_error".
func IsWorkflowUnknownRuntimeError(err errawrgo.Error) bool {
	return err != nil && err.Is(WorkflowUnknownRuntimeErrorCode)
}

// IsWorkflowUnknownRuntimeError tests whether a given error is an instance of "unknown_runtime_error".
func (External) IsWorkflowUnknownRuntimeError(err errawrgo.Error) bool {
	return IsWorkflowUnknownRuntimeError(err)
}

// WorkflowUnknownRuntimeErrorBuilder is a builder for "unknown_runtime_error" errors.
type WorkflowUnknownRuntimeErrorBuilder struct {
	arguments impl.ErrorArguments
}

// Build creates the error for the code "unknown_runtime_error" from this builder.
func (b *WorkflowUnknownRuntimeErrorBuilder) Build() Error {
	description := &impl.ErrorDescription{
		Friendly:  "an unknown error occurred",
		Technical: "an unknown error occurred",
	}

	return &impl.Error{
		ErrorArguments:   b.arguments,
		ErrorCode:        "unknown_runtime_error",
		ErrorDescription: description,
		ErrorDomain:      Domain,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSection:     WorkflowSection,
		ErrorSensitivity: errawrgo.ErrorSensitivityNone,
		ErrorTitle:       "Unknown runtime error",
		Version:          1,
	}
}

// NewWorkflowUnknownRuntimeErrorBuilder creates a new error builder for the code "unknown_runtime_error".
func NewWorkflowUnknownRuntimeErrorBuilder() *WorkflowUnknownRuntimeErrorBuilder {
	return &WorkflowUnknownRuntimeErrorBuilder{arguments: impl.ErrorArguments{}}
}

// NewWorkflowUnknownRuntimeError creates a new error with the code "unknown_runtime_error".
func NewWorkflowUnknownRuntimeError() Error {
	return NewWorkflowUnknownRuntimeErrorBuilder().Build()
}
