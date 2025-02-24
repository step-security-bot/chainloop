// Code generated by ent, DO NOT EDIT.

package workflowrun

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/chainloop-dev/chainloop/app/controlplane/internal/biz"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the workflowrun type in the database.
	Label = "workflow_run"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldFinishedAt holds the string denoting the finished_at field in the database.
	FieldFinishedAt = "finished_at"
	// FieldState holds the string denoting the state field in the database.
	FieldState = "state"
	// FieldReason holds the string denoting the reason field in the database.
	FieldReason = "reason"
	// FieldRunURL holds the string denoting the run_url field in the database.
	FieldRunURL = "run_url"
	// FieldRunnerType holds the string denoting the runner_type field in the database.
	FieldRunnerType = "runner_type"
	// FieldAttestation holds the string denoting the attestation field in the database.
	FieldAttestation = "attestation"
	// FieldAttestationDigest holds the string denoting the attestation_digest field in the database.
	FieldAttestationDigest = "attestation_digest"
	// EdgeWorkflow holds the string denoting the workflow edge name in mutations.
	EdgeWorkflow = "workflow"
	// EdgeRobotaccount holds the string denoting the robotaccount edge name in mutations.
	EdgeRobotaccount = "robotaccount"
	// EdgeContractVersion holds the string denoting the contract_version edge name in mutations.
	EdgeContractVersion = "contract_version"
	// EdgeCasBackends holds the string denoting the cas_backends edge name in mutations.
	EdgeCasBackends = "cas_backends"
	// Table holds the table name of the workflowrun in the database.
	Table = "workflow_runs"
	// WorkflowTable is the table that holds the workflow relation/edge.
	WorkflowTable = "workflow_runs"
	// WorkflowInverseTable is the table name for the Workflow entity.
	// It exists in this package in order to avoid circular dependency with the "workflow" package.
	WorkflowInverseTable = "workflows"
	// WorkflowColumn is the table column denoting the workflow relation/edge.
	WorkflowColumn = "workflow_workflowruns"
	// RobotaccountTable is the table that holds the robotaccount relation/edge.
	RobotaccountTable = "workflow_runs"
	// RobotaccountInverseTable is the table name for the RobotAccount entity.
	// It exists in this package in order to avoid circular dependency with the "robotaccount" package.
	RobotaccountInverseTable = "robot_accounts"
	// RobotaccountColumn is the table column denoting the robotaccount relation/edge.
	RobotaccountColumn = "robot_account_workflowruns"
	// ContractVersionTable is the table that holds the contract_version relation/edge.
	ContractVersionTable = "workflow_runs"
	// ContractVersionInverseTable is the table name for the WorkflowContractVersion entity.
	// It exists in this package in order to avoid circular dependency with the "workflowcontractversion" package.
	ContractVersionInverseTable = "workflow_contract_versions"
	// ContractVersionColumn is the table column denoting the contract_version relation/edge.
	ContractVersionColumn = "workflow_run_contract_version"
	// CasBackendsTable is the table that holds the cas_backends relation/edge. The primary key declared below.
	CasBackendsTable = "workflow_run_cas_backends"
	// CasBackendsInverseTable is the table name for the CASBackend entity.
	// It exists in this package in order to avoid circular dependency with the "casbackend" package.
	CasBackendsInverseTable = "cas_backends"
)

// Columns holds all SQL columns for workflowrun fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldFinishedAt,
	FieldState,
	FieldReason,
	FieldRunURL,
	FieldRunnerType,
	FieldAttestation,
	FieldAttestationDigest,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "workflow_runs"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"robot_account_workflowruns",
	"workflow_workflowruns",
	"workflow_run_contract_version",
}

var (
	// CasBackendsPrimaryKey and CasBackendsColumn2 are the table columns denoting the
	// primary key for the cas_backends relation (M2M).
	CasBackendsPrimaryKey = []string{"workflow_run_id", "cas_backend_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

const DefaultState biz.WorkflowRunStatus = "initialized"

// StateValidator is a validator for the "state" field enum values. It is called by the builders before save.
func StateValidator(s biz.WorkflowRunStatus) error {
	switch s {
	case "initialized", "success", "error", "expired", "canceled":
		return nil
	default:
		return fmt.Errorf("workflowrun: invalid enum value for state field: %q", s)
	}
}

// OrderOption defines the ordering options for the WorkflowRun queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByFinishedAt orders the results by the finished_at field.
func ByFinishedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFinishedAt, opts...).ToFunc()
}

// ByState orders the results by the state field.
func ByState(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldState, opts...).ToFunc()
}

// ByReason orders the results by the reason field.
func ByReason(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldReason, opts...).ToFunc()
}

// ByRunURL orders the results by the run_url field.
func ByRunURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRunURL, opts...).ToFunc()
}

// ByRunnerType orders the results by the runner_type field.
func ByRunnerType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRunnerType, opts...).ToFunc()
}

// ByAttestationDigest orders the results by the attestation_digest field.
func ByAttestationDigest(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAttestationDigest, opts...).ToFunc()
}

// ByWorkflowField orders the results by workflow field.
func ByWorkflowField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newWorkflowStep(), sql.OrderByField(field, opts...))
	}
}

// ByRobotaccountField orders the results by robotaccount field.
func ByRobotaccountField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newRobotaccountStep(), sql.OrderByField(field, opts...))
	}
}

// ByContractVersionField orders the results by contract_version field.
func ByContractVersionField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newContractVersionStep(), sql.OrderByField(field, opts...))
	}
}

// ByCasBackendsCount orders the results by cas_backends count.
func ByCasBackendsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCasBackendsStep(), opts...)
	}
}

// ByCasBackends orders the results by cas_backends terms.
func ByCasBackends(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCasBackendsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newWorkflowStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(WorkflowInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, WorkflowTable, WorkflowColumn),
	)
}
func newRobotaccountStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(RobotaccountInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, RobotaccountTable, RobotaccountColumn),
	)
}
func newContractVersionStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ContractVersionInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, ContractVersionTable, ContractVersionColumn),
	)
}
func newCasBackendsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CasBackendsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, CasBackendsTable, CasBackendsPrimaryKey...),
	)
}
