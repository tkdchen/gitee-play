package gitee

type Namespace struct {
	Id           *int64  `json:"id,omitempty"`
	Type         *string `json:"type,omitempty"`
	Name         *string `json:"name,omitempty"`
	Path         *string `json:"path,omitempty"`
	EnterpriseId *int64  `json:"enterprise_id,omitempty"`
	HtmlURL      *string `json:"html_url,omitempty"`
}

type User struct {
	Id                *int64  `json:"id,omitempty"`
	Login             *string `json:"login,omitempty"`
	Name              *string `json:"name,omitempty"`
	AvatarURL         *string `json:"avatar_url,omitempty"`
	URL               *string `json:"url,omitempty"`
	HtmlURL           *string `json:"html_url,omitempty"`
	Remark            *string `json:"remark,omitempty"`
	FollowersURL      *string `json:"followers_url,omitempty"`
	FollowingURL      *string `json:"following_url,omitempty"`
	GistsURL          *string `json:"gists_url,omitempty"`
	StarredURL        *string `json:"starred_url,omitempty"`
	SubscriptionsURL  *string `json:"subscriptions_url,omitempty"`
	OrganizationsURL  *string `json:"organizations_url,omitempty"`
	ReposURL          *string `json:"repos_url,omitempty"`
	EventsURL         *string `json:"events_url,omitempty"`
	ReceivedEventsURL *string `json:"received_events_url,omitempty"`
	Type              *string `json:"type,omitempty"`
}

// SlimRepository is used as nested repository information.
type SlimRepository struct {
	ID          *int64     `json:"id,omitempty"`
	FullName    *string    `json:"full_name,omitempty"`
	HumanName   *string    `json:"human_name,omitempty"`
	URL         *string    `json:"url,omitempty"`
	Namespace   *Namespace `json:"namespace,omitempty"`
	Path        *string    `json:"path,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Owner       *User      `json:"owner,omitempty"`
	Assigner    *User      `json:"assigner,omitempty"`
	Description *string    `json:"description,omitempty"`
	Private     *bool      `json:"private,omitempty"`
	Public      *bool      `json:"public,omitempty"`
	Internal    *bool      `json:"internal,omitempty"`
	Fork        *bool      `json:"fork,omitempty"`
	HtmlURL     *string    `json:"html_url,omitempty"`
	SshUrl      *string    `json:"ssh_url,omitempty"`
}

type Milestone struct {
	Id           *int64  `json:"id,omitempty"`
	Number       *int64  `json:"number,omitempty"`
	RepositoryId *int64  `json:"repository_id,omitempty"`
	State        *string `json:"state,omitempty"`
	Title        *string `json:"title,omitempty"`
	Description  *string `json:"description,omitempty"`
	UpdatedAt    *string `json:"update_at,omitempty"`
	CreatedAt    *string `json:"created_at,omitempty"`
	DueOn        *string `json:"due_on,omitempty"`
	OpenIssues   *int64  `json:"open_issues,omitempty"`
	ClosedIssues *int64  `json:"closed_issues,omitempty"`
	URL          *string `json:"url,omitempty"`
	HtmlURL      *string `json:"html_url,omitempty"`
}

type Label struct {
	Id           *int64  `json:"id,omitempty"`
	Color        *string `json:"color,omitempty"`
	Name         *string `json:"name,omitempty"`
	RepositoryId *int64  `json:"repository_id,omitempty"`
	Url          *string `json:"url,omitempty"`
	CreatedAt    *string `json:"created_at,omitempty"`
	UpdatedAt    *string `json:"updated_at,omitempty"`
}

type Link struct {
	HREF *string `json:"href,omitempty"`
}
