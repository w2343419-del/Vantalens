package comment

import (
    "encoding/json"
    "os"
    "path/filepath"
    "strings"
    "time"

    "vantalens/talentwriter/internal/config"
    "vantalens/talentwriter/internal/models"
)

func GetCommentsPath(postPath string) string {
    cfg := config.GetConfig()
    if cfg == nil {
        return ""
    }
    return filepath.Join(cfg.HugoPath, "data", "comments", postPath+".json")
}

func GetComments(postPath string) ([]models.Comment, error) {
    path := GetCommentsPath(postPath)
    if path == "" {
        return nil, nil
    }
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return []models.Comment{}, nil
    }
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var cf models.CommentsFile
    if err := json.Unmarshal(data, &cf); err != nil {
        return nil, err
    }
    return cf.Comments, nil
}

func SaveComments(postPath string, comments []models.Comment) error {
    path := GetCommentsPath(postPath)
    if path == "" {
        return nil
    }
    os.MkdirAll(filepath.Dir(path), 0755)
    cf := models.CommentsFile{Comments: comments}
    data, _ := json.MarshalIndent(cf, "", " ")
    return os.WriteFile(path, data, 0644)
}

func AddComment(postPath, author, email, content, ipAddress, userAgent, parentID string) (models.Comment, error) {
    comments, err := GetComments(postPath)
    if err != nil {
        return models.Comment{}, err
    }
    comment := models.Comment{
        ID:        generateCommentID(),
        Author:    author,
        Email:     email,
        Content:   content,
        Timestamp: time.Now().Format(time.RFC3339),
        Approved:  false,
        PostPath:  postPath,
        IPAddress: ipAddress,
        UserAgent: userAgent,
        ParentID:  parentID,
    }
    comments = append(comments, comment)
    return comment, SaveComments(postPath, comments)
}

func ApproveComment(postPath, commentID string) error {
    comments, err := GetComments(postPath)
    if err != nil {
        return err
    }
    for i := range comments {
        if comments[i].ID == commentID {
            comments[i].Approved = true
            break
        }
    }
    return SaveComments(postPath, comments)
}

func DeleteComment(postPath, commentID string) error {
    comments, err := GetComments(postPath)
    if err != nil {
        return err
    }
    var filtered []models.Comment
    for _, c := range comments {
        if c.ID != commentID {
            filtered = append(filtered, c)
        }
    }
    return SaveComments(postPath, filtered)
}

func IsBlacklisted(settings models.CommentSettings, ip, author, email, content string) bool {
    ip = strings.TrimSpace(strings.ToLower(ip))
    text := strings.ToLower(strings.Join([]string{author, email, content}, " "))
    for _, b := range settings.BlacklistIPs {
        if strings.TrimSpace(strings.ToLower(b)) != "" && ip != "" && strings.Contains(ip, strings.TrimSpace(strings.ToLower(b))) {
            return true
        }
    }
    for _, w := range settings.BlacklistWords {
        keyword := strings.TrimSpace(strings.ToLower(w))
        if keyword != "" && strings.Contains(text, keyword) {
            return true
        }
    }
    return false
}

func LoadSettings() models.CommentSettings {
    cfg := config.GetConfig()
    if cfg == nil {
        return models.CommentSettings{}
    }
    path := config.GetCommentSettingsPath(cfg.HugoPath)
    settings := models.CommentSettings{
        SMTPEnabled:     false,
        SMTPPort:        587,
        NotifyOnPending: true,
    }
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return settings
    }
    data, err := os.ReadFile(path)
    if err != nil {
        return settings
    }
    json.Unmarshal(data, &settings)
    return settings
}

func SaveSettings(settings models.CommentSettings) error {
    cfg := config.GetConfig()
    if cfg == nil {
        return nil
    }
    path := config.GetCommentSettingsPath(cfg.HugoPath)
    os.MkdirAll(filepath.Dir(path), 0755)
    data, _ := json.MarshalIndent(settings, "", " ")
    return os.WriteFile(path, data, 0644)
}

func generateCommentID() string {
    return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[time.Now().Nanosecond()%len(charset)]
    }
    return string(b)
}
