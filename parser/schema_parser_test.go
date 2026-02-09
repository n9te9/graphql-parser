package parser_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/n9te9/graphql-parser/ast"
	"github.com/n9te9/graphql-parser/lexer"
	"github.com/n9te9/graphql-parser/parser"
	"github.com/n9te9/graphql-parser/token"
)

func TestParseObjectTypeDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name: "Simple Type Definition",
			input: `
				type User {
					id: ID!
					username: String
				}
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ObjectTypeDefinition{
						Name: &ast.Name{Value: "User"},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "id"},
								Type: &ast.NonNullType{
									Type: &ast.NamedType{
										Name: &ast.Name{Value: "ID"},
									},
								},
							},
							{
								Name: &ast.Name{
									Value: "username",
								},
								Type: &ast.NamedType{
									Name: &ast.Name{Value: "String"},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Type with Implements and Directives",
			input: `type User implements Node & Entity @key(fields: "id") { id: ID! }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ObjectTypeDefinition{
						Name: &ast.Name{Value: "User"},
						Interfaces: []*ast.NamedType{
							{
								Name: &ast.Name{Value: "Node"},
							},
							{
								Name: &ast.Name{Value: "Entity"},
							},
						},
						Directives: []*ast.Directive{
							{
								Name: "key",
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "fields"},
										Value: &ast.StringValue{Value: "id"},
									},
								},
							},
						},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "id"},
								Type: &ast.NonNullType{
									Type: &ast.NamedType{
										Name: &ast.Name{Value: "ID"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Field with Arguments",
			input: `type Query { user(id: ID!): User }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ObjectTypeDefinition{
						Name: &ast.Name{Value: "Query"},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "user"},
								Arguments: []*ast.InputValueDefinition{
									{
										Name: &ast.Name{Value: "id"},
										Type: &ast.NonNullType{
											Type: &ast.NamedType{Name: &ast.Name{Value: "ID"}},
										},
									},
								},
								Type: &ast.NamedType{Name: &ast.Name{Value: "User"}},
							},
						},
					},
				},
			},
		},
		{
			name:  "Field Argument with Default Value",
			input: `type Query { posts(active: Boolean = true): [Post] }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ObjectTypeDefinition{
						Name: &ast.Name{Value: "Query"},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "posts"},
								Arguments: []*ast.InputValueDefinition{
									{
										Name:         &ast.Name{Value: "active"},
										Type:         &ast.NamedType{Name: &ast.Name{Value: "Boolean"}},
										DefaultValue: &ast.BooleanValue{Value: true},
									},
								},
								Type: &ast.ListType{
									Type: &ast.NamedType{Name: &ast.Name{Value: "Post"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Field and Argument Directives",
			input: `type User { name: String @deprecated(reason: "old") }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ObjectTypeDefinition{
						Name: &ast.Name{Value: "User"},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "name"},
								Type: &ast.NamedType{Name: &ast.Name{Value: "String"}},
								Directives: []*ast.Directive{
									{
										Name: "deprecated",
										Arguments: []*ast.Argument{
											{
												Name:  &ast.Name{Value: "reason"},
												Value: &ast.StringValue{Value: "old"},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			got := p.ParseDocument()

			errors := p.Errors()
			if tt.wantErr != "" {
				if len(errors) == 0 {
					t.Errorf("expected error %q, got none", tt.wantErr)
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

func TestParseObjectTypeDefinition_Strict(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name: "Type with Description (Block String)",
			input: `
                """
                Represents a User in the system.
                Multi-line description.
                """
                type User {
                    id: ID!
                }
            `,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ObjectTypeDefinition{
						Description: "Represents a User in the system.\nMulti-line description.",
						Name:        &ast.Name{Value: "User"},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "id"},
								Type: &ast.NonNullType{
									Type: &ast.NamedType{Name: &ast.Name{Value: "ID"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Fields and Args with Description (String Literal)",
			input: `
                type Query {
                    "Fetch user by ID"
                    user(
                        "The ID of the user"
                        id: ID!
                    ): User
                }
            `,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ObjectTypeDefinition{
						Name: &ast.Name{Value: "Query"},
						Fields: []*ast.FieldDefinition{
							{
								Description: "Fetch user by ID",
								Name:        &ast.Name{Value: "user"},
								Arguments: []*ast.InputValueDefinition{
									{
										Description: "The ID of the user",
										Name:        &ast.Name{Value: "id"},
										Type: &ast.NonNullType{
											Type: &ast.NamedType{Name: &ast.Name{Value: "ID"}},
										},
									},
								},
								Type: &ast.NamedType{Name: &ast.Name{Value: "User"}},
							},
						},
					},
				},
			},
		},
		{
			name: "Description with Directives",
			input: `
                "A deprecated type"
                type OldType @deprecated {
                    val: String
                }
            `,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ObjectTypeDefinition{
						Description: "A deprecated type",
						Name:        &ast.Name{Value: "OldType"},
						Directives: []*ast.Directive{
							{Name: "deprecated"},
						},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "val"},
								Type: &ast.NamedType{Name: &ast.Name{Value: "String"}},
							},
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
				if len(errors) == 0 {
					t.Errorf("expected error %q, got none", tt.wantErr)
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

func TestParseInterfaceTypeDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:  "Simple Interface",
			input: `interface Node { id: ID! }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.InterfaceTypeDefinition{
						Name: &ast.Name{Value: "Node"},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "id"},
								Type: &ast.NonNullType{
									Type: &ast.NamedType{Name: &ast.Name{Value: "ID"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Interface Implementing Interface",
			input: `interface Resource implements Node { id: ID! url: String }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.InterfaceTypeDefinition{
						Name: &ast.Name{Value: "Resource"},
						Interfaces: []*ast.NamedType{
							{Name: &ast.Name{Value: "Node"}},
						},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "id"},
								Type: &ast.NonNullType{
									Type: &ast.NamedType{Name: &ast.Name{Value: "ID"}},
								},
							},
							{
								Name: &ast.Name{Value: "url"},
								Type: &ast.NamedType{Name: &ast.Name{Value: "String"}},
							},
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
				if len(errors) == 0 {
					t.Errorf("expected error %q, got none", tt.wantErr)
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

func TestParseInterfaceTypeDefinition_Strict(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name: "Interface with Description (String Literal)",
			input: `
				"Common fields for Node"
				interface Node {
					id: ID!
				}
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.InterfaceTypeDefinition{
						Description: "Common fields for Node",
						Name:        &ast.Name{Value: "Node"},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "id"},
								Type: &ast.NonNullType{
									Type: &ast.NamedType{Name: &ast.Name{Value: "ID"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Interface with Block String Description and Directives",
			input: `
				"""
				An entity that has a URL.
				Useful for Relay.
				"""
				interface Resource @deprecated(reason: "Use Node instead") {
					url: String!
				}
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.InterfaceTypeDefinition{
						Description: "An entity that has a URL.\nUseful for Relay.",
						Name:        &ast.Name{Value: "Resource"},
						Directives: []*ast.Directive{
							{
								Name: "deprecated",
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "reason"},
										Value: &ast.StringValue{Value: "Use Node instead"},
									},
								},
							},
						},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "url"},
								Type: &ast.NonNullType{
									Type: &ast.NamedType{Name: &ast.Name{Value: "String"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Interface with Field Descriptions",
			input: `
				interface UserLike {
					"The username"
					username: String
				}
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.InterfaceTypeDefinition{
						Name: &ast.Name{Value: "UserLike"},
						Fields: []*ast.FieldDefinition{
							{
								Description: "The username",
								Name:        &ast.Name{Value: "username"},
								Type:        &ast.NamedType{Name: &ast.Name{Value: "String"}},
							},
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
				if len(errors) == 0 {
					t.Errorf("expected error %q, got none", tt.wantErr)
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

func TestParseUnionTypeDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:  "Simple Union",
			input: `union SearchResult = Human | Droid | Starship`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.UnionTypeDefinition{
						Name: &ast.Name{Value: "SearchResult"},
						Types: []*ast.NamedType{
							{Name: &ast.Name{Value: "Human"}},
							{Name: &ast.Name{Value: "Droid"}},
							{Name: &ast.Name{Value: "Starship"}},
						},
					},
				},
			},
		},
		{
			name: "Union with Leading Pipe and Directives",
			input: `
				"Union description"
				union Result @deprecated = | User | Error
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.UnionTypeDefinition{
						Description: "Union description",
						Name:        &ast.Name{Value: "Result"},
						Directives: []*ast.Directive{
							{Name: "deprecated"},
						},
						Types: []*ast.NamedType{
							{Name: &ast.Name{Value: "User"}},
							{Name: &ast.Name{Value: "Error"}},
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
				if len(errors) == 0 {
					t.Errorf("expected error %q, got none", tt.wantErr)
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

func TestParseEnumTypeDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:  "Simple Enum",
			input: `enum Color { RED GREEN BLUE }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.EnumTypeDefinition{
						Name: &ast.Name{Value: "Color"},
						Values: []*ast.EnumValueDefinition{
							{Name: &ast.Name{Value: "RED"}},
							{Name: &ast.Name{Value: "GREEN"}},
							{Name: &ast.Name{Value: "BLUE"}},
						},
					},
				},
			},
		},
		{
			name: "Enum with Descriptions and Directives",
			input: `
				"Episode enum"
				enum Episode {
					"Released in 1977"
					NEWHOPE @deprecated
					EMPIRE
					JEDI
				}
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.EnumTypeDefinition{
						Description: "Episode enum",
						Name:        &ast.Name{Value: "Episode"},
						Values: []*ast.EnumValueDefinition{
							{
								Description: "Released in 1977",
								Name:        &ast.Name{Value: "NEWHOPE"},
								Directives: []*ast.Directive{
									{Name: "deprecated"},
								},
							},
							{Name: &ast.Name{Value: "EMPIRE"}},
							{Name: &ast.Name{Value: "JEDI"}},
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
				if len(errors) == 0 {
					t.Errorf("expected error %q, got none", tt.wantErr)
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

func TestParseUnionTypeDefinition_Strict(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:    "Union without members",
			input:   `union EmptyUnion =`,
			wantErr: "expected",
			expect:  nil,
		},
		{
			name:    "Union with double pipe",
			input:   `union BadUnion = User || Post`,
			wantErr: "expected",
			expect:  nil,
		},
		{
			name:  "Union with keyword as type name",
			input: `union Result = type | query`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.UnionTypeDefinition{
						Name: &ast.Name{Value: "Result"},
						Types: []*ast.NamedType{
							{Name: &ast.Name{Value: "type"}},
							{Name: &ast.Name{Value: "query"}},
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
				if len(errors) == 0 {
					t.Errorf("expected error containing %q, got none", tt.wantErr)
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

func TestParseEnumTypeDefinition_Strict(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:    "Enum with boolean true",
			input:   `enum Bool { true }`,
			wantErr: "unexpected token",
			expect:  nil,
		},
		{
			name:    "Enum with null",
			input:   `enum Void { null }`,
			wantErr: "unexpected token",
			expect:  nil,
		},
		{
			name:  "Enum with keyword values",
			input: `enum Keywords { type query mutation on }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.EnumTypeDefinition{
						Name: &ast.Name{Value: "Keywords"},
						Values: []*ast.EnumValueDefinition{
							{Name: &ast.Name{Value: "type"}},
							{Name: &ast.Name{Value: "query"}},
							{Name: &ast.Name{Value: "mutation"}},
							{Name: &ast.Name{Value: "on"}},
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
				if len(errors) == 0 {
					t.Errorf("expected error containing %q, got none", tt.wantErr)
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

func TestParseScalarTypeDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:  "Simple Scalar",
			input: `scalar Date`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ScalarTypeDefinition{
						Name: &ast.Name{Value: "Date"},
					},
				},
			},
		},
		{
			name: "Scalar with Description and Directive",
			input: `
				"ISO 8601 Date"
				scalar Date @specifiedBy(url: "https://tools.ietf.org/html/rfc3339")
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ScalarTypeDefinition{
						Description: "ISO 8601 Date",
						Name:        &ast.Name{Value: "Date"},
						Directives: []*ast.Directive{
							{
								Name: "specifiedBy",
								Arguments: []*ast.Argument{
									{
										Name:  &ast.Name{Value: "url"},
										Value: &ast.StringValue{Value: "https://tools.ietf.org/html/rfc3339"},
									},
								},
							},
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
				if len(errors) == 0 {
					t.Errorf("expected error containing %q, got none", tt.wantErr)
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

func TestParseInputObjectTypeDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:  "Simple Input Object",
			input: `input UserInput { name: String! age: Int }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.InputObjectTypeDefinition{
						Name: &ast.Name{Value: "UserInput"},
						Fields: []*ast.InputValueDefinition{
							{
								Name: &ast.Name{Value: "name"},
								Type: &ast.NonNullType{
									Type: &ast.NamedType{Name: &ast.Name{Value: "String"}},
								},
							},
							{
								Name: &ast.Name{Value: "age"},
								Type: &ast.NamedType{Name: &ast.Name{Value: "Int"}},
							},
						},
					},
				},
			},
		},
		{
			name: "Input Object with Description, Defaults and Directives",
			input: `
				"Input for creating a user"
				input CreateUserInput @validate {
					"User name"
					name: String! = "Guest" @length(max: 50)
					role: Role = USER
				}
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.InputObjectTypeDefinition{
						Description: "Input for creating a user",
						Name:        &ast.Name{Value: "CreateUserInput"},
						Directives: []*ast.Directive{
							{Name: "validate"},
						},
						Fields: []*ast.InputValueDefinition{
							{
								Description: "User name",
								Name:        &ast.Name{Value: "name"},
								Type: &ast.NonNullType{
									Type: &ast.NamedType{Name: &ast.Name{Value: "String"}},
								},
								DefaultValue: &ast.StringValue{Value: "Guest"},
								Directives: []*ast.Directive{
									{
										Name: "length",
										Arguments: []*ast.Argument{
											{Name: &ast.Name{Value: "max"}, Value: &ast.IntValue{Value: 50}},
										},
									},
								},
							},
							{
								Name:         &ast.Name{Value: "role"},
								Type:         &ast.NamedType{Name: &ast.Name{Value: "Role"}},
								DefaultValue: &ast.EnumValue{Value: "USER"},
							},
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
				if len(errors) == 0 {
					t.Errorf("expected error containing %q, got none", tt.wantErr)
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

func TestParseDirectiveDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:  "Simple Directive Definition",
			input: `directive @delegate on FIELD_DEFINITION`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.DirectiveDefinition{
						Name: &ast.Name{Value: "delegate"},
						Locations: []*ast.Name{
							{Value: "FIELD_DEFINITION"},
						},
					},
				},
			},
		},
		{
			name: "Complex Directive Definition",
			input: `
				"Make a directive repeatable"
				directive @auth(role: String = "USER") repeatable on FIELD | OBJECT
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.DirectiveDefinition{
						Description: "Make a directive repeatable",
						Name:        &ast.Name{Value: "auth"},
						Arguments: []*ast.InputValueDefinition{
							{
								Name:         &ast.Name{Value: "role"},
								Type:         &ast.NamedType{Name: &ast.Name{Value: "String"}},
								DefaultValue: &ast.StringValue{Value: "USER"},
							},
						},
						Repeatable: true,
						Locations: []*ast.Name{
							{Value: "FIELD"},
							{Value: "OBJECT"},
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
				if len(errors) == 0 {
					t.Errorf("expected error containing %q, got none", tt.wantErr)
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

func TestParseSchemaDefinition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:  "Simple Schema",
			input: `schema { query: MyQuery }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.SchemaDefinition{
						OperationTypes: []*ast.OperationTypeDefinition{
							{
								Operation: token.QUERY,
								Type:      &ast.NamedType{Name: &ast.Name{Value: "MyQuery"}},
							},
						},
					},
				},
			},
		},
		{
			name: "Full Schema with Directives",
			input: `
				schema @link(url: "https://specs.apollo.dev/federation/v2.0") {
					query: RootQuery
					mutation: RootMutation
					subscription: RootSubscription
				}
			`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.SchemaDefinition{
						Directives: []*ast.Directive{
							{
								Name: "link",
								Arguments: []*ast.Argument{
									{Name: &ast.Name{Value: "url"}, Value: &ast.StringValue{Value: "https://specs.apollo.dev/federation/v2.0"}},
								},
							},
						},
						OperationTypes: []*ast.OperationTypeDefinition{
							{
								Operation: token.QUERY,
								Type:      &ast.NamedType{Name: &ast.Name{Value: "RootQuery"}},
							},
							{
								Operation: token.MUTATION,
								Type:      &ast.NamedType{Name: &ast.Name{Value: "RootMutation"}},
							},
							{
								Operation: token.SUBSCRIPTION,
								Type:      &ast.NamedType{Name: &ast.Name{Value: "RootSubscription"}},
							},
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
				if len(errors) == 0 {
					t.Errorf("expected error containing %q, got none", tt.wantErr)
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

func TestParseExtension(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
		expect  *ast.Document
	}{
		{
			name:  "Extend Schema",
			input: `extend schema @link(url: "spec") { mutation: MyMutation }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.SchemaExtension{
						Directives: []*ast.Directive{
							{
								Name: "link",
								Arguments: []*ast.Argument{
									{Name: &ast.Name{Value: "url"}, Value: &ast.StringValue{Value: "spec"}},
								},
							},
						},
						OperationTypes: []*ast.OperationTypeDefinition{
							{
								Operation: token.MUTATION,
								Type:      &ast.NamedType{Name: &ast.Name{Value: "MyMutation"}},
							},
						},
					},
				},
			},
		},
		{
			name:  "Extend Type with Fields",
			input: `extend type User { age: Int }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ObjectTypeExtension{
						Name: &ast.Name{Value: "User"},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "age"},
								Type: &ast.NamedType{Name: &ast.Name{Value: "Int"}},
							},
						},
					},
				},
			},
		},
		{
			name:  "Extend Type with Interfaces and Directives",
			input: `extend type User implements Node @key(fields: "id")`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ObjectTypeExtension{
						Name: &ast.Name{Value: "User"},
						Interfaces: []*ast.NamedType{
							{Name: &ast.Name{Value: "Node"}},
						},
						Directives: []*ast.Directive{
							{
								Name: "key",
								Arguments: []*ast.Argument{
									{Name: &ast.Name{Value: "fields"}, Value: &ast.StringValue{Value: "id"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Extend Interface",
			input: `extend interface Node @deprecated { createdAt: String }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.InterfaceTypeExtension{
						Name: &ast.Name{Value: "Node"},
						Directives: []*ast.Directive{
							{Name: "deprecated"},
						},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "createdAt"},
								Type: &ast.NamedType{Name: &ast.Name{Value: "String"}},
							},
						},
					},
				},
			},
		},
		{
			name:  "Extend Union",
			input: `extend union Result = Error | Warning`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.UnionTypeExtension{
						Name: &ast.Name{Value: "Result"},
						Types: []*ast.NamedType{
							{Name: &ast.Name{Value: "Error"}},
							{Name: &ast.Name{Value: "Warning"}},
						},
					},
				},
			},
		},
		{
			name:  "Extend Enum",
			input: `extend enum Color { YELLOW }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.EnumTypeExtension{
						Name: &ast.Name{Value: "Color"},
						Values: []*ast.EnumValueDefinition{
							{Name: &ast.Name{Value: "YELLOW"}},
						},
					},
				},
			},
		},
		{
			name:  "Extend Scalar",
			input: `extend scalar JSON @specifiedBy(url: "json")`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.ScalarTypeExtension{
						Name: &ast.Name{Value: "JSON"},
						Directives: []*ast.Directive{
							{
								Name: "specifiedBy",
								Arguments: []*ast.Argument{
									{Name: &ast.Name{Value: "url"}, Value: &ast.StringValue{Value: "json"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Extend Input",
			input: `extend input Filter { limit: Int }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.InputObjectTypeExtension{
						Name: &ast.Name{Value: "Filter"},
						Fields: []*ast.InputValueDefinition{
							{
								Name: &ast.Name{Value: "limit"},
								Type: &ast.NamedType{Name: &ast.Name{Value: "Int"}},
							},
						},
					},
				},
			},
		},
		{
			name:  "Extend Interface with Implements",
			input: `extend interface Node implements Entity { createdAt: String }`,
			expect: &ast.Document{
				Definitions: []ast.Definition{
					&ast.InterfaceTypeExtension{
						Name: &ast.Name{Value: "Node"},
						Interfaces: []*ast.NamedType{
							{Name: &ast.Name{Value: "Entity"}},
						},
						Fields: []*ast.FieldDefinition{
							{
								Name: &ast.Name{Value: "createdAt"},
								Type: &ast.NamedType{Name: &ast.Name{Value: "String"}},
							},
						},
					},
				},
			},
		},
		{
			name:    "Invalid: Extension with Description",
			input:   `"Do not describe extensions" extend type User { age: Int }`,
			wantErr: "unexpected description",
			expect:  nil,
		},
		{
			name:    "Invalid: Extend Type without anything",
			input:   `extend type User`,
			wantErr: "unexpected",
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
