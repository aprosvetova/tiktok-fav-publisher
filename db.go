package main

import "context"

var ctx = context.TODO()

func wasAlreadyPosted(id string) bool {
	return r.SAdd(ctx, "postedVideos", id).Val() == 0
}
