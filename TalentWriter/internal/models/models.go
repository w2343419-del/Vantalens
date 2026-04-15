package models

import "time"

type Post struct {
    Title       string `json:"title"`
    Lang        string `json:"lang"`
    Path        string `json:"path"`
    Date        string `json:"date"`
    Status      string `json:"status"`
    StatusColor string `json:"status_color"`
    Pinned      bool   `json:"pinned"`
}

type Comment struct {
    ID          string   `json:"id"`
    Author      string   `json:"author"`
    Email       string   `json:"email"`
    Content     string   `json:"content"`
    Timestamp   string   `json:"timestamp"`
    Approved    bool     `json:"approved"`
    PostPath    string   `json:"post_path"`
    IPAddress   string   `json:"ip_address"`
    UserAgent   string   `json:"user_agent"`
    ParentID    string   `json:"parent_id,omitempty"`
    Images      []string `json:"images,omitempty"`
    IssueNumber int      `json:"issue_number,omitempty"`
}

type CommentSettings struct {
    SMTPEnabled     bool     `json:"smtp_enabled"`
    SMTPHost        string   `json:"smtp_host"`
    SMTPPort        int      `json:"smtp_port"`
    SMTPUser        string   `json:"smtp_user"`
    SMTPPass        string   `json:"smtp_pass"`
    SMTPFrom        string   `json:"smtp_from"`
    SMTPTo          []string `json:"smtp_to"`
    NotifyOnPending bool     `json:"notify_on_pending"`
    BlacklistIPs    []string `json:"blacklist_ips"`
    BlacklistWords  []string `json:"blacklist_keywords"`
}

type CommentsFile struct {
    Comments []Comment `json:"comments"`
}

type Frontmatter struct {
    Title      string
    Draft      bool
    Date       string
    Categories []string
    Pinned     bool
}

type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Content string      `json:"content,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}

type CommentWithPost struct {
    Comment
    PostTitle string `json:"post_title"`
}

type EmailJob struct {
    Settings  CommentSettings
    Comment   Comment
    PostTitle string
    Retries   int
    CreatedAt time.Time
}

type JWTClaims struct {
    Sub string `json:"sub"`
    Iat int64  `json:"iat"`
    Exp int64  `json:"exp"`
    Jti string `json:"jti"`
    Typ string `json:"typ"`
}

type VisitorIP struct {
    IP           string `json:"ip"`
    CommentCount int    `json:"comment_count"`
    FirstSeen    string `json:"first_seen"`
    LastSeen     string `json:"last_seen"`
}

type PageStatistics struct {
    Path  string `json:"path"`
    Title string `json:"title"`
    Views int    `json:"views"`
    UV    int    `json:"uv,omitempty"`
}

type SiteStatistics struct {
    TotalPages      int              `json:"total_pages"`
    TotalViews      int              `json:"total_views"`
    TotalComments   int              `json:"total_comments"`
    PendingComments int              `json:"pending_comments"`
    UniqueIPs       int              `json:"unique_ips"`
    Pages           []PageStatistics `json:"pages,omitempty"`
    Visitors        []VisitorIP      `json:"visitors,omitempty"`
}
