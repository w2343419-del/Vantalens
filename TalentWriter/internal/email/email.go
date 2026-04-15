package email

import (
    "bytes"
    "crypto/tls"
    "fmt"
    "log"
    "net/smtp"
    "strconv"
    "strings"
    "sync"
    "time"

    "vantalens/talentwriter/internal/models"
)

const (
    maxRetries    = 3
    queueSize     = 100
    workerCount   = 2
)

var (
    queue    = make(chan models.EmailJob, queueSize)
    stats    = struct {
        sync.Mutex
        sent    int64
        failed  int64
        retried int64
    }{}
)

func StartWorkers() {
    for i := 0; i < workerCount; i++ {
        go worker(i)
    }
    log.Printf("[EMAIL] Started %d email workers", workerCount)
}

func worker(id int) {
    for job := range queue {
        log.Printf("[EMAIL-WORKER-%d] Processing email for: %s", id, job.PostTitle)
        err := sendNotification(job.Settings, job.Comment, job.PostTitle)
        if err == nil {
            stats.Lock()
            stats.sent++
            stats.Unlock()
            log.Printf("[EMAIL-WORKER-%d] Email sent", id)
            continue
        }
        job.Retries++
        if job.Retries < maxRetries {
            waitTime := time.Duration(1<<uint(job.Retries)) * time.Second
            log.Printf("[EMAIL-WORKER-%d] Retry in %v", id, waitTime)
            time.Sleep(waitTime)
            stats.Lock()
            stats.retried++
            stats.Unlock()
            select {
            case queue <- job:
            default:
                stats.Lock()
                stats.failed++
                stats.Unlock()
            }
        } else {
            stats.Lock()
            stats.failed++
            stats.Unlock()
            log.Printf("[EMAIL-WORKER-%d] Failed after %d retries", id, maxRetries)
        }
    }
}

func QueueNotification(settings models.CommentSettings, comment models.Comment, postTitle string) {
    job := models.EmailJob{
        Settings:  settings,
        Comment:   comment,
        PostTitle: postTitle,
        CreatedAt: time.Now(),
    }
    select {
    case queue <- job:
        log.Printf("[EMAIL] Queued notification for: %s", postTitle)
    default:
        log.Printf("[EMAIL] Queue full, dropping notification")
    }
}

func sendNotification(settings models.CommentSettings, comment models.Comment, postTitle string) error {
    if !settings.SMTPEnabled || !settings.NotifyOnPending {
        return nil
    }
    from := settings.SMTPFrom
    if from == "" {
        from = settings.SMTPUser
    }
    if from == "" || len(settings.SMTPTo) == 0 || settings.SMTPHost == "" {
        return nil
    }
    subject := fmt.Sprintf("New Comment - %s", postTitle)
    body := fmt.Sprintf("Post: %s\nAuthor: %s\nEmail: %s\nContent:\n%s", postTitle, comment.Author, comment.Email, comment.Content)
    msg := bytes.NewBuffer(nil)
    msg.WriteString("From: " + from + "\r\n")
    msg.WriteString("To: " + strings.Join(settings.SMTPTo, ",") + "\r\n")
    msg.WriteString("Subject: " + subject + "\r\n")
    msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
    msg.WriteString(body)
    addr := settings.SMTPHost + ":" + strconv.Itoa(settings.SMTPPort)
    auth := smtp.PlainAuth("", settings.SMTPUser, settings.SMTPPass, settings.SMTPHost)
    if settings.SMTPPort == 465 {
        tlsConfig := &tls.Config{ServerName: settings.SMTPHost}
        conn, err := tls.Dial("tcp", addr, tlsConfig)
        if err != nil { return err }
        defer conn.Close()
        client, err := smtp.NewClient(conn, settings.SMTPHost)
        if err != nil { return err }
        defer client.Close()
        client.Auth(auth)
        client.Mail(from)
        for _, to := range settings.SMTPTo { client.Rcpt(to) }
        w, _ := client.Data()
        w.Write(msg.Bytes())
        w.Close()
        return nil
    }
    client, err := smtp.Dial(addr)
    if err != nil { return err }
    defer client.Close()
    client.StartTLS(&tls.Config{ServerName: settings.SMTPHost})
    client.Auth(auth)
    client.Mail(from)
    for _, to := range settings.SMTPTo { client.Rcpt(to) }
    w, _ := client.Data()
    w.Write(msg.Bytes())
    w.Close()
    return nil
}
