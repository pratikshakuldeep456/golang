package designpattern

type Notfication interface {
	Send() error
}

type EmaiL struct {
	To      string
	Subject string
	Body    string
}

func (e *EmaiL) Send() error {
	// Logic to send email
	return nil
}

type SMS struct {
	To   string
	Body string
}

func (s *SMS) Send() error {
	// Logic to send SMS
	return nil
}

type Slack struct {
	Channel string
	Message string
}

func (s *Slack) Send() error {
	// Logic to send Slack message
	return nil
}

type NotificationService struct {
	notifier Notfication
}

func NewNotificationService(notifier Notfication) *NotificationService {
	return &NotificationService{notifier: notifier}
}

func (ns *NotificationService) Notify() error {
	return ns.notifier.Send()
}
