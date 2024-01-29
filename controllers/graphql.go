package controllers

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type GraphQL struct {
	tranController *Transaction
	userController *User
}

// NewGraphQL create a new GraphQL controller
func NewGraphQL(tranController *Transaction, userController *User) *GraphQL {
	return &GraphQL{
		tranController: tranController,
		userController: userController,
	}
}

// GraphQL ObjectTypes for Golang struct models.Transaction
var transactionType = graphql.NewObject(graphql.ObjectConfig{
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
})

// GraphQL ObjectTypes for Golang struct models.User
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"ID": &graphql.Field{
			Type: graphql.Int,
		},
		"Email": &graphql.Field{
			Type: graphql.String,
		},
		"Password": &graphql.Field{
			Type: graphql.String,
		},
		"Last": &graphql.Field{
			Type: graphql.String,
		},
		"Middle": &graphql.Field{
			Type: graphql.String,
		},
		"First": &graphql.Field{
			Type: graphql.String,
		},
		"Phone": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// newSchemaConfig create a new graphql.SchemaConfig using the resolvers for each GraphQL type define at the begging of
// the function
func (gql *GraphQL) newSchemaConfig() graphql.SchemaConfig {
	// Userful variables
	IDFieldArgument := graphql.FieldConfigArgument{
		"ID": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	}
	transactionService := gql.tranController.transService
	userService := gql.userController.userService

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

			// Read a single user
			"User": &graphql.Field{
				Type:        userType,
				Description: "Get a single user",
				Args:        IDFieldArgument,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, isOK := params.Args["ID"].(int)
					if isOK {
						return userService.ReadByID(uint(id))
					}

					return nil, errors.New("GraphQL: missing ID")
				},
			},

			// Read all users
			"AllUser": &graphql.Field{
				Type:        graphql.NewList(userType),
				Description: "Get all users",
				Resolve: func(_ graphql.ResolveParams) (interface{}, error) {
					return userService.ReadAll()
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
						return gql.tranController.NewModel(value, note, sender, receiver)
					}
					return nil, errors.New("GraphQL: missing Value, Sender, or Receiver")
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
						return gql.tranController.UpdateModel(uint(id), value, note, sender, receiver)
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

			// Create a user
			"AddUser": &graphql.Field{
				Type:        userType,
				Description: "Create a new user",
				Args: graphql.FieldConfigArgument{
					"Email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"Password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"Last": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"Middle": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"First": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"Phone": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					email, OK1 := params.Args["Email"].(string)
					password, OK2 := params.Args["Password"].(string)
					last, OK3 := params.Args["Last"].(string)
					middle := params.Args["Middle"]
					first := params.Args["First"]
					phone := params.Args["Phone"]

					if OK1 && OK2 && OK3 {
						return gql.userController.NewModel(email, password, last, middle, first, phone)
					}
					return nil, errors.New("GraphQL: missing Email, Password, or Last")
				},
			},

			// Update a user
			"UpdateUser": &graphql.Field{
				Type:        userType,
				Description: "Update a user",
				Args: graphql.FieldConfigArgument{
					"ID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"Email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"Password": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"Last": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"Middle": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"First": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"Phone": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, OK := params.Args["ID"].(int)
					email := params.Args["Email"]
					password := params.Args["Password"]
					last := params.Args["Last"]
					middle := params.Args["Middle"]
					first := params.Args["First"]
					phone := params.Args["Phone"]

					if OK {
						return gql.userController.UpdateModel(uint(id), email, password, last, middle, first, phone)
					}
					return nil, errors.New("GraphQL: missing ID")
				},
			},

			// Delete a user
			"DeleteUser": &graphql.Field{
				Type:        graphql.Int,
				Description: "Delete a user",
				Args:        IDFieldArgument,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, isOK := params.Args["ID"].(int)
					if isOK {
						return id, userService.Delete(uint(id))
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
