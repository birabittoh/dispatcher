package api

import (
	"context"
	"strconv"
	"strings"

	gitlabwebhook "github.com/flc1125/go-gitlab-webhook/v2"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var (
	_ gitlabwebhook.BuildListener    = (*telegramListener)(nil)
	_ gitlabwebhook.PushListener     = (*telegramListener)(nil)
	_ gitlabwebhook.TagListener      = (*telegramListener)(nil)
	_ gitlabwebhook.IssueListener    = (*telegramListener)(nil)
	_ gitlabwebhook.PipelineListener = (*telegramListener)(nil)
)

type telegramListener struct{}

func (l *telegramListener) OnBuild(ctx context.Context, event *gitlab.BuildEvent) error {
	tbt := ctx.Value("TelegramBotToken").(string)
	tcid := ctx.Value("TelegramChatID").(string)
	ttid := ctx.Value("TelegramThreadID").(string)

	text := `*New build*

Project: ` + event.ProjectName + `
Status: ` + event.BuildStatus + `
Ref: ` + event.Ref + `
Commit: ` + event.Commit.SHA[:8] + ` - ` + event.Commit.Message

	return sendTelegramMessage(tbt, tcid, ttid, text, true)
}

func (l *telegramListener) OnPush(ctx context.Context, event *gitlab.PushEvent) error {
	tbt := ctx.Value("TelegramBotToken").(string)
	tcid := ctx.Value("TelegramChatID").(string)
	ttid := ctx.Value("TelegramThreadID").(string)

	var text strings.Builder
	text.WriteString(`*New push*

Project: ` + event.Project.Name + `
Ref: ` + event.Ref + `
User: ` + event.UserName + `
Commits: ` + strconv.Itoa(len(event.Commits)) + "\n")

	for _, commit := range event.Commits {
		text.WriteString(`
- ` + commit.ID[:8] + ` ` + commit.Message)
	}

	return sendTelegramMessage(tbt, tcid, ttid, text.String(), true)
}

func (l *telegramListener) OnTag(ctx context.Context, event *gitlab.TagEvent) error {
	{
		tbt := ctx.Value("TelegramBotToken").(string)
		tcid := ctx.Value("TelegramChatID").(string)
		ttid := ctx.Value("TelegramThreadID").(string)

		text := `*New tag*

Project: ` + event.Project.Name + `
Tag: ` + event.Ref + `
User: ` + event.UserName

		return sendTelegramMessage(tbt, tcid, ttid, text, true)
	}
}

func (l *telegramListener) OnIssue(ctx context.Context, event *gitlab.IssueEvent) error {
	tbt := ctx.Value("TelegramBotToken").(string)
	tcid := ctx.Value("TelegramChatID").(string)
	ttid := ctx.Value("TelegramThreadID").(string)

	text := `*Issue ` + event.ObjectAttributes.Action + `*

Project: ` + event.Project.Name + `
Title: ` + event.ObjectAttributes.Title + `
User: ` + event.User.Name + `
URL: ` + event.ObjectAttributes.URL

	return sendTelegramMessage(tbt, tcid, ttid, text, false)
}

func (l *telegramListener) OnPipeline(ctx context.Context, event *gitlab.PipelineEvent) error {
	if event.ObjectAttributes.Status != "failed" {
		return nil // Only notify on failed pipelines
	}

	tbt := ctx.Value("TelegramBotToken").(string)
	tcid := ctx.Value("TelegramChatID").(string)
	ttid := ctx.Value("TelegramThreadID").(string)

	text := `*New pipeline*
Project: ` + event.Project.Name + `
Status: ` + event.ObjectAttributes.Status + `
Ref: ` + event.ObjectAttributes.Ref + `
User: ` + event.User.Name

	return sendTelegramMessage(tbt, tcid, ttid, text, false)
}
