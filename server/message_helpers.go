package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

// validateMoveOrCopy performs validation on a provided post list to determine
// if all permissions are in place to allow the for the posts to be moved or
// copied.
func (p *Plugin) validateMoveOrCopy(wpl *WranglerPostList, originalChannel *model.Channel, targetChannel *model.Channel, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	if wpl.NumPosts() == 0 {
		return nil, false, errors.New("The wrangler post list contains no posts")
	}

	config := p.getConfiguration()

	switch originalChannel.Type {
	case model.CHANNEL_PRIVATE:
		if !config.MoveThreadFromPrivateChannelEnable {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Wrangler is currently configured to not allow moving posts from private channels"), false, nil
		}
	case model.CHANNEL_DIRECT:
		if !config.MoveThreadFromDirectMessageChannelEnable {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Wrangler is currently configured to not allow moving posts from direct message channels"), false, nil
		}
	case model.CHANNEL_GROUP:
		if !config.MoveThreadFromGroupMessageChannelEnable {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Wrangler is currently configured to not allow moving posts from group message channels"), false, nil
		}
	}

	if config.MaxThreadCountMoveSizeInt() != 0 && config.MaxThreadCountMoveSizeInt() < wpl.NumPosts() {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: the thread is %d posts long, but this command is configured to only move threads of up to %d posts", wpl.NumPosts(), config.MaxThreadCountMoveSizeInt())), true, nil
	}

	if wpl.RootPost().ChannelId != extra.ChannelId {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error: this command must be run from the channel containing the post"), true, nil
	}

	_, appErr := p.API.GetChannelMember(targetChannel.Id, extra.UserId)
	if appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: channel with ID %s doesn't exist or you are not a member", targetChannel.Id)), true, nil
	}

	if !config.MoveThreadToAnotherTeamEnable && targetChannel.TeamId != originalChannel.TeamId {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Wrangler is currently configured to not allow moving messages to different teams"), false, nil
	}

	if extra.RootId == wpl.RootPost().Id || extra.ParentId == wpl.RootPost().Id {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error: this command cannot be run from inside the thread; please run directly in the channel containing the thread"), true, nil
	}

	return nil, false, nil
}

func (p *Plugin) copyWranglerPostlist(wpl *WranglerPostList, targetChannel *model.Channel) (*model.Post, error) {
	var appErr *model.AppError
	var newRootPost *model.Post

	if wpl.ContainsFileAttachments() {
		// The thread contains at least one attachment. To properly move the
		// thread, the files will have to be re-uploaded. This is completed
		// before any messages are moved.
		// TODO: check number of files that need to be re-uploaded or file size?
		p.API.LogInfo("Wrangler is re-uploading file attachments",
			"file_count", wpl.FileAttachmentCount,
		)

		for _, post := range wpl.Posts {
			var newFileIDs []string
			var fileBytes []byte
			var oldFileInfo, newFileInfo *model.FileInfo
			for _, fileID := range post.FileIds {
				oldFileInfo, appErr = p.API.GetFileInfo(fileID)
				if appErr != nil {
					return nil, errors.Wrap(appErr, "unable to lookup file info to re-upload")
				}
				fileBytes, appErr = p.API.GetFile(fileID)
				if appErr != nil {
					return nil, errors.Wrap(appErr, "unable to get file bytes to re-upload")
				}
				newFileInfo, appErr = p.API.UploadFile(fileBytes, targetChannel.Id, oldFileInfo.Name)
				if appErr != nil {
					return nil, errors.Wrap(appErr, "unable to re-upload file")
				}

				newFileIDs = append(newFileIDs, newFileInfo.Id)
			}

			post.FileIds = newFileIDs
		}
	}

	for i, post := range wpl.Posts {
		var reactions []*model.Reaction

		// Store reactions to be reapplied later.
		reactions, appErr = p.API.GetReactions(post.Id)
		if appErr != nil {
			// Reaction-based errors are logged, but do not cause the plugin to
			// abort the move thread process.
			p.API.LogError("Failed to get reactions on original post", "err", appErr)
		}

		newPost := post.Clone()
		cleanPost(newPost)
		newPost.ChannelId = targetChannel.Id

		if i == 0 {
			newPost, appErr = p.API.CreatePost(newPost)
			if appErr != nil {
				return nil, errors.Wrap(appErr, "unable to create new root post")
			}
			newRootPost = newPost.Clone()
		} else {
			newPost.RootId = newRootPost.Id
			newPost.ParentId = newRootPost.Id
			newPost, appErr = p.API.CreatePost(newPost)
			if appErr != nil {
				return nil, errors.Wrap(appErr, "unable to create new post")
			}
		}

		for _, reaction := range reactions {
			reaction.PostId = newPost.Id
			_, appErr = p.API.AddReaction(reaction)
			if appErr != nil {
				// Reaction-based errors are logged, but do not cause the plugin to
				// abort the move thread process.
				p.API.LogError("Failed to reapply reactions to post", "err", appErr)
			}
		}
	}

	return newRootPost, nil
}
