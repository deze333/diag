package diag

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/deze333/m8l"
)

//------------------------------------------------------------
//
//------------------------------------------------------------

type EmailNotifier struct {
	sender        map[string]string
	recipient     map[string]string
	subjectPrefix string
	sendProc      func(sender, recipient map[string]string, subj, body string)
}

var _emailNotifier *EmailNotifier

//------------------------------------------------------------
//
//------------------------------------------------------------

func SetEmailNotification(sender, recipient map[string]string, subjPrefix string) {

	_emailNotifier = &EmailNotifier{
		sender:        sender,
		recipient:     recipient,
		subjectPrefix: subjPrefix,
	}
}

func SetEmailNotificationProc(sender, recipient map[string]string, subjPrefix string, sendProc func(sender, recipient map[string]string, subj, body string)) {

	_emailNotifier = &EmailNotifier{
		sender:        sender,
		recipient:     recipient,
		subjectPrefix: subjPrefix,
		sendProc:      sendProc,
	}
}

//------------------------------------------------------------
//
//------------------------------------------------------------

type Email struct {
	Tag        string
	Subject    string
	ParamsHtml string
	Params     []string
}

// Email template
var _emailTpl = template.Must(template.New("deze333/diag/email").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
</head>
<body>
<div style="font-family: monospace; font-size: 13px; color: maroon;">
{{.Tag}}<br><strong>{{.Subject}}</strong>
</div>
<hr>
<div style="font-family: monospace; font-size: 11px;">
{{range .Params}}{{.}}<br>
{{end}}
</div>
</body>
</html>
`))

//------------------------------------------------------------
//
//------------------------------------------------------------

func notifyEmail(name, title string, args ...interface{}) {
	if _emailNotifier == nil {
		return
	}

	subj := "[" + _emailNotifier.subjectPrefix + "] " + name + " : " + title

	e := Email{
		Tag:     name,
		Subject: title,
		Params:  []string{},
	}

	// Convert to strings
	// Change \n to <br> in those who have it
	ss := make([]string, len(args))
	for i, v := range args {
		s := fmt.Sprint(v)
		if strings.Contains(s, "\n") {
			ss[i] = strings.Replace(s, "\n", "<br>", -1)
		} else {
			ss[i] = s
		}
	}

	if len(ss) == 1 {
		e.Params = append(e.Params, fmt.Sprint(ss[0]))
	} else {
		for i := 0; i+1 < len(ss); i += 2 {
			key := ss[i]
			val := strings.Replace(ss[i+1], "\n", "<br>", -1)
			e.Params = append(e.Params, fmt.Sprintf("<strong>%s</strong> = %s<br>", key, val))
		}
	}

	var err error
	var msg bytes.Buffer
	err = _emailTpl.Execute(&msg, e)
	if err != nil {
		ERROR("diag", "Error generating SOS email via template. Email send aborted.", "err", err)
		return
	}

	// Async send email
	if _emailNotifier.sendProc != nil {

		// Via send proc
		go _emailNotifier.sendProc(
			_emailNotifier.sender, _emailNotifier.recipient,
			subj, msg.String())

	} else {

		// Via SMTP
		email := m8l.NewEmail(subj, &msg)
		email.SetSender(_emailNotifier.sender)
		email.SetReplyTo(_emailNotifier.recipient["identity"], _emailNotifier.recipient["email"])
		email.AddTo(_emailNotifier.recipient["identity"], _emailNotifier.recipient["email"])
		if err = email.Validate(); err != nil {
			ERROR("diag", "Error validating email. Email send aborted.", "err", err)
			return
		}
		if err = email.SendAsync(); err != nil {
			ERROR("diag", "Error sending email. Email send aborted.", "err", err)
			return
		}
	}
}
