package controllers

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type GraphQL struct {
	transController *Transaction
}

// NewGraphQL create a new GraphQL controller
func NewGraphQL(transController *Transaction) *GraphQL {
	return &GraphQL{
		transController: transController,
	}
}

// newSchemaConfig create a new graphql.SchemaConfig using the resolvers for each GraphQL type define at the begging of
// the function
func (gql *GraphQL) newSchemaConfig() graphql.SchemaConfig {
	// GraphQL ObjectTypes
	transactionType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Transaction",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.Int,
			},
			"Value": &graphql.Field{
				Type: graphql.Float,
			},
			"Note": &graphql.Field{
				Type: graphql.String,
			},
			"Sender": &graphql.Field{
				Type: graphql.String,
			},
			"Receiver": &graphql.Field{
				Type: graphql.String,
			},
		},
	}) // `transactionType` for Golang struct `Transaction`

	// Userful variables
	IDFieldArgument := graphql.FieldConfigArgument{
		"ID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	}
	transactionService := gql.transController.transService

	// Root query for the SchemaConfig
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			// Read a single transaction
			"Transaction": &graphql.Field{
				Type:        transactionType,
				Description: "Get a single transaction",
				Args:        IDFieldArgument,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, isOK := params.Args["ID"].(int)
					if isOK {
						return transactionService.ReadByID(uint(id))
					}

					return nil, errors.New("GraphQL: missing ID")
				},
			},

			// Read all transactions
			"AllTransaction": &graphql.Field{
				Type:        graphql.NewList(transactionType),
				Description: "Get all transactions",
				Resolve: func(_ graphql.ResolveParams) (interface{}, error) {
					return transactionService.ReadAll()
				},
			},
		},
	})

	// Root mutation for the SchemaConfig
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			// Create a transaction
			"AddTransaction": &graphql.Field{
				Type:        transactionType,
				Description: "Create a new transaction",
				Args: graphql.FieldConfigArgument{
					"Value": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
					"Note": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"Sender": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"Receiver": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					value, OK1 := params.Args["Value"].(float64)
					sender, OK2 := params.Args["Sender"].(string)
					receiver, OK3 := params.Args["Receiver"].(string)
					note := params.Args["Note"]

					if OK1 && OK2 && OK3 {
						return gql.transController.NewModel(value, note, sender, receiver)
					}
					return nil, errors.New("GraphQL: missing value, sender, or receiver")
				},
			},

			// Update a transaction
			"UpdateTransaction": &graphql.Field{
				Type:        transactionType,
				Description: "Update a transaction",
				Args: graphql.FieldConfigArgument{
					"ID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"Value": &graphql.ArgumentConfig{
						Type: graphql.Float,
					},
					"Note": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"Sender": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"Receiver": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, OK := params.Args["ID"].(int)
					value := params.Args["Value"]
					note := params.Args["Note"]
					sender := params.Args["Sender"]
					receiver := params.Args["Receiver"]

					if OK {
						return gql.transController.UpdateModel(uint(id), value, note, sender, receiver)
					}
					return nil, errors.New("GraphQL: missing ID")
				},
			},

			// Delete a transaction
			"DeleteTransaction": &graphql.Field{
				Type:        graphql.Int,
				Description: "Delete a transaction",
				Args:        IDFieldArgument,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, isOK := params.Args["ID"].(int)
					if isOK {
						return id, transactionService.Delete(uint(id))
					}

					return 0, errors.New("GraphQL: missing ID")
				},
			},
		},
	})

	return graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	}
}

// NewHandler create a new *handler.Handler and return it.
func (gql *GraphQL) NewHandler() *handler.Handler {
	schema, _ := graphql.NewSchema(gql.newSchemaConfig())
	return handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
}
