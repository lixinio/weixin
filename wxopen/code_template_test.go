package wxopen

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTemplateDraftList(t *testing.T) {
	open := initWxOpen()
	drafts, err := open.GetTemplateDraftList(context.Background())
	require.Empty(t, err)
	fmt.Printf("%v\n", drafts)
}

func TestAddToTemplate(t *testing.T) {
	open := initWxOpen()
	drafts, err := open.GetTemplateDraftList(context.Background())
	require.Empty(t, err)
	require.NotEmpty(t, drafts)

	// get last DraftID from drafts
	draft_id := drafts[len(drafts)-1].DraftID
	fmt.Println(draft_id)
	err = open.AddToTemplate(context.Background(), draft_id)
	require.Empty(t, err)
}

func TestGetTemplateList(t *testing.T) {
	open := initWxOpen()
	drafts, err := open.GetTemplateList(context.Background())
	require.Empty(t, err)
	fmt.Printf("%v\n", drafts)
}

func TestDeleteTemplate(t *testing.T) {
	open := initWxOpen()
	// 实际测试时，请把 -1 替换成自己的模板ID
	open.DeleteTemplate(context.Background(), -1)
}
