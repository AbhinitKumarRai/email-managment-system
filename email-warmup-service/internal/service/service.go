package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type domainWorker struct {
	ratePerMinute int
	queue         chan *sendJob
	stop          chan struct{}
}

type sendJob struct {
	email     *EmailMessage
	recipient string
	tries     int
}

type Service struct {
	storage       *Storage
	domainLimits  map[string]int
	domainWorkers map[string]*domainWorker
	workerMutex   sync.Mutex
}

func NewService() *Service {
	return &Service{
		storage:       NewStorage(),
		domainLimits:  defaultDomainLimits(),
		domainWorkers: make(map[string]*domainWorker),
	}
}

func defaultDomainLimits() map[string]int {
	return map[string]int{
		"gmail.com":   10, // 10 emails/minute
		"yahoo.com":   5,
		"outlook.com": 8,
		// Add more as needed
		"default": 5,
	}
}

func (s *Service) getDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "default"
	}
	domain := strings.ToLower(parts[1])
	if _, ok := s.domainLimits[domain]; ok {
		return domain
	}
	return "default"
}

func (s *Service) getOrStartWorker(domain string) *domainWorker {
	s.workerMutex.Lock()
	defer s.workerMutex.Unlock()
	w, ok := s.domainWorkers[domain]
	if !ok {
		rate := s.domainLimits[domain]
		w = &domainWorker{
			ratePerMinute: rate,
			queue:         make(chan *sendJob, 1000),
			stop:          make(chan struct{}),
		}
		s.domainWorkers[domain] = w
		go s.runWorker(domain, w)
	}
	return w
}

func (s *Service) runWorker(domain string, w *domainWorker) {
	interval := time.Minute / time.Duration(w.ratePerMinute)
	if w.ratePerMinute <= 0 {
		interval = time.Minute / 1
	}
	for {
		select {
		case job := <-w.queue:
			if err := s.sendSMTP(job.email, job.recipient); err != nil {
				if job.tries < 3 {
					job.tries++
					time.AfterFunc(time.Second*5, func() { w.queue <- job })
					log.Printf("Retrying email %s to %s (try %d)", job.email.ID, job.recipient, job.tries)
				} else {
					log.Printf("Failed to send email %s to %s after 3 tries", job.email.ID, job.recipient)
				}
			}
			time.Sleep(interval)
		case <-w.stop:
			return
		}
	}
}

func (s *Service) RegisterEmail(addr string) error {
	if !s.storage.RegisterEmail(addr) {
		return fmt.Errorf("email already registered")
	}
	// Send all previous emails to this new address, rate-limited by domain
	domain := s.getDomain(addr)
	worker := s.getOrStartWorker(domain)
	emails := s.storage.ListEmails()
	for _, email := range emails {
		worker.queue <- &sendJob{email: email, recipient: addr, tries: 0}
	}
	return nil
}

func (s *Service) BroadcastEmail(from, subject, body string) (string, error) {
	emails := s.storage.ListRegisteredEmails()
	if len(emails) == 0 {
		return "", fmt.Errorf("no registered emails")
	}
	id := generateID()
	toList := make([]string, 0, len(emails))
	for _, e := range emails {
		toList = append(toList, e.Address)
	}
	email := &EmailMessage{
		ID:        id,
		From:      from,
		To:        toList,
		Subject:   subject,
		Body:      body,
		Timestamp: Now(),
	}
	s.storage.AddEmail(email)
	for _, to := range toList {
		go s.sendSMTP(email, to)
	}
	return id, nil
}

func (s *Service) ListEmails() []*EmailMessage {
	return s.storage.ListEmails()
}

// sendSMTP now returns error
func (s *Service) sendSMTP(email *EmailMessage, recipient string) error {
	server := "localhost:25" // Change to your SMTP relay
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Printf("SMTP connect error: %v", err)
		return err
	}
	defer conn.Close()

	r := make([]byte, 1024)
	conn.Read(r) // banner
	fmt.Fprintf(conn, "HELO localhost\r\n")
	conn.Read(r)
	fmt.Fprintf(conn, "MAIL FROM:<%s>\r\n", email.From)
	conn.Read(r)
	fmt.Fprintf(conn, "RCPT TO:<%s>\r\n", recipient)
	conn.Read(r)
	fmt.Fprintf(conn, "DATA\r\n")
	conn.Read(r)
	fmt.Fprintf(conn, "Subject: %s\r\nTo: %s\r\nFrom: %s\r\nDate: %s\r\nMessage-ID: <%s@localhost>\r\n\r\n%s\r\n.\r\n",
		email.Subject, recipient, email.From, email.Timestamp.Format(time.RFC1123Z), email.ID, email.Body)
	conn.Read(r)
	fmt.Fprintf(conn, "QUIT\r\n")
	conn.Read(r)
	log.Printf("Sent email %s to %s", email.ID, recipient)
	return nil
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
