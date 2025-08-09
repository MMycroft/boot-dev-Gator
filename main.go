// Package main holds the main program
package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/mmycroft/gator/internal/commands"
	"github.com/mmycroft/gator/internal/config"
	"github.com/mmycroft/gator/internal/database"
	"github.com/mmycroft/gator/internal/state"
)

var (
	st   *state.State
	cmds commands.Commands
)

func init() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading config")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println("error opening database")
		os.Exit(1)
	}

	st = &state.State{
		Db:  database.New(db),
		Cfg: &cfg,
	}

	cmds = commands.Commands{
		Commands: make(map[string]func(*state.State, commands.Command) error),
	}

	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerUsers)
	cmds.Register("agg", commands.HandlerAgg)
	cmds.Register("feeds", commands.HandlerFeeds)
	cmds.Register("addfeed", commands.MiddlewareLoggedIn(commands.HandlerAddFeed))
	cmds.Register("follow", commands.MiddlewareLoggedIn(commands.HandlerFollow))
	cmds.Register("following", commands.MiddlewareLoggedIn(commands.HandlerFollowing))
	cmds.Register("unfollow", commands.MiddlewareLoggedIn(commands.HandlerUnfollow))
	cmds.Register("browse", commands.MiddlewareLoggedIn(commands.HandlerBrowse))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("main: not enough arguments")
		os.Exit(1)
	}

	cmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	err := cmds.Run(st, cmd)
	if err != nil {
		fmt.Printf("main: error running command %s: %v\n", cmd.Name, err)
		os.Exit(1)
	}
}
