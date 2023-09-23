package webhook

import (
	"encoding/json"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestSubSlice(t *testing.T) {
	c := qt.New(t)

	type xd struct {
		input    []string
		upTo     int
		expected []string
	}

	tests := []xd{
		{
			input: []string{
				"a",
				"b",
			},
			upTo: 3,
			expected: []string{
				"a",
				"b",
			},
		},
		{
			input: []string{
				"a",
				"b",
				"c",
			},
			upTo: 3,
			expected: []string{
				"a",
				"b",
				"c",
			},
		},
		{
			input: []string{
				"a",
				"b",
				"c",
				"d",
			},
			upTo: 3,
			expected: []string{
				"a",
				"b",
				"c",
			},
		},
		{
			input:    []string{},
			upTo:     3,
			expected: []string{},
		},
		{
			input:    []string{},
			upTo:     0,
			expected: []string{},
		},
		{
			input: []string{
				"a",
			},
			upTo:     0,
			expected: []string{},
		},
		{
			input: []string{
				"a",
			},
			upTo: 1,
			expected: []string{
				"a",
			},
		},
		{
			input: []string{
				"a",
			},
			upTo:     0,
			expected: []string{},
		},
		{
			input: []string{
				"a",
			},
			upTo:     -1,
			expected: []string{},
		},
	}

	for _, test := range tests {
		actual := subSlice(test.input, test.upTo)
		c.Assert(actual, qt.DeepEquals, test.expected)

	}
}

func TestGenerateTwitchMessages(t *testing.T) {
	c := qt.New(t)

	type xd struct {
		body     string
		expected []string
	}

	tests := []xd{
		{
			// real
			body: `{
  "ref": "refs/heads/master",
  "before": "c71e91200a19fdab5ccfe34b39c76ac14d96cfa3",
  "after": "6860c7007e76471a5f965ebec2434e1434bb72b7",
  "repository": {
    "id": 77624593,
    "node_id": "MDEwOlJlcG9zaXRvcnk3NzYyNDU5Mw==",
    "name": "chatterino2",
    "full_name": "Chatterino/chatterino2",
    "private": false,
    "owner": {
      "name": "Chatterino",
      "email": null,
      "login": "Chatterino",
      "id": 39381366,
      "node_id": "MDEyOk9yZ2FuaXphdGlvbjM5MzgxMzY2",
      "avatar_url": "https://avatars.githubusercontent.com/u/39381366?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/Chatterino",
      "html_url": "https://github.com/Chatterino",
      "followers_url": "https://api.github.com/users/Chatterino/followers",
      "following_url": "https://api.github.com/users/Chatterino/following{/other_user}",
      "gists_url": "https://api.github.com/users/Chatterino/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/Chatterino/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/Chatterino/subscriptions",
      "organizations_url": "https://api.github.com/users/Chatterino/orgs",
      "repos_url": "https://api.github.com/users/Chatterino/repos",
      "events_url": "https://api.github.com/users/Chatterino/events{/privacy}",
      "received_events_url": "https://api.github.com/users/Chatterino/received_events",
      "type": "Organization",
      "site_admin": false
    },
    "html_url": "https://github.com/Chatterino/chatterino2",
    "description": "Chat client for https://twitch.tv",
    "fork": false,
    "url": "https://github.com/Chatterino/chatterino2",
    "forks_url": "https://api.github.com/repos/Chatterino/chatterino2/forks",
    "keys_url": "https://api.github.com/repos/Chatterino/chatterino2/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/Chatterino/chatterino2/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/Chatterino/chatterino2/teams",
    "hooks_url": "https://api.github.com/repos/Chatterino/chatterino2/hooks",
    "issue_events_url": "https://api.github.com/repos/Chatterino/chatterino2/issues/events{/number}",
    "events_url": "https://api.github.com/repos/Chatterino/chatterino2/events",
    "assignees_url": "https://api.github.com/repos/Chatterino/chatterino2/assignees{/user}",
    "branches_url": "https://api.github.com/repos/Chatterino/chatterino2/branches{/branch}",
    "tags_url": "https://api.github.com/repos/Chatterino/chatterino2/tags",
    "blobs_url": "https://api.github.com/repos/Chatterino/chatterino2/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/Chatterino/chatterino2/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/Chatterino/chatterino2/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/Chatterino/chatterino2/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/Chatterino/chatterino2/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/Chatterino/chatterino2/languages",
    "stargazers_url": "https://api.github.com/repos/Chatterino/chatterino2/stargazers",
    "contributors_url": "https://api.github.com/repos/Chatterino/chatterino2/contributors",
    "subscribers_url": "https://api.github.com/repos/Chatterino/chatterino2/subscribers",
    "subscription_url": "https://api.github.com/repos/Chatterino/chatterino2/subscription",
    "commits_url": "https://api.github.com/repos/Chatterino/chatterino2/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/Chatterino/chatterino2/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/Chatterino/chatterino2/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/Chatterino/chatterino2/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/Chatterino/chatterino2/contents/{+path}",
    "compare_url": "https://api.github.com/repos/Chatterino/chatterino2/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/Chatterino/chatterino2/merges",
    "archive_url": "https://api.github.com/repos/Chatterino/chatterino2/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/Chatterino/chatterino2/downloads",
    "issues_url": "https://api.github.com/repos/Chatterino/chatterino2/issues{/number}",
    "pulls_url": "https://api.github.com/repos/Chatterino/chatterino2/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/Chatterino/chatterino2/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/Chatterino/chatterino2/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/Chatterino/chatterino2/labels{/name}",
    "releases_url": "https://api.github.com/repos/Chatterino/chatterino2/releases{/id}",
    "deployments_url": "https://api.github.com/repos/Chatterino/chatterino2/deployments",
    "created_at": 1483028293,
    "updated_at": "2023-09-23T08:40:03Z",
    "pushed_at": 1695481796,
    "git_url": "git://github.com/Chatterino/chatterino2.git",
    "ssh_url": "git@github.com:Chatterino/chatterino2.git",
    "clone_url": "https://github.com/Chatterino/chatterino2.git",
    "svn_url": "https://github.com/Chatterino/chatterino2",
    "homepage": "",
    "size": 14459,
    "stargazers_count": 1786,
    "watchers_count": 1786,
    "language": "C++",
    "has_issues": true,
    "has_projects": false,
    "has_downloads": true,
    "has_wiki": false,
    "has_pages": false,
    "has_discussions": true,
    "forks_count": 420,
    "mirror_url": null,
    "archived": false,
    "disabled": false,
    "open_issues_count": 371,
    "license": {
      "key": "mit",
      "name": "MIT License",
      "spdx_id": "MIT",
      "url": "https://api.github.com/licenses/mit",
      "node_id": "MDc6TGljZW5zZTEz"
    },
    "allow_forking": true,
    "is_template": false,
    "web_commit_signoff_required": false,
    "topics": [
      "chatclient",
      "hacktoberfest",
      "linux",
      "macos",
      "twitch",
      "windows"
    ],
    "visibility": "public",
    "forks": 420,
    "open_issues": 371,
    "watchers": 1786,
    "default_branch": "master",
    "stargazers": 1786,
    "master_branch": "master",
    "organization": "Chatterino"
  },
  "pusher": {
    "name": "pajlada",
    "email": "rasmus.karlsson+github@pajlada.com"
  },
  "organization": {
    "login": "Chatterino",
    "id": 39381366,
    "node_id": "MDEyOk9yZ2FuaXphdGlvbjM5MzgxMzY2",
    "url": "https://api.github.com/orgs/Chatterino",
    "repos_url": "https://api.github.com/orgs/Chatterino/repos",
    "events_url": "https://api.github.com/orgs/Chatterino/events",
    "hooks_url": "https://api.github.com/orgs/Chatterino/hooks",
    "issues_url": "https://api.github.com/orgs/Chatterino/issues",
    "members_url": "https://api.github.com/orgs/Chatterino/members{/member}",
    "public_members_url": "https://api.github.com/orgs/Chatterino/public_members{/member}",
    "avatar_url": "https://avatars.githubusercontent.com/u/39381366?v=4",
    "description": ""
  },
  "sender": {
    "login": "pajlada",
    "id": 962989,
    "node_id": "MDQ6VXNlcjk2Mjk4OQ==",
    "avatar_url": "https://avatars.githubusercontent.com/u/962989?v=4",
    "gravatar_id": "",
    "url": "https://api.github.com/users/pajlada",
    "html_url": "https://github.com/pajlada",
    "followers_url": "https://api.github.com/users/pajlada/followers",
    "following_url": "https://api.github.com/users/pajlada/following{/other_user}",
    "gists_url": "https://api.github.com/users/pajlada/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/pajlada/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/pajlada/subscriptions",
    "organizations_url": "https://api.github.com/users/pajlada/orgs",
    "repos_url": "https://api.github.com/users/pajlada/repos",
    "events_url": "https://api.github.com/users/pajlada/events{/privacy}",
    "received_events_url": "https://api.github.com/users/pajlada/received_events",
    "type": "User",
    "site_admin": false
  },
  "created": false,
  "deleted": false,
  "forced": false,
  "base_ref": null,
  "compare": "https://github.com/Chatterino/chatterino2/compare/c71e91200a19...6860c7007e76",
  "commits": [
    {
      "id": "6860c7007e76471a5f965ebec2434e1434bb72b7",
      "tree_id": "71d331eaa84ab00edb4cdfe29daf3862e4193a78",
      "distinct": true,
      "message": "Fix selection rendering (#4830)\n\nThe rendering of selections was not aligned to the actual selection that took place for newlines at the end of messages, if they were the only part that was selected of that message.\r\n\r\nIn addition to that fix, we've already refactored the MessageLayoutContainer to try to make it a little bit more sane to work with in the future.\r\n\r\nCo-authored-by: Rasmus Karlsson <rasmus.karlsson@pajlada.com>",
      "timestamp": "2023-09-23T15:09:56Z",
      "url": "https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
      "author": {
        "name": "nerix",
        "email": "nerixdev@outlook.de",
        "username": "Nerixyz"
      },
      "committer": {
        "name": "GitHub",
        "email": "noreply@github.com",
        "username": "web-flow"
      },
      "added": [

      ],
      "removed": [

      ],
      "modified": [
        "CHANGELOG.md",
        "src/messages/Selection.hpp",
        "src/messages/layouts/MessageLayout.cpp",
        "src/messages/layouts/MessageLayout.hpp",
        "src/messages/layouts/MessageLayoutContainer.cpp",
        "src/messages/layouts/MessageLayoutContainer.hpp",
        "src/messages/layouts/MessageLayoutElement.cpp",
        "src/messages/layouts/MessageLayoutElement.hpp",
        "src/widgets/helper/ChannelView.cpp"
      ]
    }
  ],
  "head_commit": {
    "id": "6860c7007e76471a5f965ebec2434e1434bb72b7",
    "tree_id": "71d331eaa84ab00edb4cdfe29daf3862e4193a78",
    "distinct": true,
    "message": "Fix selection rendering (#4830)\n\nThe rendering of selections was not aligned to the actual selection that took place for newlines at the end of messages, if they were the only part that was selected of that message.\r\n\r\nIn addition to that fix, we've already refactored the MessageLayoutContainer to try to make it a little bit more sane to work with in the future.\r\n\r\nCo-authored-by: Rasmus Karlsson <rasmus.karlsson@pajlada.com>",
    "timestamp": "2023-09-23T15:09:56Z",
    "url": "https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
    "author": {
      "name": "nerix",
      "email": "nerixdev@outlook.de",
      "username": "Nerixyz"
    },
    "committer": {
      "name": "GitHub",
      "email": "noreply@github.com",
      "username": "web-flow"
    },
    "added": [

    ],
    "removed": [

    ],
    "modified": [
      "CHANGELOG.md",
      "src/messages/Selection.hpp",
      "src/messages/layouts/MessageLayout.cpp",
      "src/messages/layouts/MessageLayout.hpp",
      "src/messages/layouts/MessageLayoutContainer.cpp",
      "src/messages/layouts/MessageLayoutContainer.hpp",
      "src/messages/layouts/MessageLayoutElement.cpp",
      "src/messages/layouts/MessageLayoutElement.hpp",
      "src/widgets/helper/ChannelView.cpp"
    ]
  }
}`,
			expected: []string{
				"Nerixyz (with Rasmus Karlsson) committed to chatterino2@master (2023-09-23 15:09:56 +0000 UTC): Fix selection rendering (#4830) https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
			},
		},
		{
			// fake 1
			body: `{
  "ref": "refs/heads/master",
  "before": "c71e91200a19fdab5ccfe34b39c76ac14d96cfa3",
  "after": "6860c7007e76471a5f965ebec2434e1434bb72b7",
  "repository": {
    "id": 77624593,
    "node_id": "MDEwOlJlcG9zaXRvcnk3NzYyNDU5Mw==",
    "name": "chatterino2",
    "full_name": "Chatterino/chatterino2",
    "private": false,
    "owner": {
      "name": "Chatterino",
      "email": null,
      "login": "Chatterino",
      "id": 39381366,
      "node_id": "MDEyOk9yZ2FuaXphdGlvbjM5MzgxMzY2",
      "avatar_url": "https://avatars.githubusercontent.com/u/39381366?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/Chatterino",
      "html_url": "https://github.com/Chatterino",
      "followers_url": "https://api.github.com/users/Chatterino/followers",
      "following_url": "https://api.github.com/users/Chatterino/following{/other_user}",
      "gists_url": "https://api.github.com/users/Chatterino/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/Chatterino/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/Chatterino/subscriptions",
      "organizations_url": "https://api.github.com/users/Chatterino/orgs",
      "repos_url": "https://api.github.com/users/Chatterino/repos",
      "events_url": "https://api.github.com/users/Chatterino/events{/privacy}",
      "received_events_url": "https://api.github.com/users/Chatterino/received_events",
      "type": "Organization",
      "site_admin": false
    },
    "html_url": "https://github.com/Chatterino/chatterino2",
    "description": "Chat client for https://twitch.tv",
    "fork": false,
    "url": "https://github.com/Chatterino/chatterino2",
    "forks_url": "https://api.github.com/repos/Chatterino/chatterino2/forks",
    "keys_url": "https://api.github.com/repos/Chatterino/chatterino2/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/Chatterino/chatterino2/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/Chatterino/chatterino2/teams",
    "hooks_url": "https://api.github.com/repos/Chatterino/chatterino2/hooks",
    "issue_events_url": "https://api.github.com/repos/Chatterino/chatterino2/issues/events{/number}",
    "events_url": "https://api.github.com/repos/Chatterino/chatterino2/events",
    "assignees_url": "https://api.github.com/repos/Chatterino/chatterino2/assignees{/user}",
    "branches_url": "https://api.github.com/repos/Chatterino/chatterino2/branches{/branch}",
    "tags_url": "https://api.github.com/repos/Chatterino/chatterino2/tags",
    "blobs_url": "https://api.github.com/repos/Chatterino/chatterino2/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/Chatterino/chatterino2/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/Chatterino/chatterino2/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/Chatterino/chatterino2/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/Chatterino/chatterino2/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/Chatterino/chatterino2/languages",
    "stargazers_url": "https://api.github.com/repos/Chatterino/chatterino2/stargazers",
    "contributors_url": "https://api.github.com/repos/Chatterino/chatterino2/contributors",
    "subscribers_url": "https://api.github.com/repos/Chatterino/chatterino2/subscribers",
    "subscription_url": "https://api.github.com/repos/Chatterino/chatterino2/subscription",
    "commits_url": "https://api.github.com/repos/Chatterino/chatterino2/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/Chatterino/chatterino2/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/Chatterino/chatterino2/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/Chatterino/chatterino2/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/Chatterino/chatterino2/contents/{+path}",
    "compare_url": "https://api.github.com/repos/Chatterino/chatterino2/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/Chatterino/chatterino2/merges",
    "archive_url": "https://api.github.com/repos/Chatterino/chatterino2/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/Chatterino/chatterino2/downloads",
    "issues_url": "https://api.github.com/repos/Chatterino/chatterino2/issues{/number}",
    "pulls_url": "https://api.github.com/repos/Chatterino/chatterino2/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/Chatterino/chatterino2/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/Chatterino/chatterino2/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/Chatterino/chatterino2/labels{/name}",
    "releases_url": "https://api.github.com/repos/Chatterino/chatterino2/releases{/id}",
    "deployments_url": "https://api.github.com/repos/Chatterino/chatterino2/deployments",
    "created_at": 1483028293,
    "updated_at": "2023-09-23T08:40:03Z",
    "pushed_at": 1695481796,
    "git_url": "git://github.com/Chatterino/chatterino2.git",
    "ssh_url": "git@github.com:Chatterino/chatterino2.git",
    "clone_url": "https://github.com/Chatterino/chatterino2.git",
    "svn_url": "https://github.com/Chatterino/chatterino2",
    "homepage": "",
    "size": 14459,
    "stargazers_count": 1786,
    "watchers_count": 1786,
    "language": "C++",
    "has_issues": true,
    "has_projects": false,
    "has_downloads": true,
    "has_wiki": false,
    "has_pages": false,
    "has_discussions": true,
    "forks_count": 420,
    "mirror_url": null,
    "archived": false,
    "disabled": false,
    "open_issues_count": 371,
    "license": {
      "key": "mit",
      "name": "MIT License",
      "spdx_id": "MIT",
      "url": "https://api.github.com/licenses/mit",
      "node_id": "MDc6TGljZW5zZTEz"
    },
    "allow_forking": true,
    "is_template": false,
    "web_commit_signoff_required": false,
    "topics": [
      "chatclient",
      "hacktoberfest",
      "linux",
      "macos",
      "twitch",
      "windows"
    ],
    "visibility": "public",
    "forks": 420,
    "open_issues": 371,
    "watchers": 1786,
    "default_branch": "master",
    "stargazers": 1786,
    "master_branch": "master",
    "organization": "Chatterino"
  },
  "pusher": {
    "name": "pajlada",
    "email": "rasmus.karlsson+github@pajlada.com"
  },
  "organization": {
    "login": "Chatterino",
    "id": 39381366,
    "node_id": "MDEyOk9yZ2FuaXphdGlvbjM5MzgxMzY2",
    "url": "https://api.github.com/orgs/Chatterino",
    "repos_url": "https://api.github.com/orgs/Chatterino/repos",
    "events_url": "https://api.github.com/orgs/Chatterino/events",
    "hooks_url": "https://api.github.com/orgs/Chatterino/hooks",
    "issues_url": "https://api.github.com/orgs/Chatterino/issues",
    "members_url": "https://api.github.com/orgs/Chatterino/members{/member}",
    "public_members_url": "https://api.github.com/orgs/Chatterino/public_members{/member}",
    "avatar_url": "https://avatars.githubusercontent.com/u/39381366?v=4",
    "description": ""
  },
  "sender": {
    "login": "pajlada",
    "id": 962989,
    "node_id": "MDQ6VXNlcjk2Mjk4OQ==",
    "avatar_url": "https://avatars.githubusercontent.com/u/962989?v=4",
    "gravatar_id": "",
    "url": "https://api.github.com/users/pajlada",
    "html_url": "https://github.com/pajlada",
    "followers_url": "https://api.github.com/users/pajlada/followers",
    "following_url": "https://api.github.com/users/pajlada/following{/other_user}",
    "gists_url": "https://api.github.com/users/pajlada/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/pajlada/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/pajlada/subscriptions",
    "organizations_url": "https://api.github.com/users/pajlada/orgs",
    "repos_url": "https://api.github.com/users/pajlada/repos",
    "events_url": "https://api.github.com/users/pajlada/events{/privacy}",
    "received_events_url": "https://api.github.com/users/pajlada/received_events",
    "type": "User",
    "site_admin": false
  },
  "created": false,
  "deleted": false,
  "forced": false,
  "base_ref": null,
  "compare": "https://github.com/Chatterino/chatterino2/compare/c71e91200a19...6860c7007e76",
  "commits": [
    {
      "id": "6860c7007e76471a5f965ebec2434e1434bb72b7",
      "tree_id": "71d331eaa84ab00edb4cdfe29daf3862e4193a78",
      "distinct": true,
      "message": "Fix selection rendering (#4830)\n\nThe rendering of selections was not aligned to the actual selection that took place for newlines at the end of messages, if they were the only part that was selected of that message.\r\n\r\nIn addition to that fix, we've already refactored the MessageLayoutContainer to try to make it a little bit more sane to work with in the future.\r\n\r\nCo-authored-by: pajlada <rasmus.karlsson@pajlada.com>",
      "timestamp": "2023-09-23T15:09:56Z",
      "url": "https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
      "author": {
        "name": "nerix",
        "email": "nerixdev@outlook.de",
        "username": "Nerixyz"
      },
      "committer": {
        "name": "GitHub",
        "email": "noreply@github.com",
        "username": "web-flow"
      },
      "added": [

      ],
      "removed": [

      ],
      "modified": [
        "CHANGELOG.md",
        "src/messages/Selection.hpp",
        "src/messages/layouts/MessageLayout.cpp",
        "src/messages/layouts/MessageLayout.hpp",
        "src/messages/layouts/MessageLayoutContainer.cpp",
        "src/messages/layouts/MessageLayoutContainer.hpp",
        "src/messages/layouts/MessageLayoutElement.cpp",
        "src/messages/layouts/MessageLayoutElement.hpp",
        "src/widgets/helper/ChannelView.cpp"
      ]
    }
  ],
  "head_commit": {
    "id": "6860c7007e76471a5f965ebec2434e1434bb72b7",
    "tree_id": "71d331eaa84ab00edb4cdfe29daf3862e4193a78",
    "distinct": true,
    "message": "Fix selection rendering (#4830)\n\nThe rendering of selections was not aligned to the actual selection that took place for newlines at the end of messages, if they were the only part that was selected of that message.\r\n\r\nIn addition to that fix, we've already refactored the MessageLayoutContainer to try to make it a little bit more sane to work with in the future.\r\n\r\nCo-authored-by: Rasmus Karlsson <rasmus.karlsson@pajlada.com>",
    "timestamp": "2023-09-23T15:09:56Z",
    "url": "https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
    "author": {
      "name": "nerix",
      "email": "nerixdev@outlook.de",
      "username": "Nerixyz"
    },
    "committer": {
      "name": "GitHub",
      "email": "noreply@github.com",
      "username": "web-flow"
    },
    "added": [

    ],
    "removed": [

    ],
    "modified": [
      "CHANGELOG.md",
      "src/messages/Selection.hpp",
      "src/messages/layouts/MessageLayout.cpp",
      "src/messages/layouts/MessageLayout.hpp",
      "src/messages/layouts/MessageLayoutContainer.cpp",
      "src/messages/layouts/MessageLayoutContainer.hpp",
      "src/messages/layouts/MessageLayoutElement.cpp",
      "src/messages/layouts/MessageLayoutElement.hpp",
      "src/widgets/helper/ChannelView.cpp"
    ]
  }
}`,
			expected: []string{
				"Nerixyz (with pajlada) committed to chatterino2@master (2023-09-23 15:09:56 +0000 UTC): Fix selection rendering (#4830) https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
			},
		},
		{
			// fake 2
			body: `{
  "ref": "refs/heads/master",
  "before": "c71e91200a19fdab5ccfe34b39c76ac14d96cfa3",
  "after": "6860c7007e76471a5f965ebec2434e1434bb72b7",
  "repository": {
    "id": 77624593,
    "node_id": "MDEwOlJlcG9zaXRvcnk3NzYyNDU5Mw==",
    "name": "chatterino2",
    "full_name": "Chatterino/chatterino2",
    "private": false,
    "owner": {
      "name": "Chatterino",
      "email": null,
      "login": "Chatterino",
      "id": 39381366,
      "node_id": "MDEyOk9yZ2FuaXphdGlvbjM5MzgxMzY2",
      "avatar_url": "https://avatars.githubusercontent.com/u/39381366?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/Chatterino",
      "html_url": "https://github.com/Chatterino",
      "followers_url": "https://api.github.com/users/Chatterino/followers",
      "following_url": "https://api.github.com/users/Chatterino/following{/other_user}",
      "gists_url": "https://api.github.com/users/Chatterino/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/Chatterino/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/Chatterino/subscriptions",
      "organizations_url": "https://api.github.com/users/Chatterino/orgs",
      "repos_url": "https://api.github.com/users/Chatterino/repos",
      "events_url": "https://api.github.com/users/Chatterino/events{/privacy}",
      "received_events_url": "https://api.github.com/users/Chatterino/received_events",
      "type": "Organization",
      "site_admin": false
    },
    "html_url": "https://github.com/Chatterino/chatterino2",
    "description": "Chat client for https://twitch.tv",
    "fork": false,
    "url": "https://github.com/Chatterino/chatterino2",
    "forks_url": "https://api.github.com/repos/Chatterino/chatterino2/forks",
    "keys_url": "https://api.github.com/repos/Chatterino/chatterino2/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/Chatterino/chatterino2/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/Chatterino/chatterino2/teams",
    "hooks_url": "https://api.github.com/repos/Chatterino/chatterino2/hooks",
    "issue_events_url": "https://api.github.com/repos/Chatterino/chatterino2/issues/events{/number}",
    "events_url": "https://api.github.com/repos/Chatterino/chatterino2/events",
    "assignees_url": "https://api.github.com/repos/Chatterino/chatterino2/assignees{/user}",
    "branches_url": "https://api.github.com/repos/Chatterino/chatterino2/branches{/branch}",
    "tags_url": "https://api.github.com/repos/Chatterino/chatterino2/tags",
    "blobs_url": "https://api.github.com/repos/Chatterino/chatterino2/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/Chatterino/chatterino2/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/Chatterino/chatterino2/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/Chatterino/chatterino2/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/Chatterino/chatterino2/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/Chatterino/chatterino2/languages",
    "stargazers_url": "https://api.github.com/repos/Chatterino/chatterino2/stargazers",
    "contributors_url": "https://api.github.com/repos/Chatterino/chatterino2/contributors",
    "subscribers_url": "https://api.github.com/repos/Chatterino/chatterino2/subscribers",
    "subscription_url": "https://api.github.com/repos/Chatterino/chatterino2/subscription",
    "commits_url": "https://api.github.com/repos/Chatterino/chatterino2/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/Chatterino/chatterino2/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/Chatterino/chatterino2/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/Chatterino/chatterino2/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/Chatterino/chatterino2/contents/{+path}",
    "compare_url": "https://api.github.com/repos/Chatterino/chatterino2/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/Chatterino/chatterino2/merges",
    "archive_url": "https://api.github.com/repos/Chatterino/chatterino2/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/Chatterino/chatterino2/downloads",
    "issues_url": "https://api.github.com/repos/Chatterino/chatterino2/issues{/number}",
    "pulls_url": "https://api.github.com/repos/Chatterino/chatterino2/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/Chatterino/chatterino2/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/Chatterino/chatterino2/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/Chatterino/chatterino2/labels{/name}",
    "releases_url": "https://api.github.com/repos/Chatterino/chatterino2/releases{/id}",
    "deployments_url": "https://api.github.com/repos/Chatterino/chatterino2/deployments",
    "created_at": 1483028293,
    "updated_at": "2023-09-23T08:40:03Z",
    "pushed_at": 1695481796,
    "git_url": "git://github.com/Chatterino/chatterino2.git",
    "ssh_url": "git@github.com:Chatterino/chatterino2.git",
    "clone_url": "https://github.com/Chatterino/chatterino2.git",
    "svn_url": "https://github.com/Chatterino/chatterino2",
    "homepage": "",
    "size": 14459,
    "stargazers_count": 1786,
    "watchers_count": 1786,
    "language": "C++",
    "has_issues": true,
    "has_projects": false,
    "has_downloads": true,
    "has_wiki": false,
    "has_pages": false,
    "has_discussions": true,
    "forks_count": 420,
    "mirror_url": null,
    "archived": false,
    "disabled": false,
    "open_issues_count": 371,
    "license": {
      "key": "mit",
      "name": "MIT License",
      "spdx_id": "MIT",
      "url": "https://api.github.com/licenses/mit",
      "node_id": "MDc6TGljZW5zZTEz"
    },
    "allow_forking": true,
    "is_template": false,
    "web_commit_signoff_required": false,
    "topics": [
      "chatclient",
      "hacktoberfest",
      "linux",
      "macos",
      "twitch",
      "windows"
    ],
    "visibility": "public",
    "forks": 420,
    "open_issues": 371,
    "watchers": 1786,
    "default_branch": "master",
    "stargazers": 1786,
    "master_branch": "master",
    "organization": "Chatterino"
  },
  "pusher": {
    "name": "pajlada",
    "email": "rasmus.karlsson+github@pajlada.com"
  },
  "organization": {
    "login": "Chatterino",
    "id": 39381366,
    "node_id": "MDEyOk9yZ2FuaXphdGlvbjM5MzgxMzY2",
    "url": "https://api.github.com/orgs/Chatterino",
    "repos_url": "https://api.github.com/orgs/Chatterino/repos",
    "events_url": "https://api.github.com/orgs/Chatterino/events",
    "hooks_url": "https://api.github.com/orgs/Chatterino/hooks",
    "issues_url": "https://api.github.com/orgs/Chatterino/issues",
    "members_url": "https://api.github.com/orgs/Chatterino/members{/member}",
    "public_members_url": "https://api.github.com/orgs/Chatterino/public_members{/member}",
    "avatar_url": "https://avatars.githubusercontent.com/u/39381366?v=4",
    "description": ""
  },
  "sender": {
    "login": "pajlada",
    "id": 962989,
    "node_id": "MDQ6VXNlcjk2Mjk4OQ==",
    "avatar_url": "https://avatars.githubusercontent.com/u/962989?v=4",
    "gravatar_id": "",
    "url": "https://api.github.com/users/pajlada",
    "html_url": "https://github.com/pajlada",
    "followers_url": "https://api.github.com/users/pajlada/followers",
    "following_url": "https://api.github.com/users/pajlada/following{/other_user}",
    "gists_url": "https://api.github.com/users/pajlada/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/pajlada/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/pajlada/subscriptions",
    "organizations_url": "https://api.github.com/users/pajlada/orgs",
    "repos_url": "https://api.github.com/users/pajlada/repos",
    "events_url": "https://api.github.com/users/pajlada/events{/privacy}",
    "received_events_url": "https://api.github.com/users/pajlada/received_events",
    "type": "User",
    "site_admin": false
  },
  "created": false,
  "deleted": false,
  "forced": false,
  "base_ref": null,
  "compare": "https://github.com/Chatterino/chatterino2/compare/c71e91200a19...6860c7007e76",
  "commits": [
    {
      "id": "6860c7007e76471a5f965ebec2434e1434bb72b7",
      "tree_id": "71d331eaa84ab00edb4cdfe29daf3862e4193a78",
      "distinct": true,
      "message": "Fix selection rendering (#4830)\n\nThe rendering of selections was not aligned to the actual selection that took place for newlines at the end of messages, if they were the only part that was selected of that message.\r\n\r\nIn addition to that fix, we've already refactored the MessageLayoutContainer to try to make it a little bit more sane to work with in the future.\r\n\r\nCo-authored-by: pajlada <rasmus.karlsson@pajlada.com>\r\n\r\nCo-authored-by: pajlada2 <rasmus.karlsson@pajlada.com>\r\n\r\nCo-authored-by: pajlada3 <rasmus.karlsson@pajlada.com>",
      "timestamp": "2023-09-23T15:09:56Z",
      "url": "https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
      "author": {
        "name": "nerix",
        "email": "nerixdev@outlook.de",
        "username": "Nerixyz"
      },
      "committer": {
        "name": "GitHub",
        "email": "noreply@github.com",
        "username": "web-flow"
      },
      "added": [

      ],
      "removed": [

      ],
      "modified": [
        "CHANGELOG.md",
        "src/messages/Selection.hpp",
        "src/messages/layouts/MessageLayout.cpp",
        "src/messages/layouts/MessageLayout.hpp",
        "src/messages/layouts/MessageLayoutContainer.cpp",
        "src/messages/layouts/MessageLayoutContainer.hpp",
        "src/messages/layouts/MessageLayoutElement.cpp",
        "src/messages/layouts/MessageLayoutElement.hpp",
        "src/widgets/helper/ChannelView.cpp"
      ]
    }
  ],
  "head_commit": {
    "id": "6860c7007e76471a5f965ebec2434e1434bb72b7",
    "tree_id": "71d331eaa84ab00edb4cdfe29daf3862e4193a78",
    "distinct": true,
    "message": "Fix selection rendering (#4830)\n\nThe rendering of selections was not aligned to the actual selection that took place for newlines at the end of messages, if they were the only part that was selected of that message.\r\n\r\nIn addition to that fix, we've already refactored the MessageLayoutContainer to try to make it a little bit more sane to work with in the future.\r\n\r\nCo-authored-by: Rasmus Karlsson <rasmus.karlsson@pajlada.com>",
    "timestamp": "2023-09-23T15:09:56Z",
    "url": "https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
    "author": {
      "name": "nerix",
      "email": "nerixdev@outlook.de",
      "username": "Nerixyz"
    },
    "committer": {
      "name": "GitHub",
      "email": "noreply@github.com",
      "username": "web-flow"
    },
    "added": [

    ],
    "removed": [

    ],
    "modified": [
      "CHANGELOG.md",
      "src/messages/Selection.hpp",
      "src/messages/layouts/MessageLayout.cpp",
      "src/messages/layouts/MessageLayout.hpp",
      "src/messages/layouts/MessageLayoutContainer.cpp",
      "src/messages/layouts/MessageLayoutContainer.hpp",
      "src/messages/layouts/MessageLayoutElement.cpp",
      "src/messages/layouts/MessageLayoutElement.hpp",
      "src/widgets/helper/ChannelView.cpp"
    ]
  }
}`,
			expected: []string{
				"Nerixyz (with pajlada, pajlada2, pajlada3) committed to chatterino2@master (2023-09-23 15:09:56 +0000 UTC): Fix selection rendering (#4830) https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
			},
		},
		{
			// fake 3
			body: `{
  "ref": "refs/heads/master",
  "before": "c71e91200a19fdab5ccfe34b39c76ac14d96cfa3",
  "after": "6860c7007e76471a5f965ebec2434e1434bb72b7",
  "repository": {
    "id": 77624593,
    "node_id": "MDEwOlJlcG9zaXRvcnk3NzYyNDU5Mw==",
    "name": "chatterino2",
    "full_name": "Chatterino/chatterino2",
    "private": false,
    "owner": {
      "name": "Chatterino",
      "email": null,
      "login": "Chatterino",
      "id": 39381366,
      "node_id": "MDEyOk9yZ2FuaXphdGlvbjM5MzgxMzY2",
      "avatar_url": "https://avatars.githubusercontent.com/u/39381366?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/Chatterino",
      "html_url": "https://github.com/Chatterino",
      "followers_url": "https://api.github.com/users/Chatterino/followers",
      "following_url": "https://api.github.com/users/Chatterino/following{/other_user}",
      "gists_url": "https://api.github.com/users/Chatterino/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/Chatterino/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/Chatterino/subscriptions",
      "organizations_url": "https://api.github.com/users/Chatterino/orgs",
      "repos_url": "https://api.github.com/users/Chatterino/repos",
      "events_url": "https://api.github.com/users/Chatterino/events{/privacy}",
      "received_events_url": "https://api.github.com/users/Chatterino/received_events",
      "type": "Organization",
      "site_admin": false
    },
    "html_url": "https://github.com/Chatterino/chatterino2",
    "description": "Chat client for https://twitch.tv",
    "fork": false,
    "url": "https://github.com/Chatterino/chatterino2",
    "forks_url": "https://api.github.com/repos/Chatterino/chatterino2/forks",
    "keys_url": "https://api.github.com/repos/Chatterino/chatterino2/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/Chatterino/chatterino2/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/Chatterino/chatterino2/teams",
    "hooks_url": "https://api.github.com/repos/Chatterino/chatterino2/hooks",
    "issue_events_url": "https://api.github.com/repos/Chatterino/chatterino2/issues/events{/number}",
    "events_url": "https://api.github.com/repos/Chatterino/chatterino2/events",
    "assignees_url": "https://api.github.com/repos/Chatterino/chatterino2/assignees{/user}",
    "branches_url": "https://api.github.com/repos/Chatterino/chatterino2/branches{/branch}",
    "tags_url": "https://api.github.com/repos/Chatterino/chatterino2/tags",
    "blobs_url": "https://api.github.com/repos/Chatterino/chatterino2/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/Chatterino/chatterino2/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/Chatterino/chatterino2/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/Chatterino/chatterino2/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/Chatterino/chatterino2/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/Chatterino/chatterino2/languages",
    "stargazers_url": "https://api.github.com/repos/Chatterino/chatterino2/stargazers",
    "contributors_url": "https://api.github.com/repos/Chatterino/chatterino2/contributors",
    "subscribers_url": "https://api.github.com/repos/Chatterino/chatterino2/subscribers",
    "subscription_url": "https://api.github.com/repos/Chatterino/chatterino2/subscription",
    "commits_url": "https://api.github.com/repos/Chatterino/chatterino2/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/Chatterino/chatterino2/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/Chatterino/chatterino2/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/Chatterino/chatterino2/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/Chatterino/chatterino2/contents/{+path}",
    "compare_url": "https://api.github.com/repos/Chatterino/chatterino2/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/Chatterino/chatterino2/merges",
    "archive_url": "https://api.github.com/repos/Chatterino/chatterino2/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/Chatterino/chatterino2/downloads",
    "issues_url": "https://api.github.com/repos/Chatterino/chatterino2/issues{/number}",
    "pulls_url": "https://api.github.com/repos/Chatterino/chatterino2/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/Chatterino/chatterino2/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/Chatterino/chatterino2/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/Chatterino/chatterino2/labels{/name}",
    "releases_url": "https://api.github.com/repos/Chatterino/chatterino2/releases{/id}",
    "deployments_url": "https://api.github.com/repos/Chatterino/chatterino2/deployments",
    "created_at": 1483028293,
    "updated_at": "2023-09-23T08:40:03Z",
    "pushed_at": 1695481796,
    "git_url": "git://github.com/Chatterino/chatterino2.git",
    "ssh_url": "git@github.com:Chatterino/chatterino2.git",
    "clone_url": "https://github.com/Chatterino/chatterino2.git",
    "svn_url": "https://github.com/Chatterino/chatterino2",
    "homepage": "",
    "size": 14459,
    "stargazers_count": 1786,
    "watchers_count": 1786,
    "language": "C++",
    "has_issues": true,
    "has_projects": false,
    "has_downloads": true,
    "has_wiki": false,
    "has_pages": false,
    "has_discussions": true,
    "forks_count": 420,
    "mirror_url": null,
    "archived": false,
    "disabled": false,
    "open_issues_count": 371,
    "license": {
      "key": "mit",
      "name": "MIT License",
      "spdx_id": "MIT",
      "url": "https://api.github.com/licenses/mit",
      "node_id": "MDc6TGljZW5zZTEz"
    },
    "allow_forking": true,
    "is_template": false,
    "web_commit_signoff_required": false,
    "topics": [
      "chatclient",
      "hacktoberfest",
      "linux",
      "macos",
      "twitch",
      "windows"
    ],
    "visibility": "public",
    "forks": 420,
    "open_issues": 371,
    "watchers": 1786,
    "default_branch": "master",
    "stargazers": 1786,
    "master_branch": "master",
    "organization": "Chatterino"
  },
  "pusher": {
    "name": "pajlada",
    "email": "rasmus.karlsson+github@pajlada.com"
  },
  "organization": {
    "login": "Chatterino",
    "id": 39381366,
    "node_id": "MDEyOk9yZ2FuaXphdGlvbjM5MzgxMzY2",
    "url": "https://api.github.com/orgs/Chatterino",
    "repos_url": "https://api.github.com/orgs/Chatterino/repos",
    "events_url": "https://api.github.com/orgs/Chatterino/events",
    "hooks_url": "https://api.github.com/orgs/Chatterino/hooks",
    "issues_url": "https://api.github.com/orgs/Chatterino/issues",
    "members_url": "https://api.github.com/orgs/Chatterino/members{/member}",
    "public_members_url": "https://api.github.com/orgs/Chatterino/public_members{/member}",
    "avatar_url": "https://avatars.githubusercontent.com/u/39381366?v=4",
    "description": ""
  },
  "sender": {
    "login": "pajlada",
    "id": 962989,
    "node_id": "MDQ6VXNlcjk2Mjk4OQ==",
    "avatar_url": "https://avatars.githubusercontent.com/u/962989?v=4",
    "gravatar_id": "",
    "url": "https://api.github.com/users/pajlada",
    "html_url": "https://github.com/pajlada",
    "followers_url": "https://api.github.com/users/pajlada/followers",
    "following_url": "https://api.github.com/users/pajlada/following{/other_user}",
    "gists_url": "https://api.github.com/users/pajlada/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/pajlada/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/pajlada/subscriptions",
    "organizations_url": "https://api.github.com/users/pajlada/orgs",
    "repos_url": "https://api.github.com/users/pajlada/repos",
    "events_url": "https://api.github.com/users/pajlada/events{/privacy}",
    "received_events_url": "https://api.github.com/users/pajlada/received_events",
    "type": "User",
    "site_admin": false
  },
  "created": false,
  "deleted": false,
  "forced": false,
  "base_ref": null,
  "compare": "https://github.com/Chatterino/chatterino2/compare/c71e91200a19...6860c7007e76",
  "commits": [
    {
      "id": "6860c7007e76471a5f965ebec2434e1434bb72b7",
      "tree_id": "71d331eaa84ab00edb4cdfe29daf3862e4193a78",
      "distinct": true,
      "message": "Fix selection rendering (#4830)\n\nThe rendering of selections was not aligned to the actual selection that took place for newlines at the end of messages, if they were the only part that was selected of that message.\r\n\r\nIn addition to that fix, we've already refactored the MessageLayoutContainer to try to make it a little bit more sane to work with in the future.\r\n\r\nCO-AUTHORED-BY: pajlada <rasmus.karlsson@pajlada.com>\r\n\r\nCo-authored-by pajlada2 <rasmus.karlsson@pajlada.com>\r\n\r\nCo-authored-by: pajlada3\r\n\r\nCo-authored-by: pajlada4 <email1> <email2>\r\n\r\nCo-authored-by: pajlada5 <email1> lo li lo\r\n\r\nCo-authored-by: \r\n\r\nCo-authored-by:\r\n\r\nCo-authored-by:  <email>\r\n\r\nCo-authored-by: a<email>\r\n\r\nCo-authored-by:            <email>",
      "timestamp": "2023-09-23T15:09:56Z",
      "url": "https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
      "author": {
        "name": "nerix",
        "email": "nerixdev@outlook.de",
        "username": "Nerixyz"
      },
      "committer": {
        "name": "GitHub",
        "email": "noreply@github.com",
        "username": "web-flow"
      },
      "added": [

      ],
      "removed": [

      ],
      "modified": [
        "CHANGELOG.md",
        "src/messages/Selection.hpp",
        "src/messages/layouts/MessageLayout.cpp",
        "src/messages/layouts/MessageLayout.hpp",
        "src/messages/layouts/MessageLayoutContainer.cpp",
        "src/messages/layouts/MessageLayoutContainer.hpp",
        "src/messages/layouts/MessageLayoutElement.cpp",
        "src/messages/layouts/MessageLayoutElement.hpp",
        "src/widgets/helper/ChannelView.cpp"
      ]
    }
  ],
  "head_commit": {
    "id": "6860c7007e76471a5f965ebec2434e1434bb72b7",
    "tree_id": "71d331eaa84ab00edb4cdfe29daf3862e4193a78",
    "distinct": true,
    "message": "Fix selection rendering (#4830)\n\nThe rendering of selections was not aligned to the actual selection that took place for newlines at the end of messages, if they were the only part that was selected of that message.\r\n\r\nIn addition to that fix, we've already refactored the MessageLayoutContainer to try to make it a little bit more sane to work with in the future.\r\n\r\nCo-authored-by: Rasmus Karlsson <rasmus.karlsson@pajlada.com>",
    "timestamp": "2023-09-23T15:09:56Z",
    "url": "https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
    "author": {
      "name": "nerix",
      "email": "nerixdev@outlook.de",
      "username": "Nerixyz"
    },
    "committer": {
      "name": "GitHub",
      "email": "noreply@github.com",
      "username": "web-flow"
    },
    "added": [

    ],
    "removed": [

    ],
    "modified": [
      "CHANGELOG.md",
      "src/messages/Selection.hpp",
      "src/messages/layouts/MessageLayout.cpp",
      "src/messages/layouts/MessageLayout.hpp",
      "src/messages/layouts/MessageLayoutContainer.cpp",
      "src/messages/layouts/MessageLayoutContainer.hpp",
      "src/messages/layouts/MessageLayoutElement.cpp",
      "src/messages/layouts/MessageLayoutElement.hpp",
      "src/widgets/helper/ChannelView.cpp"
    ]
  }
}`,
			expected: []string{
				"Nerixyz (with pajlada, pajlada4, pajlada5) committed to chatterino2@master (2023-09-23 15:09:56 +0000 UTC): Fix selection rendering (#4830) https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
			},
		},
		{
			// fake
			body: `{
  "ref": "refs/heads/master",
  "before": "c71e91200a19fdab5ccfe34b39c76ac14d96cfa3",
  "after": "6860c7007e76471a5f965ebec2434e1434bb72b7",
  "repository": {
    "id": 77624593,
    "node_id": "MDEwOlJlcG9zaXRvcnk3NzYyNDU5Mw==",
    "name": "chatterino2",
    "full_name": "Chatterino/chatterino2",
    "private": false,
    "owner": {
      "name": "Chatterino",
      "email": null,
      "login": "Chatterino",
      "id": 39381366,
      "node_id": "MDEyOk9yZ2FuaXphdGlvbjM5MzgxMzY2",
      "avatar_url": "https://avatars.githubusercontent.com/u/39381366?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/Chatterino",
      "html_url": "https://github.com/Chatterino",
      "followers_url": "https://api.github.com/users/Chatterino/followers",
      "following_url": "https://api.github.com/users/Chatterino/following{/other_user}",
      "gists_url": "https://api.github.com/users/Chatterino/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/Chatterino/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/Chatterino/subscriptions",
      "organizations_url": "https://api.github.com/users/Chatterino/orgs",
      "repos_url": "https://api.github.com/users/Chatterino/repos",
      "events_url": "https://api.github.com/users/Chatterino/events{/privacy}",
      "received_events_url": "https://api.github.com/users/Chatterino/received_events",
      "type": "Organization",
      "site_admin": false
    },
    "html_url": "https://github.com/Chatterino/chatterino2",
    "description": "Chat client for https://twitch.tv",
    "fork": false,
    "url": "https://github.com/Chatterino/chatterino2",
    "forks_url": "https://api.github.com/repos/Chatterino/chatterino2/forks",
    "keys_url": "https://api.github.com/repos/Chatterino/chatterino2/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/Chatterino/chatterino2/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/Chatterino/chatterino2/teams",
    "hooks_url": "https://api.github.com/repos/Chatterino/chatterino2/hooks",
    "issue_events_url": "https://api.github.com/repos/Chatterino/chatterino2/issues/events{/number}",
    "events_url": "https://api.github.com/repos/Chatterino/chatterino2/events",
    "assignees_url": "https://api.github.com/repos/Chatterino/chatterino2/assignees{/user}",
    "branches_url": "https://api.github.com/repos/Chatterino/chatterino2/branches{/branch}",
    "tags_url": "https://api.github.com/repos/Chatterino/chatterino2/tags",
    "blobs_url": "https://api.github.com/repos/Chatterino/chatterino2/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/Chatterino/chatterino2/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/Chatterino/chatterino2/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/Chatterino/chatterino2/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/Chatterino/chatterino2/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/Chatterino/chatterino2/languages",
    "stargazers_url": "https://api.github.com/repos/Chatterino/chatterino2/stargazers",
    "contributors_url": "https://api.github.com/repos/Chatterino/chatterino2/contributors",
    "subscribers_url": "https://api.github.com/repos/Chatterino/chatterino2/subscribers",
    "subscription_url": "https://api.github.com/repos/Chatterino/chatterino2/subscription",
    "commits_url": "https://api.github.com/repos/Chatterino/chatterino2/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/Chatterino/chatterino2/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/Chatterino/chatterino2/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/Chatterino/chatterino2/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/Chatterino/chatterino2/contents/{+path}",
    "compare_url": "https://api.github.com/repos/Chatterino/chatterino2/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/Chatterino/chatterino2/merges",
    "archive_url": "https://api.github.com/repos/Chatterino/chatterino2/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/Chatterino/chatterino2/downloads",
    "issues_url": "https://api.github.com/repos/Chatterino/chatterino2/issues{/number}",
    "pulls_url": "https://api.github.com/repos/Chatterino/chatterino2/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/Chatterino/chatterino2/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/Chatterino/chatterino2/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/Chatterino/chatterino2/labels{/name}",
    "releases_url": "https://api.github.com/repos/Chatterino/chatterino2/releases{/id}",
    "deployments_url": "https://api.github.com/repos/Chatterino/chatterino2/deployments",
    "created_at": 1483028293,
    "updated_at": "2023-09-23T08:40:03Z",
    "pushed_at": 1695481796,
    "git_url": "git://github.com/Chatterino/chatterino2.git",
    "ssh_url": "git@github.com:Chatterino/chatterino2.git",
    "clone_url": "https://github.com/Chatterino/chatterino2.git",
    "svn_url": "https://github.com/Chatterino/chatterino2",
    "homepage": "",
    "size": 14459,
    "stargazers_count": 1786,
    "watchers_count": 1786,
    "language": "C++",
    "has_issues": true,
    "has_projects": false,
    "has_downloads": true,
    "has_wiki": false,
    "has_pages": false,
    "has_discussions": true,
    "forks_count": 420,
    "mirror_url": null,
    "archived": false,
    "disabled": false,
    "open_issues_count": 371,
    "license": {
      "key": "mit",
      "name": "MIT License",
      "spdx_id": "MIT",
      "url": "https://api.github.com/licenses/mit",
      "node_id": "MDc6TGljZW5zZTEz"
    },
    "allow_forking": true,
    "is_template": false,
    "web_commit_signoff_required": false,
    "topics": [
      "chatclient",
      "hacktoberfest",
      "linux",
      "macos",
      "twitch",
      "windows"
    ],
    "visibility": "public",
    "forks": 420,
    "open_issues": 371,
    "watchers": 1786,
    "default_branch": "master",
    "stargazers": 1786,
    "master_branch": "master",
    "organization": "Chatterino"
  },
  "pusher": {
    "name": "pajlada",
    "email": "rasmus.karlsson+github@pajlada.com"
  },
  "organization": {
    "login": "Chatterino",
    "id": 39381366,
    "node_id": "MDEyOk9yZ2FuaXphdGlvbjM5MzgxMzY2",
    "url": "https://api.github.com/orgs/Chatterino",
    "repos_url": "https://api.github.com/orgs/Chatterino/repos",
    "events_url": "https://api.github.com/orgs/Chatterino/events",
    "hooks_url": "https://api.github.com/orgs/Chatterino/hooks",
    "issues_url": "https://api.github.com/orgs/Chatterino/issues",
    "members_url": "https://api.github.com/orgs/Chatterino/members{/member}",
    "public_members_url": "https://api.github.com/orgs/Chatterino/public_members{/member}",
    "avatar_url": "https://avatars.githubusercontent.com/u/39381366?v=4",
    "description": ""
  },
  "sender": {
    "login": "pajlada",
    "id": 962989,
    "node_id": "MDQ6VXNlcjk2Mjk4OQ==",
    "avatar_url": "https://avatars.githubusercontent.com/u/962989?v=4",
    "gravatar_id": "",
    "url": "https://api.github.com/users/pajlada",
    "html_url": "https://github.com/pajlada",
    "followers_url": "https://api.github.com/users/pajlada/followers",
    "following_url": "https://api.github.com/users/pajlada/following{/other_user}",
    "gists_url": "https://api.github.com/users/pajlada/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/pajlada/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/pajlada/subscriptions",
    "organizations_url": "https://api.github.com/users/pajlada/orgs",
    "repos_url": "https://api.github.com/users/pajlada/repos",
    "events_url": "https://api.github.com/users/pajlada/events{/privacy}",
    "received_events_url": "https://api.github.com/users/pajlada/received_events",
    "type": "User",
    "site_admin": false
  },
  "created": false,
  "deleted": false,
  "forced": false,
  "base_ref": null,
  "compare": "https://github.com/Chatterino/chatterino2/compare/c71e91200a19...6860c7007e76",
  "commits": [
    {
      "id": "6860c7007e76471a5f965ebec2434e1434bb72b7",
      "tree_id": "71d331eaa84ab00edb4cdfe29daf3862e4193a78",
      "distinct": true,
      "message": "Fix selection rendering (#4830)\n\nThe rendering of selections was not aligned to the actual selection that took place for newlines at the end of messages, if they were the only part that was selected of that message.\r\n\r\nIn addition to that fix, we've already refactored the MessageLayoutContainer to try to make it a little bit more sane to work with in the future.\r\n\r\nCo-authored-by: pajlada <rasmus.karlsson@pajlada.com>\r\n\r\nCo-authored-by: pajlada2 <rasmus.karlsson@pajlada.com>\r\n\r\nCo-authored-by: pajlada3 <rasmus.karlsson@pajlada.com\r\n\r\nCo-authored-by: pajlada4 <rasmus.karlsson@pajlada.com\r\n\r\nCo-authored-by: pajlada5 <rasmus.karlsson@pajlada.com\r\n\r\nCo-authored-by: pajlada6 <rasmus.karlsson@pajlada.com>>>>\r\n\r\nCo-authored-by: pajlada7 <rasmus.karlsson@pajlada.com>",
      "timestamp": "2023-09-23T15:09:56Z",
      "url": "https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
      "author": {
        "name": "nerix",
        "email": "nerixdev@outlook.de",
        "username": "Nerixyz"
      },
      "committer": {
        "name": "GitHub",
        "email": "noreply@github.com",
        "username": "web-flow"
      },
      "added": [

      ],
      "removed": [

      ],
      "modified": [
        "CHANGELOG.md",
        "src/messages/Selection.hpp",
        "src/messages/layouts/MessageLayout.cpp",
        "src/messages/layouts/MessageLayout.hpp",
        "src/messages/layouts/MessageLayoutContainer.cpp",
        "src/messages/layouts/MessageLayoutContainer.hpp",
        "src/messages/layouts/MessageLayoutElement.cpp",
        "src/messages/layouts/MessageLayoutElement.hpp",
        "src/widgets/helper/ChannelView.cpp"
      ]
    }
  ],
  "head_commit": {
    "id": "6860c7007e76471a5f965ebec2434e1434bb72b7",
    "tree_id": "71d331eaa84ab00edb4cdfe29daf3862e4193a78",
    "distinct": true,
    "message": "Fix selection rendering (#4830)\n\nThe rendering of selections was not aligned to the actual selection that took place for newlines at the end of messages, if they were the only part that was selected of that message.\r\n\r\nIn addition to that fix, we've already refactored the MessageLayoutContainer to try to make it a little bit more sane to work with in the future.\r\n\r\nCo-authored-by: Rasmus Karlsson <rasmus.karlsson@pajlada.com>",
    "timestamp": "2023-09-23T15:09:56Z",
    "url": "https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
    "author": {
      "name": "nerix",
      "email": "nerixdev@outlook.de",
      "username": "Nerixyz"
    },
    "committer": {
      "name": "GitHub",
      "email": "noreply@github.com",
      "username": "web-flow"
    },
    "added": [

    ],
    "removed": [

    ],
    "modified": [
      "CHANGELOG.md",
      "src/messages/Selection.hpp",
      "src/messages/layouts/MessageLayout.cpp",
      "src/messages/layouts/MessageLayout.hpp",
      "src/messages/layouts/MessageLayoutContainer.cpp",
      "src/messages/layouts/MessageLayoutContainer.hpp",
      "src/messages/layouts/MessageLayoutElement.cpp",
      "src/messages/layouts/MessageLayoutElement.hpp",
      "src/widgets/helper/ChannelView.cpp"
    ]
  }
}`,
			expected: []string{
				"Nerixyz (with pajlada, pajlada2, pajlada3, pajlada4, pajlada5) committed to chatterino2@master (2023-09-23 15:09:56 +0000 UTC): Fix selection rendering (#4830) https://github.com/Chatterino/chatterino2/commit/6860c7007e76471a5f965ebec2434e1434bb72b7",
			},
		},
	}

	for _, test := range tests {
		var pushData PushHookResponse
		err := json.Unmarshal([]byte(test.body), &pushData)
		if err != nil {
			t.Fatal("unable to unmarshal string")
		}

		messages := GenerateTwitchMessages(pushData)

		c.Assert(messages, qt.DeepEquals, test.expected)
	}
}
