package parser_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/n9te9/graphql-parser/ast"
	"github.com/n9te9/graphql-parser/lexer"
	"github.com/n9te9/graphql-parser/parser"
	"github.com/n9te9/graphql-parser/token"
)

func TestParseOperationDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:    "Shorthand Query",
			input:   `{ myField }`,
			wantErr: "",
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "myField"},
							},
						},
					},
				},
			},
		},
		{
			name:    "Named Query",
			input:   `query MyQuery { user }`,
			wantErr: "",
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						Name:      &ast.Name{Value: "MyQuery"},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
							},
						},
					},
				},
			},
		},
		{
			name:    "Mutation",
			input:   `mutation CreateUser { createUser }`,
			wantErr: "",
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Mutation,
						Name:      &ast.Name{Value: "CreateUser"},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "createUser"},
							},
						},
					},
				},
			},
		},
		{
			name:    "Subscription",
			input:   `subscription NewMessages { messageAdded }`,
			wantErr: "",
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Subscription,
						Name:      &ast.Name{Value: "NewMessages"},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "messageAdded"},
							},
						},
					},
				},
			},
		},
		{
			name:    "Missing Closing Brace",
			input:   `query { user`,
			wantErr: "expected next token to be }, got EOF instead",
			expect:  nil,
		},
		{
			name:    "Missing Selection Set",
			input:   `query MyQuery`,
			wantErr: "expected next token to be {, got EOF instead",
			expect:  nil,
		},
		{
			name:    "Invalid Token at Start",
			input:   `123 query { user }`,
			wantErr: "Unexpected token at top level: 123",
			expect:  nil,
		},

		{
			name:  "Query with Variables",
			input: `query getUser($id: ID!) { user(id: $id) }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						Name:      &ast.Name{Value: "getUser"},
						VariableDefinitions: []*ast.VariableDefinition{
							{
								Variable: &ast.Variable{Name: "id"},
								Type: &ast.NonNullType{
									Type: &ast.NamedType{Name: &ast.Name{Value: "ID"}},
								},
							},
						},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "id"},
										Value: &ast.Variable{Name: "id"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Field with Literal Arguments",
			input: `{ user(id: 123, status: "ACTIVE") }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "id"},
										Value: &ast.IntValue{Value: 123},
									},
									{
										Name:  &ast.Name{Value: "status"},
										Value: &ast.StringValue{Value: "ACTIVE"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Directives on Query and Field",
			input: `query @cached(ttl: 60) { user @skip(if: true) }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						Directives: []*ast.Directive{
							{
								Name: "cached",
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "ttl"},
										Value: &ast.IntValue{Value: 60},
									},
								},
							},
						},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
								Directives: []*ast.Directive{
									{
										Name: "skip",
										Arguments: []*ast.Argument{
											{
												Name:  &ast.Name{Value: "if"},
												Value: &ast.BooleanValue{Value: true},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, {
			name: "Complex Query with Fragments and Objects",
			input: `
                query GetUser($input: UserFilter) {
                    user(filter: $input) {
                        id
                        ... on Admin {
                            role
                        }
                        ...UserFields
                    }
                }
            `,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						Name:      &ast.Name{Value: "GetUser"},
						VariableDefinitions: []*ast.VariableDefinition{
							{
								Variable: &ast.Variable{Name: "input"},
								Type:     &ast.NamedType{Name: &ast.Name{Value: "UserFilter"}},
							},
						},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "filter"},
										Value: &ast.Variable{Name: "input"},
									},
								},
								SelectionSet: []ast.Selection{
									&ast.Field{Name: &ast.Name{Value: "id"}},
									&ast.InlineFragment{
										TypeCondition: &ast.NamedType{Name: &ast.Name{Value: "Admin"}},
										SelectionSet: []ast.Selection{
											&ast.Field{Name: &ast.Name{Value: "role"}},
										},
									},
									&ast.FragmentSpread{Name: &ast.Name{Value: "UserFields"}},
								},
							},
						},
					},
				},
			},
		}, {
			name: "Recursive Input Values (Nested Objects and Lists)",
			input: `
                query {
                    complexField(
                        filter: { 
                            tags: ["go", "graphql"], 
                            metadata: { count: 10, active: true } 
                        },
                        matrix: [[1, 2], [3, 4]]
                    )
                }
            `,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "complexField"},
								Arguments: []*ast.Argument{
									{
										Name: &ast.Name{Value: "filter"},
										Value: &ast.ObjectValue{
											Fields: []*ast.ObjectField{
												{
													Name: &ast.Name{Value: "tags"},
													Value: &ast.ListValue{
														Values: []ast.Value{
															&ast.StringValue{Value: "go"},
															&ast.StringValue{Value: "graphql"},
														},
													},
												},
												{
													Name: &ast.Name{Value: "metadata"},
													Value: &ast.ObjectValue{
														Fields: []*ast.ObjectField{
															{
																Name:  &ast.Name{Value: "count"},
																Value: &ast.IntValue{Value: 10},
															},
															{
																Name:  &ast.Name{Value: "active"},
																Value: &ast.BooleanValue{Value: true},
															},
														},
													},
												},
											},
										},
									},
									{
										Name: &ast.Name{Value: "matrix"},
										Value: &ast.ListValue{
											Values: []ast.Value{
												&ast.ListValue{
													Values: []ast.Value{
														&ast.IntValue{Value: 1},
														&ast.IntValue{Value: 2},
													},
												},
												&ast.ListValue{
													Values: []ast.Value{
														&ast.IntValue{Value: 3},
														&ast.IntValue{Value: 4},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Float and Null Literals",
			input: `{ check(val: 123.45, old: null) }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "check"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "val"},
										Value: &ast.FloatValue{Value: 123.45},
									},
									{
										Name:  &ast.Name{Value: "old"},
										Value: &ast.NullValue{},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Alias and Nested Selection Set",
			input: `{ user: currentUser { id, name: fullName } }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Alias: &ast.Name{Value: "user"},
								Name:  &ast.Name{Value: "currentUser"},
								SelectionSet: []ast.Selection{
									&ast.Field{Name: &ast.Name{Value: "id"}},
									&ast.Field{
										Alias: &ast.Name{Value: "name"},
										Name:  &ast.Name{Value: "fullName"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Variables with Default Values",
			input: `query ($id: ID = "default-id", $limit: Int = 10) { user(id: $id) }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						VariableDefinitions: []*ast.VariableDefinition{
							{
								Variable:     &ast.Variable{Name: "id"},
								Type:         &ast.NamedType{Name: &ast.Name{Value: "ID"}},
								DefaultValue: &ast.StringValue{Value: "default-id"},
							},
							{
								Variable:     &ast.Variable{Name: "limit"},
								Type:         &ast.NamedType{Name: &ast.Name{Value: "Int"}},
								DefaultValue: &ast.IntValue{Value: 10},
							},
						},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "id"},
										Value: &ast.Variable{Name: "id"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Inline Fragment without Type Condition",
			input: `{ user { id ... @include(if: $expanded) { email } } }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
								SelectionSet: []ast.Selection{
									&ast.Field{Name: &ast.Name{Value: "id"}},
									&ast.InlineFragment{
										TypeCondition: nil,
										Directives: []*ast.Directive{
											{
												Name: "include",
												Arguments: []*ast.Argument{
													{Name: &ast.Name{Value: "if"}, Value: &ast.Variable{Name: "expanded"}},
												},
											},
										},
										SelectionSet: []ast.Selection{
											&ast.Field{
												Name: &ast.Name{Value: "email"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Fragment with Directives",
			input: `fragment UserFields on User @deprecated(reason: "use newFields") { id }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.FragmentDefinition{
						Name:          &ast.Name{Value: "UserFields"},
						TypeCondition: &ast.NamedType{Name: &ast.Name{Value: "User"}},
						Directives: []*ast.Directive{
							{
								Name: "deprecated",
								Arguments: []*ast.Argument{
									{Name: &ast.Name{Value: "reason"}, Value: &ast.StringValue{Value: "use newFields"}},
								},
							},
						},
						SelectionSet: []ast.Selection{
							&ast.Field{Name: &ast.Name{Value: "id"}},
						},
					},
				},
			},
		},
		{
			name: "Multiple Operations in One Document",
			input: `
        query GetUser { user { id } }
        mutation UpdateUser { updateUser { id } }
    `,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query, Name: &ast.Name{Value: "GetUser"},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
								SelectionSet: []ast.Selection{
									&ast.Field{Name: &ast.Name{Value: "id"}},
								},
							},
						},
					},
					&ast.OperationDefinition{
						Operation: ast.Mutation, Name: &ast.Name{Value: "UpdateUser"},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "updateUser"},
								SelectionSet: []ast.Selection{
									&ast.Field{Name: &ast.Name{Value: "id"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Variable Definition with Directive",
			input: `query ($id: ID! @deprecated(reason: "use uuid")) { user(id: $id) }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						VariableDefinitions: []*ast.VariableDefinition{
							{
								Variable: &ast.Variable{Name: "id"},
								Type: &ast.NonNullType{
									Type: &ast.NamedType{Name: &ast.Name{Value: "ID"}},
								},
								Directives: []*ast.Directive{
									{
										Name: "deprecated",
										Arguments: []*ast.Argument{
											{
												Name:  &ast.Name{Value: "reason"},
												Value: &ast.StringValue{Value: "use uuid"},
											},
										},
									},
								},
							},
						},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "id"},
										Value: &ast.Variable{Name: "id"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Document with BOM",
			input: "\uFEFFquery { id }",
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{Name: &ast.Name{Value: "id"}},
						},
					},
				},
			},
		},
		{
			name:  "Fragment Named on",
			input: `fragment on on User { id }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.FragmentDefinition{
						Name:          &ast.Name{Value: "on"}, // 名前が "on"
						TypeCondition: &ast.NamedType{Name: &ast.Name{Value: "User"}},
						SelectionSet: []ast.Selection{
							&ast.Field{Name: &ast.Name{Value: "id"}},
						},
					},
				},
			},
		},
		{
			name:  "Fragment Spread with Directives",
			input: `query { user { ...UserFields @include(if: $verbose) } }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
								SelectionSet: []ast.Selection{
									&ast.FragmentSpread{
										Name: &ast.Name{Value: "UserFields"},
										Directives: []*ast.Directive{
											{
												Name: "include",
												Arguments: []*ast.Argument{
													{Name: &ast.Name{Value: "if"}, Value: &ast.Variable{Name: "verbose"}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Variable with Default Value and Directive",
			input: `query ($limit: Int = 10 @deprecated) { users(limit: $limit) }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						VariableDefinitions: []*ast.VariableDefinition{
							{
								Variable:     &ast.Variable{Name: "limit"},
								Type:         &ast.NamedType{Name: &ast.Name{Value: "Int"}},
								DefaultValue: &ast.IntValue{Value: 10},
								Directives: []*ast.Directive{
									{Name: "deprecated"},
								},
							},
						},
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "users"},
								Arguments: []*ast.Argument{
									{Name: &ast.Name{Value: "limit"}, Value: &ast.Variable{Name: "limit"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "Empty Selection Set",
			input:   `query { user { } }`,
			wantErr: "empty selection set",
			expect:  nil,
		},
		{
			name:    "Invalid Variable Definition (Missing Colon)",
			input:   `query getUser($id ID!) { user(id: $id) }`,
			wantErr: "expected next token to be :, got IDENT instead",
			expect:  nil,
		},
		{
			name:    "Invalid Argument (Missing Value)",
			input:   `{ user(id: ) }`,
			wantErr: "unexpected token for value",
			expect:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			got := p.ParseDocument()

			errors := p.Errors()

			if tt.wantErr != "" {
				if len(errors) == 0 {
					t.Errorf("expected error containing %q, got none", tt.wantErr)
					return
				}
				found := false
				for _, err := range errors {
					if strings.Contains(err, tt.wantErr) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got %v", tt.wantErr, errors)
				}
				return
			} else {
				if len(errors) > 0 {
					t.Fatalf("unexpected parser errors: %v", errors)
				}
			}

			opts := []cmp.Option{
				cmpopts.IgnoreTypes(token.Token{}),
			}

			if diff := cmp.Diff(tt.expect, got, opts...); diff != "" {
				t.Errorf("ParseDocument() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:    "Keywords as Field Names",
			input:   `{ type query fragment on }`,
			wantErr: "",
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{Name: &ast.Name{Value: "type"}},
							&ast.Field{Name: &ast.Name{Value: "query"}},
							&ast.Field{Name: &ast.Name{Value: "fragment"}},
							&ast.Field{Name: &ast.Name{Value: "on"}},
						},
					},
				},
			},
		},
		{
			name:    "Negative Int and Float",
			input:   `{ calculate(diff: -5, factor: -1.5) }`,
			wantErr: "",
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "calculate"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "diff"},
										Value: &ast.IntValue{Value: -5},
									},
									{
										Name:  &ast.Name{Value: "factor"},
										Value: &ast.FloatValue{Value: -1.5},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "Empty List and Object",
			input:   `{ search(ids: [], filter: {}) }`,
			wantErr: "",
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "search"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "ids"},
										Value: &ast.ListValue{Values: nil},
									},
									{
										Name:  &ast.Name{Value: "filter"},
										Value: &ast.ObjectValue{Fields: nil},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "Boolean vs Enum",
			input:   `{ check(a: true, b: TRUE, c: null, d: NULL) }`,
			wantErr: "",
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "check"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "a"},
										Value: &ast.BooleanValue{Value: true},
									},
									{
										Name:  &ast.Name{Value: "b"},
										Value: &ast.EnumValue{Value: "TRUE"},
									},
									{
										Name:  &ast.Name{Value: "c"},
										Value: &ast.NullValue{},
									},
									{
										Name:  &ast.Name{Value: "d"},
										Value: &ast.EnumValue{Value: "NULL"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "Invalid Syntax (Missing Value)",
			input:   `{ user(id: ) }`,
			wantErr: "unexpected token",
			expect:  nil,
		},
		{
			name:    "Invalid Variable Definition",
			input:   `query($id) { user }`,
			wantErr: "expected",
			expect:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			got := p.ParseDocument()

			errors := p.Errors()

			if tt.wantErr != "" {
				if len(errors) == 0 {
					t.Errorf("expected error containing %q, got none", tt.wantErr)
					return
				}

				found := false
				for _, err := range errors {
					if strings.Contains(err, tt.wantErr) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got %v", tt.wantErr, errors)
				}
				return
			}

			if len(errors) > 0 {
				t.Fatalf("unexpected parser errors: %v", errors)
			}

			opts := []cmp.Option{
				cmpopts.IgnoreTypes(token.Token{}),
			}
			if diff := cmp.Diff(tt.expect, got, opts...); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseStrictSpecCompliance(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name: "Block String with Indentation",
			input: `
				{
					description(text: """
						Hello,
						  World!
					""")
				}
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "description"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "text"},
										Value: &ast.StringValue{Value: "Hello,\n  World!"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "String with Unicode Escapes",
			input: `{ user(name: "\u004E\u0061\u006E\u0061") }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "name"},
										Value: &ast.StringValue{Value: "Nana"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Comments Everywhere",
			input: `
				query { # This is a comment
					user # comment after field
					(id: 1) # comment inside arguments
				}
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{
								Name: &ast.Name{Value: "user"},
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "id"},
										Value: &ast.IntValue{Value: 1},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Complex Default Value",
			input: `query ($filter: Filter = { active: true, tags: ["a", "b"] }) { search }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						VariableDefinitions: []*ast.VariableDefinition{
							{
								Variable: &ast.Variable{Name: "filter"},
								Type:     &ast.NamedType{Name: &ast.Name{Value: "Filter"}},
								DefaultValue: &ast.ObjectValue{
									Fields: []*ast.ObjectField{
										{
											Name:  &ast.Name{Value: "active"},
											Value: &ast.BooleanValue{Value: true},
										},
										{
											Name: &ast.Name{Value: "tags"},
											Value: &ast.ListValue{
												Values: []ast.Value{
													&ast.StringValue{Value: "a"},
													&ast.StringValue{Value: "b"},
												},
											},
										},
									},
								},
							},
						},
						SelectionSet: []ast.Selection{
							&ast.Field{Name: &ast.Name{Value: "search"}},
						},
					},
				},
			},
		},
		{
			name:  "Query Document with BOM",
			input: "\uFEFF{ me }",
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.OperationDefinition{
						Operation: ast.Query,
						SelectionSet: []ast.Selection{
							&ast.Field{Name: &ast.Name{Value: "me"}},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			got := p.ParseDocument()

			errors := p.Errors()
			if tt.wantErr != "" {
				return
			}
			if len(errors) > 0 {
				t.Fatalf("unexpected parser errors: %v", errors)
			}

			opts := []cmp.Option{
				cmpopts.IgnoreTypes(token.Token{}),
				cmpopts.EquateEmpty(),
			}
			if diff := cmp.Diff(tt.expect, got, opts...); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseFragmentDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:    "Valid Fragment",
			input:   `fragment UserFields on User { id name }`,
			wantErr: "",
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.FragmentDefinition{
						Name: &ast.Name{Value: "UserFields"},
						TypeCondition: &ast.NamedType{
							Name: &ast.Name{Value: "User"},
						},
						SelectionSet: []ast.Selection{
							&ast.Field{Name: &ast.Name{Value: "id"}},
							&ast.Field{Name: &ast.Name{Value: "name"}},
						},
					},
				},
			},
		},
		{
			name:    "Fragment Missing On",
			input:   `fragment UserFields User { id }`,
			wantErr: "expected next token to be on, got IDENT instead",
			expect:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			got := p.ParseDocument()

			errors := p.Errors()

			if tt.wantErr != "" {
				if len(errors) == 0 {
					t.Errorf("expected error containing %q, got none", tt.wantErr)
					return
				}
				found := false
				for _, err := range errors {
					if strings.Contains(err, tt.wantErr) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got %v", tt.wantErr, errors)
				}
				return
			} else {
				if len(errors) > 0 {
					t.Fatalf("unexpected parser errors: %v", errors)
				}
			}

			opts := []cmp.Option{
				cmpopts.IgnoreTypes(token.Token{}),
			}

			if diff := cmp.Diff(tt.expect, got, opts...); diff != "" {
				t.Errorf("ParseDocument() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
