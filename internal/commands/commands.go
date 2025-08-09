// Package commands holds commands
package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mmycroft/gator/internal/database"
	"github.com/mmycroft/gator/internal/feed"
	"github.com/mmycroft/gator/internal/state"
)

type Command struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type Commands struct {
	Commands map[string]func(*state.State, Command) error
}

func HandlerLogin(st *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("not enough args to retrieve username")
	}

	user, err := st.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error getting user from cfg: %w", err)
	}

	err = st.Cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error setting user on cfg: %w", err)
	}

	fmt.Printf("set user %s %v: %v\n", user.Name, user.ID, user)

	return nil
}

func HandlerRegister(st *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("not enough args to retrieve username")
	}

	userName := cmd.Args[0]

	id := uuid.New()
	now := time.Now()

	userParams := database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      userName,
	}

	user, err := st.Db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("error creating user in database: %w", err)
	}

	err = st.Cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("error setting user in database: %w", err)
	}

	fmt.Printf("created user %s: %v\n", user.Name, user)

	return nil
}

func HandlerReset(st *state.State, cmd Command) error {
	err := st.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error removing users from database: %w", err)
	}

	fmt.Println("reset database")

	return nil
}

func HandlerUsers(st *state.State, cmd Command) error {
	users, err := st.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users from database: %w", err)
	}

	for i, user := range users {
		var label string
		if st.Cfg.CurrentUserName == user.Name {
			label = "(current)"
		}
		fmt.Println(i+1, user.Name, label)
	}
	return nil
}

func HandlerAgg(st *state.State, cmd Command) error {
	url := "https://www.wagslane.dev/index.xml"

	rss, err := feed.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching rss feed: %w", err)
	}

	fmt.Println(rss)

	return nil
}

func HandlerAddFeed(st *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("no enough arguments to retrieve feed name and url")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	id := uuid.New()
	now := time.Now()

	feedParams := database.CreateFeedParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := st.Db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("error creating feed in database: %w", err)
	}

	fmt.Printf("created feed %s: %v\n", feed.Name, feed)

	id = uuid.New()
	now = time.Now()

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := st.Db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("error creating feed follow in database: %w", err)
	}

	fmt.Printf("created feed follow %s for %s: %v\n", feedFollow.ID, feedFollow.UserName, feedFollow)

	return nil
}

func HandlerFeeds(st *state.State, cmd Command) error {
	feeds, err := st.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds from database: %w", err)
	}

	for i, feed := range feeds {
		fmt.Println(i+1, feed.Name, feed.Url, feed.UserName)
	}
	return nil
}

func HandlerFollow(st *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("not enough args to retrieve url")
	}

	url := cmd.Args[0]
	feed, err := st.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error retrieving feed by url %s from database: %w", url, err)
	}

	id := uuid.New()
	now := time.Now()

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := st.Db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("error creating feed follow in database: %w", err)
	}

	fmt.Printf("created feed follow %s for %s: %v\n", feedFollow.ID, feedFollow.UserName, feedFollow)

	return nil
}

func HandlerFollowing(st *state.State, cmd Command, user database.User) error {
	userFeedFollows, err := st.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feeds for user from database: %w", err)
	}

	for i, feedFollow := range userFeedFollows {
		fmt.Println(i+1, feedFollow.FeedName, feedFollow.UserName)
	}
	return nil
}

func HandlerUnfollow(st *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("not enough args to retrieve url")
	}
	feedUrl := cmd.Args[0]

	deleteFeedParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    feedUrl,
	}

	err := st.Db.DeleteFeedFollow(context.Background(), deleteFeedParams)
	if err != nil {
		return fmt.Errorf("error deleting feed from database: %w", err)
	}

	return nil
}

func MiddlewareLoggedIn(handler func(st *state.State, cmd Command, user database.User) error) func(*state.State, Command) error {
	return func(st *state.State, cmd Command) error {
		userName := st.Cfg.CurrentUserName
		currentUser, err := st.Db.GetUser(context.Background(), userName)
		if err != nil {
			return fmt.Errorf("error retrieving %s from database: %w", userName, err)
		}
		return handler(st, cmd, currentUser)
	}
}

func (c *Commands) Run(st *state.State, cmd Command) error {
	command, ok := c.Commands[cmd.Name]
	if !ok {
		return fmt.Errorf("command %s is not present in commands", cmd.Name)
	}
	return command(st, cmd)
}

func (c *Commands) Register(name string, fn func(*state.State, Command) error) {
	c.Commands[name] = fn
}
