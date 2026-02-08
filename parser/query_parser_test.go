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
						// ... SelectionSet
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
										TypeCondition: nil, // ここがnilでパースされるか
										Directives: []*ast.Directive{
											{
												Name: "include",
												Arguments: []*ast.Argument{
													{Name: &ast.Name{Value: "if"}, Value: &ast.Variable{Name: "expanded"}},
												},
											},
										},
										SelectionSet: []ast.Selection{
											&ast.Field{Name: &ast.Name{Value: "email"}},
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
					&ast.OperationDefinition{Operation: ast.Query, Name: &ast.Name{Value: "GetUser"}},
					&ast.OperationDefinition{Operation: ast.Mutation, Name: &ast.Name{Value: "UpdateUser"}},
				},
			},
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
