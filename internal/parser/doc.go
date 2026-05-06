// Package parser provides utilities for reading and parsing .env files.
//
// A .env file is a plain-text file containing KEY=VALUE pairs, one per line.
// Lines beginning with '#' are treated as comments and ignored. Blank lines
// are also skipped. Values may optionally be wrapped in single or double
// quotes, which are stripped during parsing.
//
// Example usage:
//
//	env, err := parser.ParseFile(".env.production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(env["DATABASE_URL"])
package parser
